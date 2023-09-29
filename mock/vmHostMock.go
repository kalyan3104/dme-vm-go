package mock

import (
	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/wasmer"
)

var _ arwen.VMHost = (*VmHostMock)(nil)

type VmHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     vmcommon.CryptoHook

	EthInput []byte

	BlockchainContext arwen.BlockchainContext
	RuntimeContext    arwen.RuntimeContext
	OutputContext     arwen.OutputContext
	MeteringContext   arwen.MeteringContext
	StorageContext    arwen.StorageContext
	BigIntContext     arwen.BigIntContext

	SCAPIMethods *wasmer.Imports
}

func (host *VmHostMock) Crypto() vmcommon.CryptoHook {
	return host.CryptoHook
}

func (host *VmHostMock) Blockchain() arwen.BlockchainContext {
	return host.BlockchainContext
}

func (host *VmHostMock) Runtime() arwen.RuntimeContext {
	return host.RuntimeContext
}

func (host *VmHostMock) Output() arwen.OutputContext {
	return host.OutputContext
}

func (host *VmHostMock) Metering() arwen.MeteringContext {
	return host.MeteringContext
}

func (host *VmHostMock) Storage() arwen.StorageContext {
	return host.StorageContext
}

func (host *VmHostMock) BigInt() arwen.BigIntContext {
	return host.BigIntContext
}

func (host *VmHostMock) CreateNewContract(input *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, nil
}

func (host *VmHostMock) ExecuteOnSameContext(input *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error) {
	return nil, nil
}

func (host *VmHostMock) ExecuteOnDestContext(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, error) {
	return nil, nil, nil
}

func (host *VmHostMock) EthereumCallData() []byte {
	return host.EthInput
}

func (host *VmHostMock) InitState() {
}

func (host *VmHostMock) PushState() {
}

func (host *VmHostMock) PopState() {
}

func (host *VmHostMock) ClearStateStack() {
}

func (host *VmHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

func (host *VmHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}
