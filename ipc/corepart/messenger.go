package corepart

import (
	"os"

	"github.com/Dharitri-org/sme-core-vm-go/ipc/common"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/marshaling"
)

// CoreMessenger is the messenger on Core's part of the pipe
type CoreMessenger struct {
	common.Messenger
}

// NewCoreMessenger creates a new messenger
func NewCoreMessenger(reader *os.File, writer *os.File, marshalizer marshaling.Marshalizer) *CoreMessenger {
	return &CoreMessenger{
		Messenger: *common.NewMessengerPipes("CORE", reader, writer, marshalizer),
	}
}

// ReceiveNodeRequest waits for a request from Node
func (messenger *CoreMessenger) ReceiveNodeRequest() (common.MessageHandler, error) {
	message, err := messenger.Receive(0)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// SendContractResponse sends a contract response to the Node
func (messenger *CoreMessenger) SendContractResponse(response common.MessageHandler) error {
	log.Trace("[CORE]: SendContractResponse", "response", response.DebugString())

	err := messenger.Send(response)
	if err != nil {
		return err
	}

	return nil
}

// SendHookCallRequest makes a hook call (over the pipe) and waits for the response
func (messenger *CoreMessenger) SendHookCallRequest(request common.MessageHandler) (common.MessageHandler, error) {
	log.Trace("[CORE]: SendHookCallRequest", "request", request.DebugString())

	err := messenger.Send(request)
	if err != nil {
		return nil, common.ErrCannotSendHookCallRequest
	}

	response, err := messenger.Receive(0)
	if err != nil {
		return nil, common.ErrCannotReceiveHookCallResponse
	}

	return response, nil
}
