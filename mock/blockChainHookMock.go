package mock

import (
	"errors"
	"math/big"

	vmcommon "github.com/kalyan3104/dme-vm-common"
)

var ErrAccountDoesntExist = errors.New("account does not exist")

var _ vmcommon.BlockchainHook = (*BlockchainHookMock)(nil)

// AccountMap is a map from address to account
type AccountsMap map[string]*AccountMock

type BlockchainHookMock struct {
	Accounts      AccountsMap
	BlockHash     []byte
	LNonce        uint64
	LRound        uint64
	CNonce        uint64
	CRound        uint64
	LTimeStamp    uint64
	CTimeStamp    uint64
	LRandomSeed   []byte
	CRandomSeed   []byte
	LEpoch        uint32
	CEpoch        uint32
	StateRootHash []byte
	NewAddr       []byte
	Value         *big.Int
	Gas           uint64
	Err           error
}

func NewBlockchainHookMock() *BlockchainHookMock {
	return &BlockchainHookMock{
		Accounts: make(AccountsMap),
	}
}

func (b *BlockchainHookMock) AddAccount(account *AccountMock) {
	if account.Storage == nil {
		account.Storage = make(map[string][]byte)
	}
	if account.Balance == nil {
		account.Balance = big.NewInt(0)
	}
	b.Accounts[string(account.Address)] = account
}

func (b *BlockchainHookMock) AddAccounts(accounts []*AccountMock) {
	for _, account := range accounts {
		b.AddAccount(account)
	}
}

func (b *BlockchainHookMock) NewAddress(creatorAddress []byte, creatorNonce uint64, vmType []byte) ([]byte, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	return b.NewAddr, nil
}

func (b *BlockchainHookMock) GetStorageData(address []byte, index []byte) ([]byte, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	account, ok := b.Accounts[string(address)]
	if !ok {
		return []byte{}, ErrAccountDoesntExist
	}

	if account.Err != nil {
		return nil, account.Err
	}

	return account.Storage[string(index)], nil
}

func (b *BlockchainHookMock) GetBlockhash(nonce uint64) ([]byte, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	return b.BlockHash, nil
}

func (b *BlockchainHookMock) LastNonce() uint64 {
	return b.LNonce
}

func (b *BlockchainHookMock) LastRound() uint64 {
	return b.LRound
}

func (b *BlockchainHookMock) LastTimeStamp() uint64 {
	return b.LTimeStamp
}

func (b *BlockchainHookMock) LastRandomSeed() []byte {
	return b.LRandomSeed
}

func (b *BlockchainHookMock) LastEpoch() uint32 {
	return b.LEpoch
}

func (b *BlockchainHookMock) GetStateRootHash() []byte {
	return b.StateRootHash
}

func (b *BlockchainHookMock) CurrentNonce() uint64 {
	return b.CNonce
}

func (b *BlockchainHookMock) CurrentRound() uint64 {
	return b.CRound
}

func (b *BlockchainHookMock) CurrentTimeStamp() uint64 {
	return b.CTimeStamp
}

func (b *BlockchainHookMock) CurrentRandomSeed() []byte {
	return b.CRandomSeed
}

func (b *BlockchainHookMock) CurrentEpoch() uint32 {
	return b.CEpoch
}

func (b *BlockchainHookMock) ProcessBuiltInFunction(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	outPutAccounts := make(map[string]*vmcommon.OutputAccount)
	outPutAccounts[string(input.CallerAddr)] = &vmcommon.OutputAccount{BalanceDelta: b.Value}

	return &vmcommon.VMOutput{
		GasRemaining:   b.Gas,
		OutputAccounts: outPutAccounts,
	}, b.Err
}

func (b *BlockchainHookMock) GetBuiltinFunctionNames() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

func (b *BlockchainHookMock) GetAllState(_ []byte) (map[string][]byte, error) {
	return nil, nil
}

func (b *BlockchainHookMock) GetUserAccount(address []byte) (vmcommon.UserAccountHandler, error) {
	if b.Err != nil {
		return nil, b.Err
	}

	account, ok := b.Accounts[string(address)]
	if !ok {
		return nil, ErrAccountDoesntExist
	}

	if account.Err != nil {
		return nil, account.Err
	}

	return account, nil
}

func (b *BlockchainHookMock) GetShardOfAddress(address []byte) uint32 {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return 0
	}

	return account.ShardID
}

func (b *BlockchainHookMock) IsSmartContract(address []byte) bool {
	account, ok := b.Accounts[string(address)]
	if !ok {
		return false
	}

	return len(account.Code) > 0
}

func (b *BlockchainHookMock) UpdateAccounts(outputAccounts map[string]*vmcommon.OutputAccount) {
	for strAddress, outputAccount := range outputAccounts {
		account, exists := b.Accounts[strAddress]
		if !exists {
			account = &AccountMock{
				Address: outputAccount.Address,
				Balance: big.NewInt(0),
				Code:    nil,
				Storage: make(map[string][]byte),
				Nonce:   0,
			}
		}

		if outputAccount.Nonce > account.Nonce {
			account.Nonce = outputAccount.Nonce
		}
		account.Balance.Add(account.Balance, outputAccount.BalanceDelta)
		if len(outputAccount.Code) > 0 {
			account.Code = outputAccount.Code
		}

		mergeStorageUpdates(account, outputAccount)
		b.Accounts[strAddress] = account
	}
}

func mergeStorageUpdates(
	leftAccount *AccountMock,
	rightAccount *vmcommon.OutputAccount,
) {
	if leftAccount.Storage == nil {
		leftAccount.Storage = make(map[string][]byte)
	}
	for key, update := range rightAccount.StorageUpdates {
		leftAccount.Storage[key] = update.Data
	}
}
