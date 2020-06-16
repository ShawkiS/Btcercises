package utils

import (
	"encoding/hex"
	"math/big"
)

// p = 2^256 – 232 – 977
func GeneratePrimeValue() big.Int {
	var two256, two32, p big.Int
	two256.Exp(big.NewInt(2), big.NewInt(256), nil)
	two32.Exp(big.NewInt(2), big.NewInt(32), nil)

	p.Sub(&two256, &two32)
	p.Sub(&p, big.NewInt(977))

	return p
}

func IsSameField(num1 *big.Int, num2 *big.Int) bool {
	if num1.Cmp(num2) == 0 {
		return true
	}
	return false
}

func HexToBigInt(str string) *big.Int {

	hexStr := str[2:len(str)]

	decByte, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	z := new(big.Int)
	z.SetBytes(decByte)

	return z
}
