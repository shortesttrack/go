package pubsub

import (
	"context"
)

type Publisher interface {
	Publish(context.Context, Message) error
}

type Subscriber interface {
	Start() <-chan Message
	Err() error
	Stop()
}

type Message interface {
	ID() string
	Attributes() map[string]string
	Event() string
	Data() []byte

	Ack()
	Nack()
}
