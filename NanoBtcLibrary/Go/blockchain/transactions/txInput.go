package transactions

import (
	"encoding/binary"

	"github.com/Btcercises/NanoBtcLibrary/Go/blockchain/utils"
)

var SequenceDefaultVal = "0xffffffff"

type TxInput struct {
	Hash          []byte
	Index         uint32
	PrevIndex     int
	PrevTx        Transaction
	Script        Script
	Sequence      uint32
	ScriptWitness [][]byte
	Value         int
}

func (in TxInput) Binary() []byte {
	bin := make([]byte, 0)
	bin = append(bin, in.Hash...)

	index := make([]byte, 4)
	binary.LittleEndian.PutUint32(index, uint32(in.Index))
	bin = append(bin, index...)

	scriptLength := utils.Varint(uint64(len(in.Script)))
	bin = append(bin, scriptLength...)
	bin = append(bin, in.Script...)

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, uint32(in.Sequence))
	bin = append(bin, sequence...)

	return bin
}
