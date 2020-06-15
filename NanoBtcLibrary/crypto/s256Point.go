package crypto

import (
	"fmt"
	"math/big"

	btcMath "github.com/Btcercises/NanoBtcLibrary/math"
	"github.com/Btcercises/NanoBtcLibrary/math/utils"

	cryptoUtils "github.com/Btcercises/NanoBtcLibrary/crypto/utils"
	"github.com/Btcercises/NanoBtcLibrary/math"
	base58 "github.com/btcsuite/btcutil/base58"
)

// S256Point struct representation of s256 point
type S256Point struct {
	point *math.FFPoint
}

const N = "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"

const Gx = "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"

const Gy = "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

func NewS256Point(x, y big.Int) S256Point {

	cmp := x.Cmp(big.NewInt(0))
	cmpZ := y.Cmp(big.NewInt(0))

	if cmp == 0 || cmpZ == 0 {
		NewP := btcMath.CreateNewPoint(*big.NewInt(0), *big.NewInt(0))
		newSP := S256Point{&NewP}
		return newSP
	}

	NewP := btcMath.CreateNewPoint(x, y)
	newSP := S256Point{&NewP}
	return newSP
}

func (sp *S256Point) S256RMul(coef big.Int) *S256Point {

	decByte := utils.HexToBigInt(N)

	var cf big.Int
	cf.Mod(&coef, decByte)

	res := sp.point.RMul(cf)

	r256 := S256Point{res}
	return &r256
}

func (sp *S256Point) verify(z, s, r *big.Int) bool {

	var sInv, u, v, zsInv, rsInv big.Int

	n := utils.HexToBigInt(N)
	nField := btcMath.CreateFieldElement(*n)

	nMinTwo := nField.Sub(*btcMath.CreateFieldElement(*big.NewInt(2)))

	sInv.Exp(s, nMinTwo.Num, n)

	zsInv.Mul(&sInv, z)
	u.Mod(&zsInv, n)

	rsInv.Mul(r, &sInv)

	v.Mod(&rsInv, n)

	G := gValue()

	total2 := G.S256RMul(u)
	total1 := sp.S256RMul(v)

	res := total2.point.Add(total1.point)

	return res.X.Num.Cmp(r) == 0

}

func gValue() *S256Point {
	xHex := utils.HexToBigInt(Gx)
	yHex := utils.HexToBigInt(Gy)

	G := NewS256Point(*xHex, *yHex)
	return &G
}

func Sec(p S256Point, compressed bool) string {
	if compressed == true {
		if new(big.Int).Div(p.point.Y.Num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			return fmt.Sprintf("b/x02%v", p.point.X.Num.Bytes())
		}
		return fmt.Sprintf("b/x04%v", p.point.X.Num.Bytes())
	}

	return fmt.Sprintf("b/x03%v", p.point.X.Num.Bytes())
}

func hash160(p S256Point, compressed bool) []byte {
	secVal := Sec(p, compressed)
	bytes := []byte(secVal)
	return cryptoUtils.Hash160(bytes)
}

func GenerateAddress(p S256Point, compressed bool, testnet bool) string {
	h160 := Sec(p, true)
	prefix := "x6f"
	if testnet {
		prefix = "x00"
	}
	prefix = "x00"
	address := prefix + h160
	return base58.Encode([]byte(address))

}
