package comms

import "context"

type MessageHandler func(message string)

type Communicator interface {
	Start(ctx context.Context, handler MessageHandler) error
	Send(message string) error
	Close() error
}

func NewCommunicator() Communicator {
	return newPlatformSpecificCommunicator()
}
