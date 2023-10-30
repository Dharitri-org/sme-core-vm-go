package nodepart

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/marshaling"
	logger "github.com/Dharitri-org/sme-logger"
	"github.com/Dharitri-org/sme-logger/pipes"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var log = logger.GetOrCreate("coreDriver")

var _ vmcommon.VMExecutionHandler = (*CoreDriver)(nil)

// CoreDriver manages the execution of the Core process
type CoreDriver struct {
	blockchainHook      vmcommon.BlockchainHook
	coreArguments       common.CoreArguments
	config              Config
	logsMarshalizer     marshaling.Marshalizer
	messagesMarshalizer marshaling.Marshalizer

	coreInitRead    *os.File
	coreInitWrite   *os.File
	coreInputRead   *os.File
	coreInputWrite  *os.File
	coreOutputRead  *os.File
	coreOutputWrite *os.File

	counterDeploy uint64
	counterCall   uint64

	command  *exec.Cmd
	part     *NodePart
	logsPart ParentLogsPart
}

// NewCoreDriver creates a new driver
func NewCoreDriver(
	blockchainHook vmcommon.BlockchainHook,
	coreArguments common.CoreArguments,
	config Config,
) (*CoreDriver, error) {
	driver := &CoreDriver{
		blockchainHook:      blockchainHook,
		coreArguments:       coreArguments,
		config:              config,
		logsMarshalizer:     marshaling.CreateMarshalizer(coreArguments.LogsMarshalizer),
		messagesMarshalizer: marshaling.CreateMarshalizer(coreArguments.MessagesMarshalizer),
	}

	err := driver.startCore()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (driver *CoreDriver) startCore() error {
	log.Info("CoreDriver.startCore()")

	logsProfileReader, logsWriter, err := driver.resetLogsPart()
	if err != nil {
		return err
	}

	err = driver.resetPipeStreams()
	if err != nil {
		return err
	}

	corePath, err := driver.getCorePath()
	if err != nil {
		return err
	}

	driver.command = exec.Command(corePath)
	driver.command.ExtraFiles = []*os.File{
		driver.coreInitRead,
		driver.coreInputRead,
		driver.coreOutputWrite,
		logsProfileReader,
		logsWriter,
	}

	coreStdout, err := driver.command.StdoutPipe()
	if err != nil {
		return err
	}

	coreStderr, err := driver.command.StderrPipe()
	if err != nil {
		return err
	}

	err = driver.command.Start()
	if err != nil {
		return err
	}

	err = common.SendCoreArguments(driver.coreInitWrite, driver.coreArguments)
	if err != nil {
		return err
	}

	driver.part, err = NewNodePart(
		driver.coreOutputRead,
		driver.coreInputWrite,
		driver.blockchainHook,
		driver.config,
		driver.messagesMarshalizer,
	)
	if err != nil {
		return err
	}

	driver.logsPart.StartLoop(coreStdout, coreStderr)

	return nil
}

func (driver *CoreDriver) resetLogsPart() (*os.File, *os.File, error) {
	logsPart, err := pipes.NewParentPart("Core", driver.logsMarshalizer)
	if err != nil {
		return nil, nil, err
	}

	driver.logsPart = logsPart
	readProfile, writeLogs := logsPart.GetChildPipes()
	return readProfile, writeLogs, nil
}

func (driver *CoreDriver) resetPipeStreams() error {
	closeFile(driver.coreInitRead)
	closeFile(driver.coreInitWrite)
	closeFile(driver.coreInputRead)
	closeFile(driver.coreInputWrite)
	closeFile(driver.coreOutputRead)
	closeFile(driver.coreOutputWrite)

	var err error

	driver.coreInitRead, driver.coreInitWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.coreInputRead, driver.coreInputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.coreOutputRead, driver.coreOutputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	return nil
}

func closeFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot close file.\n")
		}
	}
}

// RestartCoreIfNecessary restarts Core if the process is closed
func (driver *CoreDriver) RestartCoreIfNecessary() error {
	if !driver.IsClosed() {
		return nil
	}

	err := driver.startCore()
	return err
}

// IsClosed checks whether the Core process is closed
func (driver *CoreDriver) IsClosed() bool {
	pid := driver.command.Process.Pid
	process, err := os.FindProcess(pid)
	if err != nil {
		return true
	}

	err = process.Signal(syscall.Signal(0))
	return err != nil
}

// RunSmartContractCreate sends a deploy request to Core and waits for the output
func (driver *CoreDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	driver.counterDeploy++
	log.Trace("RunSmartContractCreate", "counter", driver.counterDeploy)

	err := driver.RestartCoreIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractDeployRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Warn("RunSmartContractCreate", "err", err)
		_ = driver.Close()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.VMOutput, response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// RunSmartContractCall sends an execution request to Core and waits for the output
func (driver *CoreDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	driver.counterCall++
	log.Trace("RunSmartContractCall", "counter", driver.counterCall, "func", input.Function, "sc", input.RecipientAddr)

	err := driver.RestartCoreIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractCallRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Warn("RunSmartContractCall", "err", err)
		_ = driver.Close()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.VMOutput, response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// DiagnoseWait sends a diagnose message to Core
func (driver *CoreDriver) DiagnoseWait(milliseconds uint32) error {
	err := driver.RestartCoreIfNecessary()
	if err != nil {
		return common.WrapCriticalError(err)
	}

	request := common.NewMessageDiagnoseWaitRequest(milliseconds)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		log.Error("DiagnoseWait", "err", err)
		_ = driver.Close()
		return common.WrapCriticalError(err)
	}

	return response.GetError()
}

// Close stops Core
func (driver *CoreDriver) Close() error {
	driver.logsPart.StopLoop()

	err := driver.stopCore()
	if err != nil {
		log.Error("CoreDriver.Close()", "err", err)
		return err
	}

	return nil
}

func (driver *CoreDriver) stopCore() error {
	err := driver.command.Process.Kill()
	if err != nil {
		return err
	}

	_, err = driver.command.Process.Wait()
	if err != nil {
		return err
	}

	return nil
}
