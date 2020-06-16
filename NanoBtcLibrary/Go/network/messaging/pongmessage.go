package messaging

import (
	"bytes"
)

type PongMessage struct {
	Nonce [8]byte
}

func NewPongMessage(nonce []byte) *PongMessage {
	result := new(PongMessage)
	copy(result.Nonce[:], nonce)
	return result
}

func (*PongMessage) Command() []byte {
	return []byte("pong")
}

func (msg *PongMessage) Serialize() []byte {
	return msg.Nonce[:]
}

func (msg *PongMessage) Parse(reader *bytes.Reader) Message {
	reader.Read(msg.Nonce[:])
	return msg
}
