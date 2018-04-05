package pubsub

import "cloud.google.com/go/pubsub"

const (
	EventTypeKey = "event_type"
)

func NewMessage(event string, data []byte) Message {
	attr := make(map[string]string)
	attr[EventTypeKey] = event
	pm := &pubsub.Message{
		Data: data,
		Attributes: attr,
	}
	return &message{pm}
}

type message struct {
	*pubsub.Message
}

func (m *message) ID() string {
	return m.Message.ID
}

func (m *message) Attributes() map[string]string {
	return m.Message.Attributes
}

func (m *message) Event() string {
	return m.Message.Attributes[EventTypeKey]
}

func (m *message) Data() []byte {
	return m.Message.Data
}
