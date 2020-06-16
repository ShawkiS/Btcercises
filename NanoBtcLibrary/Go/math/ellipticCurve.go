package math

func CheckECPointSameCurve(p1, p2 ECPoint) bool {
	if (p1.A != p2.A) || (p1.B != p2.B) {
		return false
	}

	return true
}

func CheckFFPointSameCurve(p1 *FFPoint, p2 *FFPoint) bool {
	p1a := p1.A
	p1b := p1.B
	p2a := p2.A
	p2b := p2.B

	if (!p1a.IsEqual(p2a)) || (!p1b.IsEqual(p2b)) {
		return false
	}

	return true
}
