package messaging

import (
	"bytes"
	"reflect"
)

type ReceiveMessageTypeOption func() reflect.Type

type VerackMessage struct {
}

func VerackMessageOption() ReceiveMessageTypeOption {
	return func() reflect.Type {
		return reflect.TypeOf((*VerackMessage)(nil))
	}
}

func NewVerackMessage() *VerackMessage {
	return new(VerackMessage)
}

func (*VerackMessage) Command() []byte {
	return []byte("v")
}

func (msg *VerackMessage) Serialize() []byte {
	return make([]byte, 0)
}

func (msg *VerackMessage) Parse(reader *bytes.Reader) Message {
	return msg
}
