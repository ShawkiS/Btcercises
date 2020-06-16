package transactions

import (
	"encoding/binary"

	"github.com/Btcercises/NanoBtcLibrary/blockchain/utils"
)

type TxOutput struct {
	Amount int64
	Script Script
}

func (out TxOutput) Binary() Script {
	bin := make([]byte, 0)

	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(out.Amount))
	bin = append(bin, value...)

	scriptLength := utils.Varint(uint64(len(out.Script)))
	bin = append(bin, scriptLength...)
	bin = append(bin, out.Script...)

	return bin
}
