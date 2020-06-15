package crypto

import (
	"crypto/rand"
	"math/big"

	btcMath "github.com/Btcercises/NanoBtcLibrary/math"
	"github.com/Btcercises/NanoBtcLibrary/math/utils"
)

var G = gValue()

type PrivateKey struct {
	secret *big.Int
	point  *S256Point
}

func NewPrivateKey(secret *big.Int) PrivateKey {

	np := G.S256RMul(*secret)

	privK := PrivateKey{
		secret: secret,
		point:  np,
	}

	return privK
}

func (pk *PrivateKey) sign(z *big.Int) Signature {
	var kInv, sFinal big.Int
	n := utils.HexToBigInt(N)
	nField := btcMath.CreateFieldElement(*n)

	nMinTwo := nField.Sub(*btcMath.CreateFieldElement(*big.NewInt(2)))

	k := pk.deterministicK(z)
	r := G.S256RMul(*k)
	kInv.Exp(k, nMinTwo.Num, n)

	zField := btcMath.CreateFieldElement(*z)
	s := r.point.X.Add(*zField)
	sPoint := pk.point.S256RMul(*s.Num)
	sPoint = sPoint.S256RMul(kInv)

	sFinal.Mod(sPoint.point.X.Num, n)

	nDiv := nField.Div(*btcMath.CreateFieldElement(*big.NewInt(2)))

	if sFinal.Cmp(nDiv.Num) == 1 {
		sRet := nField.Sub(*btcMath.CreateFieldElement(sFinal))
		sFinal = *sRet.Num
	}

	return Signature{r.point.X.Num, &sFinal}
}

func (pk *PrivateKey) deterministicK(z *big.Int) *big.Int {
	rNum, _ := rand.Int(rand.Reader, z)

	return rNum
}
