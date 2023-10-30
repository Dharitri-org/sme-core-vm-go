package coremandos

import (
	"github.com/Dharitri-org/sme-core-vm-go/config"
	core "github.com/Dharitri-org/sme-core-vm-go/core"
	coreHost "github.com/Dharitri-org/sme-core-vm-go/core/host"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
	vmi "github.com/Dharitri-org/sme-vm-common"
	worldhook "github.com/Dharitri-org/sme-vm-util/mock-hook-blockchain"
	cryptohook "github.com/Dharitri-org/sme-vm-util/mock-hook-crypto"
	mc "github.com/Dharitri-org/sme-vm-util/test-util/mandos/controller"
	mjparse "github.com/Dharitri-org/sme-vm-util/test-util/mandos/json/parse"
)

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// CoreTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Core.
type CoreTestExecutor struct {
	fileResolver mjparse.FileResolver
	World        *worldhook.BlockchainHookMock
	vm           vmi.VMExecutionHandler
	checkGas     bool
}

var _ mc.TestExecutor = (*CoreTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*CoreTestExecutor)(nil)

// NewCoreTestExecutor prepares a new CoreTestExecutor instance.
func NewCoreTestExecutor() (*CoreTestExecutor, error) {
	world := worldhook.NewMock()
	world.EnableMockAddressGeneration()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMapForTests()
	vm, err := coreHost.NewCoreVM(world, cryptohook.KryptoHookMockInstance, &core.VMHostParameters{
		VMType:                     TestVMType,
		BlockGasLimit:              blockGasLimit,
		GasSchedule:                gasSchedule,
		ProtocolBuiltinFunctions:   make(vmcommon.FunctionNames),
		DharitriProtectedKeyPrefix: []byte(DharitriProtectedKeyPrefix),
	})
	if err != nil {
		return nil, err
	}
	return &CoreTestExecutor{
		fileResolver: nil,
		World:        world,
		vm:           vm,
		checkGas:     true,
	}, nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *CoreTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}
