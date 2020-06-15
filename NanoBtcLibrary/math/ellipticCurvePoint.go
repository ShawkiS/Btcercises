package math

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type ECPoint struct {
	X    float64
	Y    float64
	A    float64
	B    float64
	Infx bool
	Infy bool
}

func CreatePoint(x string, y string, a, b float64) (ECPoint, error) {

	if x == "" && y == "" {
		return ECPoint{0, 0, a, b, true, true}, nil
	}

	floatX, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return ECPoint{}, fmt.Errorf("converting string to float64")
	}
	floatY, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return ECPoint{}, fmt.Errorf("converting string to float64")
	}

	if math.Pow(floatY, 2) != math.Pow(floatX, 3)+a*floatX+b {
		return ECPoint{}, fmt.Errorf("(%v, %v) is not on the curve", x, y)
	}

	xInt, err := strconv.ParseFloat(x, 64)
	yInt, err := strconv.ParseFloat(y, 64)

	return ECPoint{xInt, yInt, a, b, false, false}, nil
}

func (p ECPoint) Add(po ECPoint) (ECPoint, error) {

	result := CheckECPointSameCurve(p, po)
	if result != true {
		return ECPoint{}, fmt.Errorf("")
	}

	if p.Infx {
		return po, nil
	}
	if po.Infx {
		return p, nil
	}

	if (p.X == po.X) && (p.Y != po.Y) {
		return ECPoint{p.A, p.B, 0, 0, true, true}, nil
	}

	if p.X != po.X {
		var s = (po.Y - p.Y) / (po.X - p.X)
		var x = math.Pow(s, 2) - p.X - po.X
		var y = s*(p.X-x) - po.Y

		return ECPoint{p.A, p.B, x, y, false, false}, nil
	}

	if (p == po) && (p.Y == 0*p.X) {
		return ECPoint{p.A, p.B, 0, 0, false, false}, nil
	}

	if p == po {
		var s = (3*math.Pow(p.X, 2) + p.A) / (2 * p.Y)
		var x = math.Pow(s, 2) - 2*p.X
		var y = s*(p.X-x) - p.Y

		return ECPoint{p.A, p.B, x, y, false, false}, nil
	}

	return ECPoint{}, errors.New("Point addition exemption")
}

func (p *ECPoint) IsEqual(po ECPoint) bool {
	return ((p.X == po.X) && (p.Y == po.Y)) && ((p.A == po.A) && (p.B == po.B))
}
func IsCheckAdditiveInverses(p1, p2 ECPoint) bool {

	if p1.X == p2.X && p1.Y != p2.Y {
		return true
	}

	return false
}
