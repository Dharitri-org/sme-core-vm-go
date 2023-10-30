package common

import (
	"os"

	"github.com/Dharitri-org/sme-core-vm-go/core"
	"github.com/Dharitri-org/sme-core-vm-go/ipc/marshaling"
)

// CoreArguments represents the initialization arguments required by Core, passed through the initialization pipe
type CoreArguments struct {
	core.VMHostParameters
	LogsMarshalizer     marshaling.MarshalizerKind
	MessagesMarshalizer marshaling.MarshalizerKind
}

// SendCoreArguments sends initialization arguments through a pipe
func SendCoreArguments(pipe *os.File, pipeArguments CoreArguments) error {
	sender := NewSender(pipe, createArgumentsMarshalizer())
	message := NewMessageInitialize(pipeArguments)
	_, err := sender.Send(message)
	return err
}

// GetCoreArguments reads initialization arguments from the pipe
func GetCoreArguments(pipe *os.File) (*CoreArguments, error) {
	receiver := NewReceiver(pipe, createArgumentsMarshalizer())
	message, _, err := receiver.Receive(0)
	if err != nil {
		return nil, err
	}

	typedMessage := message.(*MessageInitialize)
	return &typedMessage.Arguments, nil
}

// For the arguments, the marshalizer is fixed to JSON
func createArgumentsMarshalizer() marshaling.Marshalizer {
	return marshaling.CreateMarshalizer(marshaling.JSON)
}
