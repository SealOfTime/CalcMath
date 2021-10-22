package main

import m "math"

type SecantMethod struct {
	x0 float64
	x1 float64
	e float64
	eq *equation
}

type SecantMethodStep struct {
	prevX, fPrevX float64
	x, fx float64
	newX, fNewX float64
	increment float64
}
func (s *SecantMethod) Solve() (steps []SecantMethodStep) {
	prevX, x, newX := s.x0, s.x1, 0.0
	fPrevX, fx, fNewX := s.eq.exec(prevX), s.eq.exec(x), 0.0

	increment := m.Abs(x - prevX)
	for increment > s.e && m.Abs(fx) > s.e {
		newX = x - (x - prevX)/(fx - fPrevX) * fx
		fNewX = s.eq.exec(newX)

		increment = m.Abs(newX - x)

		steps = append(steps, SecantMethodStep{
			prevX: prevX, fPrevX: fPrevX,
			x: x, fx: fx,
			newX: newX, fNewX: fNewX,
			increment: increment,
		})

		prevX, fPrevX = x, fx
		x, fx = newX, fNewX
	}

	return steps
}