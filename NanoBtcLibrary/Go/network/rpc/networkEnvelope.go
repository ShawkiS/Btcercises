package rpc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Btcercises/NanoBtcLibrary/Go/network/util"
)

var (
	networkMagic     = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	testNetworkMagic = [4]byte{0x0b, 0x11, 0x09, 0x07}
)

type NetworkEnvelope struct {
	Command     []byte
	Payload     []byte
	Magic       [4]byte
	TestNetwork bool
}

func NewEnvelope(command []byte, payload []byte, testnet bool) *NetworkEnvelope {
	var magic [4]byte
	if testnet {
		magic = testNetworkMagic
	} else {
		magic = networkMagic
	}
	return &NetworkEnvelope{Command: command, Payload: payload, Magic: magic}
}

func ParseEnvelope(reader *bytes.Reader, testnet bool) *NetworkEnvelope {
	if reader.Len() == 0 {
		panic("Connection reset!")
	}
	magic := make([]byte, 4)
	reader.Read(magic)
	var expectedMagic [4]byte
	if testnet {
		expectedMagic = testNetworkMagic
	} else {
		expectedMagic = networkMagic
	}
	if !bytes.Equal(magic[:], expectedMagic[:]) {
		panic(fmt.Sprintf("magic is not right %v vs %v", hex.EncodeToString(magic[:]), hex.EncodeToString(expectedMagic[:])))
	}
	command := make([]byte, 12)
	reader.Read(command)
	command = bytes.TrimRightFunc(command, func(r rune) bool { return r == 0 })
	buffer := make([]byte, 4)
	reader.Read(buffer)
	payloadLength := util.LittleEndianToInt32(buffer)
	checksum := make([]byte, 4)
	reader.Read(checksum)
	payload := make([]byte, payloadLength)
	reader.Read(payload)
	if !bytes.Equal(util.Hash256(payload)[:4], checksum) {
		fmt.Fprintf(os.Stderr, "%x %x\n", util.Hash256(payload)[:4], checksum)
		panic("Invalid checksum!")
	}
	var returnedMagic [4]byte
	return &NetworkEnvelope{command, payload, returnedMagic, testnet}
}

func (env *NetworkEnvelope) Serialize() []byte {
	command := make([]byte, 12)
	copy(command, env.Command)
	payloadLength := len(env.Payload)
	checksum := util.Hash256(env.Payload)[:4]
	result := make([]byte, payloadLength+24)
	copy(result[:4], env.Magic[:])
	copy(result[4:16], command)
	copy(result[16:20], util.Int32ToLittleEndian(uint32(payloadLength)))
	copy(result[20:24], checksum)
	copy(result[24:], env.Payload)
	return result
}
func (env *NetworkEnvelope) Stream() *bytes.Reader {
	return bytes.NewReader(env.Payload)
}
