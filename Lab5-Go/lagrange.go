package main

type LagrangePolynomial struct {
	DataSet
}

func (s *LagrangePolynomial) EvaluateAt(x float64) (y float64) {
	n := len(s.Y)
	for i := 0; i < n; i++ {
		y += s.Y[i] * s.li(x, i)
	}
	return
}

func (s *LagrangePolynomial) li(x float64, i int) float64 {
	li := 1.0
	for j := 0; j < len(s.Y); j++ {
		if i == j {
			continue
		}
		li *= (x - s.X[j]) / (s.X[i] - s.X[j])
	}
	return li
}
