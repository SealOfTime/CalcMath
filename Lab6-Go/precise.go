package main

type preciseSolver struct {
	rightBorder float64
	step float64
	precision float64
}

func (s *preciseSolver) SetStep(h float64) {
	s.step = h
}

func (s *preciseSolver) Solve(eq *equation) []Point {
	noSteps := int((s.rightBorder - eq.x0)/s.step)

	P := make([]Point, noSteps+1)
	P[0].Y = eq.y0
	P[0].X = eq.x0
	for i := 0; i < noSteps; i++ {
		P[i+1].X = P[i].X + s.step
		P[i+1].Y = eq.preciseSolution(P[i+1].X)
	}

	return P
}


