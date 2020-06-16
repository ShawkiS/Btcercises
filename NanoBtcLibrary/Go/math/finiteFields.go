package math

import (
	"fmt"
	"math/big"

	utils "github.com/Btcercises/NanoBtcLibrary/Go/math/utils"
)

var prime = utils.GeneratePrimeValue()

type FieldElement struct {
	Num   *big.Int
	Prime *big.Int
}

func CreateFieldElement(num big.Int) *FieldElement {
	localPrime := &prime
	localNum := &num
	numCmpPrime := num.Cmp(localPrime)
	numbCmpZero := num.Cmp(big.NewInt(0))

	if (numCmpPrime >= 1) || (numbCmpZero == -1) {
		panic(fmt.Errorf("Num %v not in field range 0 to %d\n", num, prime))
	}

	return &FieldElement{localNum, localPrime}
}

func (f *FieldElement) IsEqual(fe *FieldElement) bool {

	if f.Num.Cmp(fe.Num) == 0 {
		return true
	}

	return false
}

func (f *FieldElement) Add(fe FieldElement) FieldElement {
	var res, mod big.Int

	if !utils.IsSameField(f.Prime, fe.Prime) {
		panic("Sorry, not the same field")
	}

	res.Add(f.Num, fe.Num)
	mod.Mod(&res, f.Prime)

	fld := FieldElement{
		Num:   &mod,
		Prime: f.Prime,
	}

	return fld
}

func (f *FieldElement) Sub(fe FieldElement) FieldElement {
	var res, mod big.Int

	if !utils.IsSameField(f.Prime, fe.Prime) {
		panic("Not members of same field")
	}

	res.Sub(f.Num, fe.Num)
	mod.Mod(&res, f.Prime)

	fld := FieldElement{
		Num:   &mod,
		Prime: f.Prime,
	}

	return fld
}

func (f *FieldElement) Mul(fe FieldElement) FieldElement {
	var res, mod big.Int

	if !utils.IsSameField(f.Prime, fe.Prime) {
		panic("Not member of same field")
	}

	res.Mul(f.Num, fe.Num)
	mod.Mod(&res, f.Prime)

	fld := FieldElement{
		Num:   &mod,
		Prime: f.Prime,
	}

	return fld
}

func (f *FieldElement) Pow(exp big.Int) FieldElement {

	var num, n, fPrime big.Int
	fPrime.Sub(f.Prime, big.NewInt(1))
	n.Mod(&exp, &fPrime)

	num.Exp(f.Num, &n, f.Prime)

	fld := FieldElement{
		Num:   &num,
		Prime: f.Prime,
	}

	return fld
}

func (f *FieldElement) Div(fe FieldElement) FieldElement {

	if !utils.IsSameField(f.Prime, fe.Prime) {
		panic("Not member of same field")
	}

	var prime2, res, mod big.Int

	prime2.Sub(f.Prime, big.NewInt(2))

	res.Exp(fe.Num, &prime2, f.Prime)
	res.Mul(&res, f.Num)

	mod.Mod(&res, f.Prime)

	fld := FieldElement{
		Num:   &mod,
		Prime: f.Prime,
	}

	return fld
}
