package transactions

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/Btcercises/NanoBtcLibrary/blockchain/utils"
)

type Transaction struct {
	Id       []byte
	Version  int32
	Input    []TxInput
	Output   []TxOutput
	Locktime uint32
	Testnet  bool
}

func Print(tx Transaction) {
	fmt.Printf("tx: %v\nversion: %v\ntx_ins:\n%vtx_outs:\n%vlocktime: %v",
		tx.Id,
		tx.Version,
		tx.Input,
		tx.Output,
		tx.Locktime)
}

func CreateTransaction(version int32, input []TxInput, output []TxOutput, locktime uint32, testnet bool) Transaction {
	tx := Transaction{nil, version, input, output, locktime, testnet}
	tx.Id = GenerateTransactionId(tx)
	return tx
}

func GenerateTransactionId(tx Transaction) []byte {
	if tx.Id != nil {
		return tx.Id
	}
	bin := make([]byte, 0)
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(tx.Version))
	bin = append(bin, version...)

	vinLength := utils.Varint(uint64(len(tx.Input)))
	bin = append(bin, vinLength...)
	for _, in := range tx.Input {
		bin = append(bin, in.Binary()...)
	}

	voutLength := utils.Varint(uint64(len(tx.Output)))
	bin = append(bin, voutLength...)
	for _, out := range tx.Output {
		bin = append(bin, out.Binary()...)
	}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.Locktime)
	bin = append(bin, locktime...)

	tx.Id = utils.DoubleSha256(bin)
	return tx.Id
}

func IsCoinbaseTx(tx Transaction) bool {
	if len(tx.Input) != 1 {
		return false
	}

	firstInput := tx.Input[0]

	if string(firstInput.PrevTx.Id) != strings.Repeat("00", 32) {
		return false
	}

	if firstInput.PrevIndex != -1 {
		return false
	}

	return true
}

func CreateCoinbaseTransaction(tx Transaction) {
	isCoinbaseTx := IsCoinbaseTx(tx)

	if !isCoinbaseTx {
		panic("Duplicated Coinbase Transaction")
	}

	// element = tx.Input[0].script_sig.cmds[0]
	// return little_endian_to_int(element)
}

func GetUrl(testnet bool) string {
	if testnet {
		return "http://testnet.programmingbitcoin.com"
	}
	return "http://mainnet.programmingbitcoin.com/"
}
