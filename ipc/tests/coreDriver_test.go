package tests

import (
	"testing"

	"github.com/Dharitri-org/sme-core-vm-go/config"
	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/nodepart"
	"github.com/Dharitri-org/sme-core-vm-go/mock"
	logger "github.com/Dharitri-org/sme-logger"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
	"github.com/stretchr/testify/require"
)

var coreVirtualMachine = []byte{5, 0}

func TestCoreDriver_DiagnoseWait(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(t, blockchain)

	err := driver.DiagnoseWait(100)
	require.Nil(t, err)
}

func TestCoreDriver_DiagnoseWaitWithTimeout(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(t, blockchain)

	err := driver.DiagnoseWait(5000)
	require.True(t, common.IsCriticalError(err))
	require.Contains(t, err.Error(), "timeout")
	require.True(t, driver.IsClosed())
}

func TestCoreDriver_RestartsIfStopped(t *testing.T) {
	logger.ToggleLoggerName(true)
	_ = logger.SetLogLevel("*:TRACE")

	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(t, blockchain)

	blockchain.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &mock.AccountMock{Code: bytecodeCounter}, nil
	}

	vmOutput, err := driver.RunSmartContractCreate(createDeployInput(bytecodeCounter))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	vmOutput, err = driver.RunSmartContractCall(createCallInput("increment"))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)

	require.False(t, driver.IsClosed())
	driver.Close()
	require.True(t, driver.IsClosed())

	// Per this request, Core is restarted
	vmOutput, err = driver.RunSmartContractCreate(createDeployInput(bytecodeCounter))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.False(t, driver.IsClosed())
}

func BenchmarkCoreDriver_RestartsIfStopped(b *testing.B) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(b, blockchain)

	for i := 0; i < b.N; i++ {
		driver.Close()
		require.True(b, driver.IsClosed())
		_ = driver.RestartCoreIfNecessary()
		require.False(b, driver.IsClosed())
	}
}

func BenchmarkCoreDriver_RestartCoreIfNecessary(b *testing.B) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(b, blockchain)

	for i := 0; i < b.N; i++ {
		_ = driver.RestartCoreIfNecessary()
	}
}

func newDriver(tb testing.TB, blockchain *mock.BlockchainHookStub) *nodepart.CoreDriver {
	driver, err := nodepart.NewCoreDriver(
		blockchain,
		common.CoreArguments{
			VMHostParameters: core.VMHostParameters{
				VMType:                     coreVirtualMachine,
				BlockGasLimit:              uint64(10000000),
				GasSchedule:                config.MakeGasMapForTests(),
				DharitriProtectedKeyPrefix: []byte("DHARITRI"),
			},
		},
		nodepart.Config{MaxLoopTime: 1000},
	)
	require.Nil(tb, err)
	require.NotNil(tb, driver)
	require.False(tb, driver.IsClosed())
	return driver
}

