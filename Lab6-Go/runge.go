package main

type rungeKuttaDiffSolver struct {
	rightBorder float64
	step float64
	precision float64
}

func (s *rungeKuttaDiffSolver) SetStep(h float64) {
	s.step = h
}

func (s *rungeKuttaDiffSolver) Solve(eq *equation) []Point {
	noSteps := int((s.rightBorder - eq.x0)/s.step)

	P := make([]Point, noSteps+1)
	P[0].Y = eq.y0
	P[0].X = eq.x0

	var xi, yi, k1, k2, k3, k4 float64
	for i := 0; i < noSteps; i++ {
		xi, yi = P[i].X, P[i].Y
		k1 = s.step * eq.exec(xi, yi)
		k2 = s.step * eq.exec(xi + s.step/2, yi + k1/2)
		k3 = s.step * eq.exec(xi + s.step/2, yi + k2/2)
		k4 = s.step * eq.exec(xi + s.step, yi + k3)

		P[i+1].X = xi + s.step
		P[i+1].Y = yi + (k1 + 2*k2 + 2*k3 + k4)/6
	}

	return P
}