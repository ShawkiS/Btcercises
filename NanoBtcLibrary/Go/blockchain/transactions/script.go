package transactions

import (
	cryptoUtils "github.com/Btcercises/NanoBtcLibrary/Go/crypto/utils"
	utils "github.com/Btcercises/NanoBtcLibrary/Go/utils"
)

var OP_CODE_FUNCTIONS = []OP_CODE_FUNCTION{{118, "Op_dup", "0x76"}, {169, "Op_hash160", "0xa9"}}

type Script []byte

type OP_CODE_FUNCTION struct {
	code     int
	function string
	hex      string
}

func Op_dup(stack utils.Stack) bool {
	if stack.Len() < 1 {
		return false
	}
	stack.Push(stack.Peek())
	return true
}

func Op_hash160(stack utils.Stack) bool {
	if stack.Len() < 1 {
		return false
	}
	element := stack.Pop()
	h160 := cryptoUtils.Hash160([]byte(element))
	stack.Push(string(h160))
	return true
}
