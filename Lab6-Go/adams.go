package main

import (
	m "math"
)

type adamsDiffSolver struct {
	rightBorder float64
	step float64
	precision float64
}

func (s *adamsDiffSolver) SetStep(h float64) {
	s.step = h
}

func (s *adamsDiffSolver) Solve(eq *equation) []Point {
	noSteps := int((s.rightBorder - eq.x0)/s.step)

	P := make([]Point, noSteps+1)
	startPoints := (&rungeKuttaDiffSolver{
		rightBorder: eq.x0 + 4 * s.step,
		step:        s.step,
		precision:   s.precision,
	}).Solve(eq)

	copy(P, startPoints)

	f := make([]float64, noSteps+1)
	for i := 0; i < 4; i++ {
		f[i] = eq.exec(P[i].X, P[i].Y)
	}
	for i := 3; i < noSteps; i++ {
		P[i+1].X = P[i].X + s.step

		P[i+1].Y = P[i].Y + s.step*f[i]
		P[i+1].Y += m.Pow(s.step, 2)/2 * (f[i] - f[i-1])
		P[i+1].Y += 5*m.Pow(s.step, 3)/12 * (f[i] - 2*f[i-1] + f[i-2])
		P[i+1].Y += 3*m.Pow(s.step, 4)/8 * (f[i] - 3*f[i-1] + 3*f[i-2] - f[i-3])

		f[i+1] = eq.exec(P[i+1].X, P[i+1].Y)
	}
	return P
}
