package math

import (
	"fmt"
	"math/big"

	"github.com/Btcercises/NanoBtcLibrary/math/utils"
)

var A = big.NewInt(0)
var B = big.NewInt(7)

type FFPoint struct {
	X *FieldElement
	Y *FieldElement
	A *FieldElement
	B *FieldElement
}

func (p *FFPoint) printPoint() {

	x := p.X
	y := p.Y

	xx := x.Num.String()
	yy := y.Num.String()
	pr := prime.String()

	if xx == "0" || yy == "0" {
		fmt.Println("Sorry, it's infinity")
		return
	}

	fmt.Printf("Point (%s, %s)_%d_%d FieldElement(%s)\n", xx, yy, A, B, pr)
}

func CreateNewPoint(x, y big.Int) FFPoint {

	pnt := FFPoint{
		X: CreateFieldElement(x),
		Y: CreateFieldElement(y),
		A: CreateFieldElement(*A),
		B: CreateFieldElement(*B),
	}

	if (x.Cmp(big.NewInt(0)) == 0) || (y.Cmp(big.NewInt(0)) == 0) {
		return pnt
	}

	onCurve := utils.IsSameField(&x, &y)
	if !onCurve {

		xx := x.String()
		yy := y.String()

		resStr := fmt.Sprintf("Point (%s, %s)_%d_%d is not on the curve", xx, yy, A, B)
		panic(resStr)
	}

	return pnt
}

func (p *FFPoint) Add(po *FFPoint) *FFPoint {

	isSameCurve := CheckFFPointSameCurve(p, po)
	if !isSameCurve {
		panic(isSameCurve)
	}

	if p.X.Num.Cmp(big.NewInt(0)) == 0 {

		return po
	}

	if po.X.Num.Cmp(big.NewInt(0)) == 0 {

		return p
	}

	if (p.X.Num.Cmp(po.X.Num) == 0) && (p.Y.Num.Cmp(po.Y.Num) != 0) {

		infa := CreateFieldElement(*big.NewInt(0))
		infb := CreateFieldElement(*big.NewInt(0))
		return &FFPoint{infa, infb, CreateFieldElement(*A), CreateFieldElement(*B)}
	}

	if !p.X.IsEqual(po.X) {

		ss1 := po.Y.Sub(*p.Y)
		ss2 := po.X.Sub(*p.X)

		sDiv := ss1.Div(ss2)

		x3 := sDiv.Pow(*big.NewInt(2))
		x3 = x3.Sub(*p.X)
		x3 = x3.Sub(*po.X)

		y := p.X.Sub(x3)
		y = sDiv.Mul(y)
		y = y.Sub(*p.Y)

		return &FFPoint{&x3, &y, p.A, p.B}

	}

	if p.IsEqual(*po) {

		x12 := p.X.Pow(*big.NewInt(2))
		fe3 := CreateFieldElement(*big.NewInt(3))
		sNom := x12.Mul(*fe3)
		sNom = sNom.Add(*p.A)

		fe2 := CreateFieldElement(*big.NewInt(2))
		sDom := p.Y.Mul(*fe2)

		sDiv := sNom.Div(sDom)

		x3 := sDiv.Pow(*big.NewInt(2))
		xx := p.X.Mul(*fe2)
		x3 = x3.Sub(xx)

		y := p.X.Sub(x3)
		y = sDiv.Mul(y)
		y = y.Sub(*p.Y)

		return &FFPoint{&x3, &y, p.A, p.B}
	}

	if (p.IsEqual(*po)) && (p.Y.Num.Cmp(big.NewInt(0)) == 0) {

		infa := CreateFieldElement(*big.NewInt(0))
		infb := CreateFieldElement(*big.NewInt(0))
		return &FFPoint{infa, infb, CreateFieldElement(*A), CreateFieldElement(*B)}
	}

	panic("Point addition exemption: no condition fulfilled")
}

func (p *FFPoint) IsEqual(po FFPoint) bool {

	if !p.X.IsEqual(po.X) || !p.Y.IsEqual(po.Y) || !p.A.IsEqual(po.A) || !p.B.IsEqual(po.B) {
		return false
	}

	return true
}

func (p *FFPoint) RMul(coef big.Int) *FFPoint {
	current := p

	newPoint := CreateNewPoint(*big.NewInt(0), *big.NewInt(0))
	result := &newPoint

	if coef.Cmp(big.NewInt(1)) == -1 {
		return result
	}

	result = &newPoint

	if coef.Cmp(big.NewInt(1)) == 0 {
		return result
	}

	var coefBit big.Int
BitShift:
	for {

		coefBit.And(&coef, big.NewInt(1))
		if coefBit.Cmp(big.NewInt(1)) == 0 {
			result = result.Add(current)
		}
		current = current.Add(current)

		if coef.Cmp(big.NewInt(0)) == 0 {
			break BitShift
		}

	}

	return result
}
