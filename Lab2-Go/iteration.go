package main

import (
	"errors"
	m "math"
)

type IterationMethod struct {
	a float64
	b float64
	x0 float64
	e float64
	lambda float64
	eq *equation
}

type IterationMethodStep struct {
	x, newX float64
	phiX, fx float64
	increment float64
}
func (s *IterationMethod) Solve() (steps []IterationMethodStep, err error) {
	s.lambda = s.findPhi()
	q := s.findQ()
	if q > 1 {
		return steps, errors.New("на данном интервале итерационный метод не сходится")
	}

	var convCrit float64
	if q > 0.5 {
		convCrit = s.e * (1 - q) / q
	} else {
		convCrit = s.e
	}

	x, newX := s.x0, 0.0

	steps = make([]IterationMethodStep, 0, 20)
	increment := 1.0
	for increment >= convCrit {
		newX = x + s.lambda * s.eq.exec(x)
		increment = m.Abs(newX - x)

		steps = append(steps, IterationMethodStep{
			x:         x,
			newX:      newX,
			phiX:      x + s.lambda * s.eq.exec(x),
			fx:        s.eq.exec(x),
			increment: increment,
		})

		x = newX
	}

	return
}

func (s *IterationMethod) findPhi() (lambda float64) {
	maxDerivative := 0.0
	for x := s.a; x < s.b; x += s.e {
		d := s.eq.d(x)
		if d > maxDerivative {
			maxDerivative = d
		}
	}

	return -1/maxDerivative
}

func (s *IterationMethod) findQ() (q float64) {
	var dPhi, maxDerivative float64
	for x := s.a; x < s.b; x += s.e {
		dPhi = 1 + s.lambda * s.eq.d(x)
		if m.Abs(dPhi) > maxDerivative {
			maxDerivative = m.Abs(dPhi)
		}
	}
	return maxDerivative
}