package messaging

import "bytes"

type Message interface {
	Command() []byte
	Serialize() []byte
	Parse(reader *bytes.Reader) Message
}
