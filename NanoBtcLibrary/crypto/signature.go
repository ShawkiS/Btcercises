package crypto

import (
	"fmt"
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func NewSignature(r, s *big.Int) Signature {
	return Signature{r, s}
}

func (f *Signature) print() {
	fmt.Printf("Signature(%s, %s)\n", f.s.String(), f.r.String())
}
