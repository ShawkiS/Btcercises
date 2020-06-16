package crypto

import (
	"math/big"

	"github.com/Btcercises/NanoBtcLibrary/Go/math"
)

type S256Field struct {
	field *math.FieldElement
}

func NewS256Field(num big.Int) S256Field {

	fe := math.CreateFieldElement(num)

	fld := S256Field{fe}

	return fld
}
