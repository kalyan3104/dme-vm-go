package tests

import (
	"testing"

	logger "github.com/kalyan3104/dme-logger-go"
	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/config"
	"github.com/kalyan3104/dme-vm-go/ipc/common"
	"github.com/kalyan3104/dme-vm-go/ipc/nodepart"
	"github.com/kalyan3104/dme-vm-go/mock"
	"github.com/stretchr/testify/require"
)

var arwenVirtualMachine = []byte{5, 0}

func TestArwenDriver_DiagnoseWait(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(t, blockchain)

	err := driver.DiagnoseWait(100)
	require.Nil(t, err)
}

func TestArwenDriver_DiagnoseWaitWithTimeout(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(t, blockchain)

	err := driver.DiagnoseWait(5000)
	require.True(t, common.IsCriticalError(err))
	require.Contains(t, err.Error(), "timeout")
	require.True(t, driver.IsClosed())
}

func TestArwenDriver_RestartsIfStopped(t *testing.T) {
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

	// Per this request, Arwen is restarted
	vmOutput, err = driver.RunSmartContractCreate(createDeployInput(bytecodeCounter))
	require.Nil(t, err)
	require.NotNil(t, vmOutput)
	require.False(t, driver.IsClosed())
}

func BenchmarkArwenDriver_RestartsIfStopped(b *testing.B) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(b, blockchain)

	for i := 0; i < b.N; i++ {
		driver.Close()
		require.True(b, driver.IsClosed())
		_ = driver.RestartArwenIfNecessary()
		require.False(b, driver.IsClosed())
	}
}

func BenchmarkArwenDriver_RestartArwenIfNecessary(b *testing.B) {
	blockchain := &mock.BlockchainHookStub{}
	driver := newDriver(b, blockchain)

	for i := 0; i < b.N; i++ {
		_ = driver.RestartArwenIfNecessary()
	}
}

func newDriver(tb testing.TB, blockchain *mock.BlockchainHookStub) *nodepart.ArwenDriver {
	driver, err := nodepart.NewArwenDriver(
		blockchain,
		common.ArwenArguments{
			VMHostParameters: arwen.VMHostParameters{
				VMType:                       arwenVirtualMachine,
				BlockGasLimit:                uint64(10000000),
				GasSchedule:                  config.MakeGasMapForTests(),
				Kalyan3104ProtectedKeyPrefix: []byte("KALYAN3104"),
			},
		},
		nodepart.Config{MaxLoopTime: 1000},
	)
	require.Nil(tb, err)
	require.NotNil(tb, driver)
	require.False(tb, driver.IsClosed())
	return driver
}
