package crypto

import (
	"math/big"

	"github.com/Btcercises/NanoBtcLibrary/math"
	"github.com/Btcercises/NanoBtcLibrary/math/utils"
)

var prime = utils.GeneratePrimeValue()

func CheckIfOnCurve(x *big.Int, y *big.Int) bool {

	var y2, x3, reEq, y2Mod, reEqMod big.Int
	var e2 = big.NewInt(2)
	var e3 = big.NewInt(3)

	y2.Exp(y, e2, nil)
	y2Mod.Mod(&y2, &prime)
	x3.Exp(x, e3, nil)
	reEq.Mul(big.NewInt(0), x)
	reEq.Add(&reEq, &x3)
	reEq.Add(&reEq, big.NewInt(7))
	reEqMod.Mod(&reEq, &prime)

	res := y2Mod.Cmp(&reEqMod)
	if res != 0 {
		return false
	}

	return true
}

func IsSameCurve(p1 math.FFPoint, p2 math.FFPoint) bool {
	p1a := p1.A
	p1b := p1.B
	p2a := p2.A
	p2b := p2.B

	if (!p1a.IsEqual(p2a)) || (!p1b.IsEqual(p2b)) {
		return false
	}

	return true
}
