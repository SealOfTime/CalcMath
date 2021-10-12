package main

import (
	"fmt"
	"math"
)

type NewtonPolynomialFinDiff struct {
	//h это шаг интерполяционной формулы
	h float64

	//i это начальный (конечный для интерполирования назад) индекс интервала,
	//на котором применяется интерполяционная формула.
	i int

	//backwards показывает направления интерполяции. True = интерполирование назад по II-й формуле Ньютона
	backwards bool

	DataSet
}

func (s *NewtonPolynomialFinDiff) EvaluateAt(x float64) float64 {
	if x < s.X[0] {
		fmt.Printf("Экстраполируем назад значение F(%f). i=%d. I-й формулой Ньютона\n", x, 0)
		return s.firstEvalAt(x, 0)
	}

	if x > s.X[len(s.X)-1] {
		fmt.Printf("Экстраполируем вперёд значение F(%f). i=%d. II-й формулой Ньютона\n", x, len(s.X) - 1)
		return s.secondEvalAt(x, len(s.X)-1)
	}

	half := math.Mod(x, s.h)
	if half > 0.05 {
		return s.EvaluateAtSecondHalf(x)
	} else {
		return s.EvaluateAtFirstHalf(x)
	}
}

func (s *NewtonPolynomialFinDiff) EvaluateAtFirstHalf(x float64) (y float64) {
	i, xi := 0, 0.0
	for i, xi = range s.X {
		if x < xi {
			i -= 1
			break
		}
	}
	fmt.Printf("Интерполируем вперёд значение F(%f). i=%d I-й формулой Ньютона  \n", x, i)
	return s.firstEvalAt(x, i)
}

func (s *NewtonPolynomialFinDiff) EvaluateAtSecondHalf(x float64) (y float64) {
	i, xi := 0, 0.0
	for i, xi = range s.X {
		if x < xi {
			break
		}
	}
	fmt.Printf("Интерполируем назад значение F(%f). i=%d. II-й формулой Ньютона\n", x, i)
	return s.secondEvalAt(x, i)
}

func (s *NewtonPolynomialFinDiff) secondEvalAt(x float64, i int) (y float64) {
	t := (x - s.X[i]) / s.h

	step := 1.0
	for k := i; k >= 0; k-- {
		y += s.FinDiff(i-k, k) * step
		step *= (t + float64(i-k)) / float64(i-k+1)
	}
	return
}

func (s *NewtonPolynomialFinDiff) firstEvalAt(x float64, i int) (y float64) {
	n := float64(len(s.Y) - i)
	t := (x - s.X[i]) / s.h

	step := 1.0
	for k := 0.0; k < n; k++ {
		y += s.FinDiff(int(k), i) * step
		step *= (t - k) / (k + 1)
	}
	return
}

func (s *NewtonPolynomialFinDiff) FinDiff(k, i int) float64 {
	if k == 0 {
		return s.Y[i]
	}

	return s.FinDiff(k-1, i+1) - s.FinDiff(k-1, i)
}
