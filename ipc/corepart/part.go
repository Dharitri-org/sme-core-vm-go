package corepart

import (
	"os"
	"time"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/core/host"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/marshaling"
	logger "github.com/Dharitri-org/sme-logger"
	vmcommon "github.com/Dharitri-org/sme-vm-common"
)

var log = logger.GetOrCreate("core/part")

// CorePart is the endpoint that implements the message loop on Core's side
type CorePart struct {
	Messenger *CoreMessenger
	VMHost    vmcommon.VMExecutionHandler
	Repliers  []common.MessageReplier
}

// NewCorePart creates the Core part
func NewCorePart(
	input *os.File,
	output *os.File,
	vmHostParameters *core.VMHostParameters,
	marshalizer marshaling.Marshalizer,
) (*CorePart, error) {
	messenger := NewCoreMessenger(input, output, marshalizer)
	blockchain := NewBlockchainHookGateway(messenger)
	crypto := NewCryptoHookGateway()

	newCoreHost, err := host.NewCoreVM(
		blockchain,
		crypto,
		vmHostParameters,
	)
	if err != nil {
		return nil, err
	}

	part := &CorePart{
		Messenger: messenger,
		VMHost:    newCoreHost,
	}

	part.Repliers = common.CreateReplySlots(part.noopReplier)
	part.Repliers[common.ContractDeployRequest] = part.replyToRunSmartContractCreate
	part.Repliers[common.ContractCallRequest] = part.replyToRunSmartContractCall
	part.Repliers[common.DiagnoseWaitRequest] = part.replyToDiagnoseWait

	return part, nil
}

func (part *CorePart) noopReplier(_ common.MessageHandler) common.MessageHandler {
	log.Error("noopReplier called")
	return common.CreateMessage(common.UndefinedRequestOrResponse)
}

// StartLoop runs the main loop
func (part *CorePart) StartLoop() error {
	part.Messenger.Reset()
	err := part.doLoop()
	part.Messenger.Shutdown()
	log.Error("end of loop", "err", err)
	return err
}

// doLoop ends only when a critical failure takes place
func (part *CorePart) doLoop() error {
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

func (part *CorePart) replyToNodeRequest(request common.MessageHandler) common.MessageHandler {
	replier := part.Repliers[request.GetKind()]
	return replier(request)
}

func (part *CorePart) replyToRunSmartContractCreate(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractDeployRequest)
	vmOutput, err := part.VMHost.RunSmartContractCreate(typedRequest.CreateInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *CorePart) replyToRunSmartContractCall(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageContractCallRequest)
	vmOutput, err := part.VMHost.RunSmartContractCall(typedRequest.CallInput)
	return common.NewMessageContractResponse(vmOutput, err)
}

func (part *CorePart) replyToDiagnoseWait(request common.MessageHandler) common.MessageHandler {
	typedRequest := request.(*common.MessageDiagnoseWaitRequest)
	duration := time.Duration(int64(typedRequest.Milliseconds) * int64(time.Millisecond))
	time.Sleep(duration)
	return common.NewMessageDiagnoseWaitResponse()
}
