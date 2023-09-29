package tests

import (
	"os"
	"sync"
	"testing"

	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/config"
	"github.com/kalyan3104/dme-vm-go/ipc/arwenpart"
	"github.com/kalyan3104/dme-vm-go/ipc/common"
	"github.com/kalyan3104/dme-vm-go/ipc/marshaling"
	"github.com/kalyan3104/dme-vm-go/ipc/nodepart"
	"github.com/kalyan3104/dme-vm-go/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testFiles struct {
	outputOfNode  *os.File
	inputOfArwen  *os.File
	outputOfArwen *os.File
	inputOfNode   *os.File
}

func TestArwenPart_SendDeployRequest(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}

	response, err := doContractRequest(t, "2", createDeployRequest(bytecodeCounter), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequestWhenNoContract(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}

	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func TestArwenPart_SendCallRequest(t *testing.T) {
	blockchain := &mock.BlockchainHookStub{}

	blockchain.GetUserAccountCalled = func(address []byte) (vmcommon.UserAccountHandler, error) {
		return &mock.AccountMock{Code: bytecodeCounter}, nil
	}

	response, err := doContractRequest(t, "3", createCallRequest("increment"), blockchain)
	require.NotNil(t, response)
	require.Nil(t, err)
}

func doContractRequest(
	t *testing.T,
	tag string,
	request common.MessageHandler,
	blockchain vmcommon.BlockchainHook,
) (common.MessageHandler, error) {
	files := createTestFiles(t, tag)
	var response common.MessageHandler
	var responseError error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		vmHostParameters := &arwen.VMHostParameters{
			VMType:                       []byte{5, 0},
			BlockGasLimit:                uint64(10000000),
			GasSchedule:                  config.MakeGasMapForTests(),
			Kalyan3104ProtectedKeyPrefix: []byte("KALYAN3104"),
		}

		part, err := arwenpart.NewArwenPart(
			files.inputOfArwen,
			files.outputOfArwen,
			vmHostParameters,
			marshaling.CreateMarshalizer(marshaling.JSON),
		)
		assert.Nil(t, err)
		_ = part.StartLoop()
		wg.Done()
	}()

	go func() {
		part, err := nodepart.NewNodePart(
			files.inputOfNode,
			files.outputOfNode,
			blockchain,
			nodepart.Config{MaxLoopTime: 1000},
			marshaling.CreateMarshalizer(marshaling.JSON),
		)
		assert.Nil(t, err)
		response, responseError = part.StartLoop(request)
		_ = part.SendStopSignal()
		wg.Done()
	}()

	wg.Wait()

	return response, responseError
}

func createTestFiles(t *testing.T, tag string) testFiles {
	files := testFiles{}

	var err error
	files.inputOfArwen, files.outputOfNode, err = os.Pipe()
	require.Nil(t, err)
	files.inputOfNode, files.outputOfArwen, err = os.Pipe()
	require.Nil(t, err)

	return files
}

func createDeployRequest(contractCode []byte) common.MessageHandler {
	return common.NewMessageContractDeployRequest(createDeployInput(contractCode))
}

func createCallRequest(function string) common.MessageHandler {
	return common.NewMessageContractCallRequest(createCallInput(function))
}
