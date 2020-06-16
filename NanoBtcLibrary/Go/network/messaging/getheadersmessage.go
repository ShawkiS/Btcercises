package messaging

import (
	"bytes"

	"github.com/Btcercises/NanoBtcLibrary/Go/network/util"
)

type GetHeadersMessage struct {
	Version    uint32
	NumHashes  int
	StartBlock [32]byte
	EndBlock   [32]byte
}

const (
	defaultVersion   uint32 = 70015
	defaultNumHashes int    = 1
)

func NewGetHeadersMessage(startBlock []byte) *GetHeadersMessage {
	var startBlockData [32]byte
	var endBlockData [32]byte
	copy(startBlockData[:], startBlock[:32])
	return &GetHeadersMessage{
		Version:    defaultVersion,
		NumHashes:  defaultNumHashes,
		StartBlock: startBlockData,
		EndBlock:   endBlockData,
	}
}

func (*GetHeadersMessage) Command() []byte {
	return []byte("getheaders")
}

func (msg *GetHeadersMessage) Serialize() []byte {
	version := util.Int32ToLittleEndian(msg.Version)
	numHashes := util.EncodeVarInt(msg.NumHashes)
	startBlock := make([]byte, 32)
	endBlock := make([]byte, 32)
	copy(startBlock, msg.StartBlock[:])
	copy(endBlock, msg.EndBlock[:])
	util.ReverseByteArray(startBlock)
	util.ReverseByteArray(endBlock)
	result := make([]byte, 68+len(numHashes))
	copy(result[:4], version)
	copy(result[4:4+len(numHashes)], numHashes)
	copy(result[4+len(numHashes):36+len(numHashes)], startBlock)
	copy(result[36+len(numHashes):], endBlock)
	return result
}

func (msg *GetHeadersMessage) Parse(reader *bytes.Reader) Message {
	version := make([]byte, 4)
	reader.Read(version)
	msg.Version = util.LittleEndianToInt32(version)
	msg.NumHashes = util.ReadVarInt(reader)
	blockData := make([]byte, 32)
	reader.Read(blockData)
	copy(msg.StartBlock[:], util.ReverseByteArray(blockData))
	reader.Read(blockData)
	copy(msg.EndBlock[:], util.ReverseByteArray(blockData))
	return msg
}
