package host

import (
	"fmt"

	"github.com/Dharitri-org/sme-core-vm-go/config"
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/core/contexts"
	"github.com/Dharitri-org/sme-core-vm-go/core/cryptoapi"
	"github.com/Dharitri-org/sme-core-vm-go/core/dharitriapi"
	"github.com/Dharitri-org/sme-core-vm-go/crypto"
	"github.com/Dharitri-org/sme-core-vm-go/wasmer"
	"github.com/Dharitri-org/sme-dharitri/core/atomic"
	logger "github.com/Dharitri-org/sme-logger"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var log = logger.GetOrCreate("core/host")

var MaximumWasmerInstanceCount = uint64(10)

// TryFunction corresponds to the try() part of a try / catch block
type TryFunction func()

// CatchFunction corresponds to the catch() part of a try / catch block
type CatchFunction func(error)

// vmHost implements HostContext interface.
type vmHost struct {
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     crypto.VMCrypto

	ethInput []byte

	blockchainContext core.BlockchainContext
	runtimeContext    core.RuntimeContext
	outputContext     core.OutputContext
	meteringContext   core.MeteringContext
	storageContext    core.StorageContext
	bigIntContext     core.BigIntContext

	scAPIMethods             *wasmer.Imports
	protocolBuiltinFunctions vmcommon.FunctionNames

	coreV2EnableEpoch uint32
	flagCoreV2        atomic.Flag
}

// NewCoreVM creates a new Core vmHost
func NewCoreVM(
	blockChainHook vmcommon.BlockchainHook,
	hostParameters *core.VMHostParameters,
) (*vmHost, error) {

	cryptoHook := crypto.NewVMCrypto()
	host := &vmHost{
		blockChainHook:           blockChainHook,
		cryptoHook:               cryptoHook,
		meteringContext:          nil,
		runtimeContext:           nil,
		blockchainContext:        nil,
		storageContext:           nil,
		bigIntContext:            nil,
		scAPIMethods:             nil,
		protocolBuiltinFunctions: hostParameters.ProtocolBuiltinFunctions,
		coreV2EnableEpoch:        hostParameters.CoreV2EnableEpoch,
	}

	var err error

	imports, err := dharitriapi.DharitriEIImports()
	if err != nil {
		return nil, err
	}

	imports, err = dharitriapi.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = dharitriapi.SmallIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	err = wasmer.SetImports(imports)
	if err != nil {
		return nil, err
	}

	host.scAPIMethods = imports

	host.blockchainContext, err = contexts.NewBlockchainContext(host, blockChainHook)
	if err != nil {
		return nil, err
	}

	host.runtimeContext, err = contexts.NewRuntimeContext(
		host,
		hostParameters.VMType,
		hostParameters.UseWarmInstance,
	)
	if err != nil {
		return nil, err
	}

	host.meteringContext, err = contexts.NewMeteringContext(host, hostParameters.GasSchedule, hostParameters.BlockGasLimit)
	if err != nil {
		return nil, err
	}

	host.outputContext, err = contexts.NewOutputContext(host)
	if err != nil {
		return nil, err
	}

	host.storageContext, err = contexts.NewStorageContext(host, blockChainHook, hostParameters.DharitriProtectedKeyPrefix)
	if err != nil {
		return nil, err
	}

	host.bigIntContext, err = contexts.NewBigIntContext()
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(hostParameters.GasSchedule)
	if err != nil {
		return nil, err
	}

	host.runtimeContext.SetMaxInstanceCount(MaximumWasmerInstanceCount)

	opcodeCosts := gasCostConfig.WASMOpcodeCost.ToOpcodeCostsArray()
	wasmer.SetOpcodeCosts(&opcodeCosts)

	host.initContexts()

	return host, nil
}

func (host *vmHost) Crypto() crypto.VMCrypto {
	return host.cryptoHook
}

func (host *vmHost) Blockchain() core.BlockchainContext {
	return host.blockchainContext
}

func (host *vmHost) Runtime() core.RuntimeContext {
	return host.runtimeContext
}

func (host *vmHost) Output() core.OutputContext {
	return host.outputContext
}

func (host *vmHost) Metering() core.MeteringContext {
	return host.meteringContext
}

func (host *vmHost) Storage() core.StorageContext {
	return host.storageContext
}

func (host *vmHost) BigInt() core.BigIntContext {
	return host.bigIntContext
}

func (host *vmHost) IsCoreV2Enabled() bool {
	return host.flagCoreV2.IsSet()
}

func (host *vmHost) GetContexts() (
	core.BigIntContext,
	core.BlockchainContext,
	core.MeteringContext,
	core.OutputContext,
	core.RuntimeContext,
	core.StorageContext,
) {
	return host.bigIntContext,
		host.blockchainContext,
		host.meteringContext,
		host.outputContext,
		host.runtimeContext,
		host.storageContext
}

func (host *vmHost) InitState() {
	host.initContexts()
	host.flagCoreV2.Toggle(host.blockChainHook.CurrentEpoch() >= host.coreV2EnableEpoch)
	log.Trace("coreV2", "enabled", host.flagCoreV2.IsSet())
}

func (host *vmHost) initContexts() {
	host.ClearContextStateStack()
	host.bigIntContext.InitState()
	host.outputContext.InitState()
	host.runtimeContext.InitState()
	host.storageContext.InitState()
	host.ethInput = nil
}

func (host *vmHost) ClearContextStateStack() {
	host.bigIntContext.ClearStateStack()
	host.outputContext.ClearStateStack()
	host.runtimeContext.ClearStateStack()
	host.storageContext.ClearStateStack()
}

func (host *vmHost) Clean() {
	if host.runtimeContext.IsWarmInstance() {
		return
	}
	host.runtimeContext.CleanWasmerInstance()
	core.RemoveAllHostContexts()
}

func (host *vmHost) GetAPIMethods() *wasmer.Imports {
	return host.scAPIMethods
}

func (host *vmHost) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return host.protocolBuiltinFunctions
}

