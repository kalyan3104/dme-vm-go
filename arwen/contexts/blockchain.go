package contexts

import (
	"math/big"

	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
)

type blockchainContext struct {
	host           arwen.VMHost
	blockChainHook vmcommon.BlockchainHook
}

// NewBlockchainContext creates a new blockchainContext
func NewBlockchainContext(
	host arwen.VMHost,
	blockChainHook vmcommon.BlockchainHook,
) (*blockchainContext, error) {

	context := &blockchainContext{
		blockChainHook: blockChainHook,
		host:           host,
	}

	return context, nil
}

func (context *blockchainContext) NewAddress(creatorAddress []byte) ([]byte, error) {
	nonce, err := context.GetNonce(creatorAddress)
	if err != nil {
		return nil, err
	}

	if nonce > 0 {
		nonce--
	}

	vmType := context.host.Runtime().GetVMType()
	return context.blockChainHook.NewAddress(creatorAddress, nonce, vmType)
}

func (context *blockchainContext) AccountExists(address []byte) bool {
	account, err := context.blockChainHook.GetUserAccount(address)
	if err != nil {
		return false
	}

	exists := !arwen.IfNil(account)
	return exists
}

func (context *blockchainContext) GetBalance(address []byte) []byte {
	return context.GetBalanceBigInt(address).Bytes()
}

func (context *blockchainContext) GetBalanceBigInt(address []byte) *big.Int {
	outputAccount, isNew := context.host.Output().GetOutputAccount(address)
	if !isNew {
		if outputAccount.Balance == nil {
			account, err := context.blockChainHook.GetUserAccount(address)
			if err != nil || arwen.IfNil(account) {
				return big.NewInt(0)
			}

			outputAccount.Balance = account.GetBalance()
		}

		balance := big.NewInt(0).Add(outputAccount.Balance, outputAccount.BalanceDelta)
		return balance
	}

	account, err := context.blockChainHook.GetUserAccount(address)
	if err != nil || arwen.IfNil(account) {
		return big.NewInt(0)
	}

	balance := account.GetBalance()
	outputAccount.Balance = balance

	return balance
}

func (context *blockchainContext) GetNonce(address []byte) (uint64, error) {
	// TODO verify if Nonce is 0, which means the outputAccount was cached, but
	// its Nonce not yet retrieved from the BlockchainHook; more generally,
	// create a list of accounts that have been cached, but not yet fully updated
	// from the BlockchainHook (they might have uninitialized Nonce and Balance).
	outputAccount, isNew := context.host.Output().GetOutputAccount(address)
	if !isNew {
		return outputAccount.Nonce, nil
	}

	account, err := context.blockChainHook.GetUserAccount(address)
	if err != nil || arwen.IfNil(account) {
		return 0, err
	}

	nonce := account.GetNonce()
	outputAccount.Nonce = nonce

	return nonce, nil
}

func (context *blockchainContext) IncreaseNonce(address []byte) {
	nonce, _ := context.GetNonce(address)
	outputAccount, _ := context.host.Output().GetOutputAccount(address)
	outputAccount.Nonce = nonce + 1
}

func (context *blockchainContext) GetCodeHash(addr []byte) ([]byte, error) {
	code, err := context.GetCode(addr)
	if err != nil {
		return nil, err
	}

	return context.host.Crypto().Keccak256(code)
}

func (context *blockchainContext) GetCode(address []byte) ([]byte, error) {
	account, err := context.blockChainHook.GetUserAccount(address)
	if err != nil {
		return nil, err
	}
	if arwen.IfNil(account) {
		return nil, arwen.ErrInvalidAccount
	}

	code := account.GetCode()
	if len(code) == 0 {
		return nil, arwen.ErrContractNotFound
	}

	return account.GetCode(), nil
}

func (context *blockchainContext) GetCodeSize(address []byte) (int32, error) {
	account, err := context.blockChainHook.GetUserAccount(address)
	if err != nil || arwen.IfNil(account) {
		return 0, err
	}

	code := account.GetCode()
	result := int32(len(code))
	return result, nil
}

func (context *blockchainContext) BlockHash(number int64) []byte {
	if number < 0 {
		return nil
	}

	block, err := context.blockChainHook.GetBlockhash(uint64(number))
	if err != nil {
		return nil
	}

	return block
}

func (context *blockchainContext) CurrentEpoch() uint32 {
	return context.blockChainHook.CurrentEpoch()
}

func (context *blockchainContext) CurrentNonce() uint64 {
	return context.blockChainHook.CurrentNonce()
}

func (context *blockchainContext) GetStateRootHash() []byte {
	return context.blockChainHook.GetStateRootHash()
}

func (context *blockchainContext) LastTimeStamp() uint64 {
	return context.blockChainHook.LastTimeStamp()
}

func (context *blockchainContext) LastNonce() uint64 {
	return context.blockChainHook.LastNonce()
}

func (context *blockchainContext) LastRound() uint64 {
	return context.blockChainHook.LastRound()
}

func (context *blockchainContext) LastEpoch() uint32 {
	return context.blockChainHook.LastEpoch()
}

func (context *blockchainContext) CurrentRound() uint64 {
	return context.blockChainHook.CurrentRound()
}

func (context *blockchainContext) CurrentTimeStamp() uint64 {
	return context.blockChainHook.CurrentTimeStamp()
}

func (context *blockchainContext) LastRandomSeed() []byte {
	return context.blockChainHook.LastRandomSeed()
}

func (context *blockchainContext) CurrentRandomSeed() []byte {
	return context.blockChainHook.CurrentRandomSeed()
}

func (context *blockchainContext) GetOwnerAddress() ([]byte, error) {
	scAddress := context.host.Runtime().GetSCAddress()
	scAccount, err := context.blockChainHook.GetUserAccount(scAddress)
	if err != nil || arwen.IfNil(scAccount) {
		return nil, err
	}

	return scAccount.GetOwnerAddress(), nil
}

func (context *blockchainContext) GetShardOfAddress(addr []byte) uint32 {
	return context.blockChainHook.GetShardOfAddress(addr)
}

func (context *blockchainContext) IsSmartContract(addr []byte) bool {
	return context.blockChainHook.IsSmartContract(addr)
}
