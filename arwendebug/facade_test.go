package arwendebug

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var databasePath = "./testdata/db"
var wasmCounterPath = "../test/contracts/counter/counter.wasm"
var wasmErc20Path = "../test/contracts/erc20/erc20.wasm"

func init() {
	_ = os.RemoveAll(databasePath)
}

func TestFacade_CreateAccount(t *testing.T) {
	context := newTestContext(t)
	context.createAccount(newDummyAddress("alice").hex, "42")

	require.True(t, context.accountExists(newDummyAddress("alice").raw))
}

func TestFacade_RunContract_Counter(t *testing.T) {
	context := newTestContext(t)

	counterKey := []byte{'m', 'y', 'c', 'o', 'u', 'n', 't', 'e', 'r', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	alice := newDummyAddress("alice")
	context.createAccount(alice.hex, "42")
	deployResponse := context.deployContract(wasmCounterPath, alice.hex)
	contractAddress := deployResponse.ContractAddress
	contractAddressHex := deployResponse.ContractAddressHex
	require.Equal(t, "contract0000000000000000000alice", string(contractAddress))
	require.True(t, context.accountExists([]byte(contractAddress)))

	context.runContract(contractAddressHex, alice.hex, "increment")
	counterValue := context.queryContract(contractAddressHex, alice.hex, "get").getFirstResultAsInt64()
	require.Equal(t, int64(2), counterValue)

	world := context.loadWorld()
	state, err := world.blockchainHook.GetAllState([]byte(contractAddress))
	require.Nil(t, err)
	require.NotNil(t, state)
	require.Equal(t, []byte{2}, state[string(counterKey)])
}

func TestFacade_RunContract_ERC20(t *testing.T) {
	context := newTestContext(t)

	alice := newDummyAddress("alice")
	bob := newDummyAddress("bob")
	carol := newDummyAddress("carol")

	context.createAccount(alice.hex, "42")
	deployResponse := context.deployContract(wasmErc20Path, alice.hex, "64")
	contractAddress := deployResponse.ContractAddress
	contractAddressHex := deployResponse.ContractAddressHex
	require.Equal(t, "contract0000000000000000000alice", string(contractAddress))

	// Initial state
	totalSupply := context.queryContract(contractAddressHex, alice.hex, "totalSupply").getFirstResultAsInt64()
	balanceOfAlice := context.queryContract(contractAddressHex, alice.hex, "balanceOf", alice.hex).getFirstResultAsInt64()
	balanceOfBob := context.queryContract(contractAddressHex, alice.hex, "balanceOf", bob.hex).getFirstResultAsInt64()
	balanceOfCarol := context.queryContract(contractAddressHex, alice.hex, "balanceOf", carol.hex).getFirstResultAsInt64()
	require.Equal(t, int64(100), totalSupply)
	require.Equal(t, int64(100), balanceOfAlice)
	require.Equal(t, int64(0), balanceOfBob)
	require.Equal(t, int64(0), balanceOfCarol)

	// Transfers
	context.runContract(contractAddressHex, alice.hex, "transferToken", alice.hex, "0A")
	context.runContract(contractAddressHex, alice.hex, "transferToken", bob.hex, "0A")
	context.runContract(contractAddressHex, alice.hex, "transferToken", carol.hex, "0A")
	context.runContract(contractAddressHex, bob.hex, "transferToken", carol.hex, "05")

	balanceOfAlice = context.queryContract(contractAddressHex, alice.hex, "balanceOf", alice.hex).getFirstResultAsInt64()
	balanceOfBob = context.queryContract(contractAddressHex, alice.hex, "balanceOf", bob.hex).getFirstResultAsInt64()
	balanceOfCarol = context.queryContract(contractAddressHex, alice.hex, "balanceOf", carol.hex).getFirstResultAsInt64()
	require.Equal(t, int64(80), balanceOfAlice)
	require.Equal(t, int64(5), balanceOfBob)
	require.Equal(t, int64(15), balanceOfCarol)
}

func TestFacade_UpgradeContract_CounterToERC20(t *testing.T) {
	context := newTestContext(t)

	alice := newDummyAddress("alice")
	bob := newDummyAddress("bob")
	context.createAccount(alice.hex, "42")

	// Deploy counter & smoke test
	deployResponse := context.deployContract(wasmCounterPath, alice.hex)
	contractAddressHex := deployResponse.ContractAddressHex
	context.runContract(contractAddressHex, alice.hex, "increment")
	counterValue := context.queryContract(contractAddressHex, alice.hex, "get").getFirstResultAsInt64()
	require.Equal(t, int64(2), counterValue)

	// Upgrade to ERC20 & smoke test
	_ = context.upgradeContract(contractAddressHex, wasmErc20Path, alice.hex, "64")
	context.runContract(contractAddressHex, alice.hex, "transferToken", bob.hex, "0A")
	balanceOfAlice := context.queryContract(contractAddressHex, alice.hex, "balanceOf", alice.hex).getFirstResultAsInt64()
	balanceOfBob := context.queryContract(contractAddressHex, alice.hex, "balanceOf", bob.hex).getFirstResultAsInt64()
	require.Equal(t, int64(90), balanceOfAlice)
	require.Equal(t, int64(10), balanceOfBob)
}