func (host *vmHost) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (vmOutput *vmcommon.VMOutput, err error) {
	log.Trace("RunSmartContractCreate begin", "len(code)", len(input.ContractCode), "metadata", input.ContractCodeMetadata)

	try := func() {
		vmOutput = host.doRunSmartContractCreate(input)
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCreate", "error", err)
	}

	TryCatch(try, catch, "core.RunSmartContractCreate")
	if vmOutput != nil {
		log.Trace("RunSmartContractCreate end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage)
	}

	return
}

func (host *vmHost) RunSmartContractCall(input *vmcommon.ContractCallInput) (vmOutput *vmcommon.VMOutput, err error) {
	log.Trace("RunSmartContractCall begin", "function", input.Function)

	tryUpgrade := func() {
		vmOutput = host.doRunSmartContractUpgrade(input)
	}

	tryCall := func() {
		vmOutput = host.doRunSmartContractCall(input)

		if host.hasRetriableExecutionError(vmOutput) {
			log.Error("Retriable execution error detected. Will reset warm Wasmer instance.")
			host.runtimeContext.ResetWarmInstance()
		}
	}

	catch := func(caught error) {
		err = caught
		log.Error("RunSmartContractCall", "error", err)
	}

	isUpgrade := input.Function == core.UpgradeFunctionName
	if isUpgrade {
		TryCatch(tryUpgrade, catch, "core.RunSmartContractUpgrade")
	} else {
		TryCatch(tryCall, catch, "core.RunSmartContractCall")
	}

	if vmOutput != nil {
		log.Debug("RunSmartContractCall end", "returnCode", vmOutput.ReturnCode, "returnMessage", vmOutput.ReturnMessage, "function", input.Function)
	}

	return
}

// TryCatch simulates a try/catch block using golang's recover() functionality
func TryCatch(try TryFunction, catch CatchFunction, catchFallbackMessage string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%s, panic: %v", catchFallbackMessage, r)
			}

			catch(err)
		}
	}()

	try()
}

func (host *vmHost) hasRetriableExecutionError(vmOutput *vmcommon.VMOutput) bool {
	if !host.runtimeContext.IsWarmInstance() {
		return false
	}

	return vmOutput.ReturnMessage == "allocation error"
}
