package arwenpart

import (
	"os"
	"time"

	logger "github.com/kalyan3104/dme-logger-go"
	vmcommon "github.com/kalyan3104/dme-vm-common"
	"github.com/kalyan3104/dme-vm-go/arwen"
	"github.com/kalyan3104/dme-vm-go/arwen/host"
	"github.com/kalyan3104/dme-vm-go/ipc/common"
	"github.com/kalyan3104/dme-vm-go/ipc/marshaling"
)

var log = logger.GetOrCreate("arwen/part")

// ArwenPart is the endpoint that implements the message loop on Arwen's side
type ArwenPart struct {
	Messenger *ArwenMessenger
	VMHost    vmcommon.VMExecutionHandler
	Repliers  []common.MessageReplier
}

// NewArwenPart creates the Arwen part
func NewArwenPart(
	input *os.File,
	output *os.File,
	vmHostParameters *arwen.VMHostParameters,
	marshalizer marshaling.Marshalizer,
) (*ArwenPart, error) {
	messenger := NewArwenMessenger(input, output, marshalizer)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()

	newArwenHost, err := host.NewArwenVM(
		blockchain,
		crypto,
		vmHostParameters,
	)
	if err != nil {
		return nil, err
	}

	part := &ArwenPart{
		Messenger: messenger,
		VMHost:    newArwenHost,
	}

	part.Repliers = common.CreateReplySlots(part.noopReplier)
	part.Repliers[common.ContractDeployRequest] = part.replyToRunSmartContractCreate
	part.Repliers[common.ContractCallRequest] = part.replyToRunSmartContractCall
	part.Repliers[common.DiagnoseWaitRequest] = part.replyToDiagnoseWait

	return part, nil
}

func (part *ArwenPart) noopReplier(_ common.MessageHandler) common.MessageHandler {
	log.Error("noopReplier called")
	return common.CreateMessage(common.UndefinedRequestOrResponse)
}

// StartLoop runs the main loop
func (part *ArwenPart) StartLoop() error {
	part.Messenger.Reset()
	err := part.doLoop()
	part.Messenger.Shutdown()
	log.Error("end of loop", "err", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *ArwenPart) doLoop() error {
	for {
		request, err := part.Messenger.ReceiveNodeRequest()
		if err != nil {
			return err
		}
		if common.IsStopRequest(request) {
			return common.ErrStopPerNodeRequest
		}

		response := part.replyToNodeRequest(request)

		// Successful execution, send response
		err = part.Messenger.SendContractResponse(response)
		if err != nil {
			return err
		}

		part.Messenger.ResetDialogue()
	}
}

func (part *ArwenPart) replyToNodeRequest(request common.MessageHandler) common.MessageHandler {
	replier := part.Repliers[request.GetKind()]
	return replier(request)
}

func (part *ArwenPart) replyToRunSmartContractCreate(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractDeployRequest)
	vmOutput, err := part.VMHost.RunSmartContractCreate(typedRequest.CreateInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *ArwenPart) replyToRunSmartContractCall(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractCallRequest)
	vmOutput, err := part.VMHost.RunSmartContractCall(typedRequest.CallInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *ArwenPart) replyToDiagnoseWait(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageDiagnoseWaitRequest)
	duration := time.Duration(int64(typedRequest.Milliseconds) * int64(time.Millisecond))
	time.Sleep(duration)
	return common.NewMessageDiagnoseWaitResponse()
}
