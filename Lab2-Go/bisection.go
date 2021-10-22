package main

import (
	"log"
	m "math"
)

type BisectionMethod struct {
	a float64
	b float64
	e float64
	eq *equation
}

type BisectionMethodStep struct {
	a, fa float64
	b, fb float64
	x, fx float64
	interval float64
}
func (s *BisectionMethod) Solve() (steps []BisectionMethodStep) {
	apprNoOfSteps := int(m.Log2(m.Abs(s.a - s.b)/s.e))
	steps = make([]BisectionMethodStep, 0, apprNoOfSteps+1)

	a, b := s.a, s.b
	fa, fb := s.eq.exec(a), s.eq.exec(b)
	interval := m.Abs(a - b)

	x := (a+ b) / 2
	fx := s.eq.exec(x)

	steps = append(steps, BisectionMethodStep{
		a: a, fa: fa,
		b: b, fb: fb,
		x: x, fx: fx,
		interval: interval,
	}) //0th step
	for m.Abs(fx) > s.e && interval > s.e {
		switch {
		case fa*fx < 0: b, fb = x, fx
		case fb*fx < 0: a, fa = x, fx
		default: log.Fatalln("Не было интервала с разными знаками")
		}
		x = (a + b) / 2
		fx = s.eq.exec(x)
		interval = m.Abs(b-a)
		steps = append(steps, BisectionMethodStep{
			a: a, fa: fa,
			b: b, fb: fb,
			x: x, fx: fx,
			interval: interval,
		})
	}
	return steps
}
