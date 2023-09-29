package nodepart

import (
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/ipc/common"
)

func (part *NodePart) replyToBlockchainNewAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainNewAddressRequest)
	result, err := part.blockchain.NewAddress(typedRequest.CreatorAddress, typedRequest.CreatorNonce, typedRequest.VmType)
	response := common.NewMessageBlockchainNewAddressResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainGetStorageData(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetStorageDataRequest)
	data, err := part.blockchain.GetStorageData(typedRequest.Address, typedRequest.Index)
	response := common.NewMessageBlockchainGetStorageDataResponse(data, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBlockhash(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetBlockhashRequest)
	result, err := part.blockchain.GetBlockhash(typedRequest.Nonce)
	response := common.NewMessageBlockchainGetBlockhashResponse(result, err)
	return response
}

func (part *NodePart) replyToBlockchainLastNonce(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastNonce()
	response := common.NewMessageBlockchainLastNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRound(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRound()
	response := common.NewMessageBlockchainLastRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastTimeStamp(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastTimeStamp()
	response := common.NewMessageBlockchainLastTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastRandomSeed(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastRandomSeed()
	response := common.NewMessageBlockchainLastRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainLastEpoch(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.LastEpoch()
	response := common.NewMessageBlockchainLastEpochResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainGetStateRootHash(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.GetStateRootHash()
	response := common.NewMessageBlockchainGetStateRootHashResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentNonce(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentNonce()
	response := common.NewMessageBlockchainCurrentNonceResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRound(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRound()
	response := common.NewMessageBlockchainCurrentRoundResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentTimeStamp(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentTimeStamp()
	response := common.NewMessageBlockchainCurrentTimeStampResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentRandomSeed(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentRandomSeed()
	response := common.NewMessageBlockchainCurrentRandomSeedResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainCurrentEpoch(request common.MessageHandler) common.MessageHandler {
	result := part.blockchain.CurrentEpoch()
	response := common.NewMessageBlockchainCurrentEpochResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainProcessBuiltinFunction(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainProcessBuiltinFunctionRequest)
	vmOutput, err := part.blockchain.ProcessBuiltInFunction(&typedRequest.CallInput)
	response := common.NewMessageBlockchainProcessBuiltinFunctionResponse(vmOutput, err)
	return response
}

func (part *NodePart) replyToBlockchainGetBuiltinFunctionNames(request common.MessageHandler) common.MessageHandler {
	functionNames := part.blockchain.GetBuiltinFunctionNames()
	response := common.NewMessageBlockchainGetBuiltinFunctionNamesResponse(functionNames)
	return response
}

func (part *NodePart) replyToBlockchainGetAllState(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetAllStateRequest)
	state, err := part.blockchain.GetAllState(typedRequest.Address)
	response := common.NewMessageBlockchainGetAllStateResponse(state, err)
	return response
}

func (part *NodePart) replyToBlockchainGetUserAccount(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetUserAccountRequest)
	account, err := part.blockchain.GetUserAccount(typedRequest.Address)

	if arwen.IfNil(account) {
		return common.NewMessageBlockchainGetUserAccountResponse(nil, err)
	}

	return common.NewMessageBlockchainGetUserAccountResponse(&common.Account{
		Nonce:           account.GetNonce(),
		Address:         account.AddressBytes(),
		Balance:         account.GetBalance(),
		Code:            account.GetCode(),
		CodeMetadata:    account.GetCodeMetadata(),
		CodeHash:        account.GetCodeHash(),
		RootHash:        account.GetRootHash(),
		DeveloperReward: account.GetDeveloperReward(),
		OwnerAddress:    account.GetOwnerAddress(),
		UserName:        account.GetUserName(),
	}, err)
}

func (part *NodePart) replyToBlockchainGetShardOfAddress(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainGetShardOfAddressRequest)
	result := part.blockchain.GetShardOfAddress(typedRequest.Address)
	response := common.NewMessageBlockchainGetShardOfAddressResponse(result)
	return response
}

func (part *NodePart) replyToBlockchainIsSmartContract(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageBlockchainIsSmartContractRequest)
	result := part.blockchain.IsSmartContract(typedRequest.Address)
	response := common.NewMessageBlockchainIsSmartContractResponse(result)
	return response
}
