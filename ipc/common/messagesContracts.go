package common

import (
	vmcommon "github.com/kalyan3104/dme-vm-common"
)

// MessageContractDeployRequest is deploy request message (from Node)
type MessageContractDeployRequest struct {
	Message
	CreateInput *vmcommon.ContractCreateInput
}

// NewMessageContractDeployRequest creates a message
func NewMessageContractDeployRequest(input *vmcommon.ContractCreateInput) *MessageContractDeployRequest {
	message := &MessageContractDeployRequest{}
	message.Kind = ContractDeployRequest
	message.CreateInput = input
	return message
}

// MessageContractCallRequest is call request message (from Node)
type MessageContractCallRequest struct {
	Message
	CallInput *vmcommon.ContractCallInput
}

// NewMessageContractCallRequest creates a message
func NewMessageContractCallRequest(input *vmcommon.ContractCallInput) *MessageContractCallRequest {
	message := &MessageContractCallRequest{}
	message.Kind = ContractCallRequest
	message.CallInput = input
	return message
}

// MessageContractResponse is a contract response message (from Arwen)
type MessageContractResponse struct {
	Message
	VMOutput *vmcommon.VMOutput
}

// NewMessageContractResponse creates a message
func NewMessageContractResponse(vmOutput *vmcommon.VMOutput, err error) *MessageContractResponse {
	message := &MessageContractResponse{}
	message.Kind = ContractResponse
	message.VMOutput = vmOutput
	message.SetError(err)
	return message
}
