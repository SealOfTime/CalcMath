package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	m "math"
	"os"
	"regexp"
	"strings"
)

type app struct {
	eq *equation
	a, b, e float64
}

//начальное приближение f(a) * f``(a) = +
var nac int
func main() {
	flag.IntVar(&nac,"p", 5, "number of signs after comma for floats in output")
	flag.Parse()

	var r *bufio.Reader
	var interactive bool

	inputFileName := flag.Arg(0)
	if len(inputFileName) != 0{
		f, err := os.Open(inputFileName)
		if err != nil {
			log.Fatalf("Файл с именем '%s' не был найден.", inputFileName)
		}
		r = bufio.NewReader(f)
		interactive = false
	} else {
		r = bufio.NewReader(os.Stdin)
		interactive = true
	}

	app := app{}
	if interactive{
		app.eq = promptEquation(r)
		app.e = promptPrecision(r)
		app.a, app.b = promptInterval(r, &app)
	} else {
		app.eq = readEquationFromFile(r)
		app.e = readPrecisionFromFile(r)
		app.a, app.b = readIntervalFromFile(r, &app)
	}

	if m.Abs(app.eq.exec(app.b)) <= app.e {
		fmt.Printf("Решение в данном случае тривиально. Корень вашего уравнения: %f\n", app.b)
		fmt.Println("Если вы не этого ожидали, то попробуйте снова")
	}

	if m.Abs(app.eq.exec(app.a)) <= app.e {
		fmt.Printf("Решение в данном случае тривиально. Корень вашего уравнения: %f\n", app.a)
		fmt.Println("Если вы не этого ожидали, то попробуйте снова")
	}
	var me int
	if interactive {
		me = promptMethod(r)
	} else {
		me = readMethodFromFile(r)
	}
	switch me {
	case 1:{
		s := BisectionMethod{a:  app.a, b:  app.b, e:  app.e, eq: app.eq}
		steps := s.Solve()

		fmt.Println("Поиск корня методом половинного деления: ")
		header := fmt.Sprintf(
			"%6s | %10s | %10s | %10s | %10s | %10s | %10s | %10s |\n",
		"№ шага", "a", "b", "x", "fa", "fb", "fx", "|a-b|",
		)
		separator := strings.Repeat("-", len(header))
		fmt.Println(header, separator)

		format := fmt.Sprintf("%%6d | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df |\n", nac, nac, nac, nac, nac, nac, nac)
		for i, s := range steps {
			fmt.Printf(
				format,
				i, s.a, s.b, s.x, s.fa, s.fb, s.fx, s.interval,
			)
		}

		fmt.Printf(fmt.Sprintf("Ответ: x ≈ %%10.%df\n", nac), steps[len(steps)-1].x)
	}
	case 2:{
		var x0, x1 float64
		if app.eq.exec(app.a) * app.eq.dd(app.a) > 0 {
			x0, x1 = app.a, app.a + 2 * app.e
			fmt.Printf("f(a) * f''(x) > 0:\nx0 = a = %10.5f; x1 = %10.5f\n", x0, x1)
		} else {
			x0, x1 = app.b, app.b - 2 * app.e
			fmt.Printf("x0 = b = %10.5f; x1 = %10.5f\n", x0, x1)
		}

		if m.Abs(app.eq.exec(x0)) <= app.e {
			fmt.Printf("Решение в данном случае тривиально. Корень вашего уравнения: %f\n", x0)
			fmt.Println("Если вы не этого ожидали, то попробуйте снова")
			return
		}

		if m.Abs(app.eq.exec(x1)) <= app.e {
			fmt.Printf("Решение в данном случае тривиально. Корень вашего уравнения: %f\n", x1)
			fmt.Println("Если вы не этого ожидали, то попробуйте снова")
			return
		}

		fmt.Println("Поиск корня методом секущих: ")
		s := SecantMethod{x0: x0, x1: x1, e:  app.e, eq: app.eq}
		steps := s.Solve()

		header := fmt.Sprintf(
			"%6s | %10s | %10s | %10s | %10s | %10s | %10s | %10s |\n",
			"№ шага", "x_(k-1)", "f(x_(k-1))", "x_k", "f(x_k)", "x_(k+1)", "f(x_(k+1))", "|x_k - x_(k+1)|",
		)
		separator := strings.Repeat("-", len(header))
		fmt.Println(header, separator)

		format := fmt.Sprintf("%%6d | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df |\n", nac, nac, nac, nac, nac, nac, nac)
		for i, s := range steps {
			fmt.Printf(format, i, s.prevX, s.fPrevX, s.x, s.fx, s.newX, s.fNewX, s.increment)
		}

		fmt.Printf(fmt.Sprintf("Ответ: x ≈ %%10.%df\n", nac), steps[len(steps)-1].newX)
	}
	case 3:{
		var x0 float64
		if app.eq.exec(app.a) * app.eq.dd(app.a) > 0 {
			x0 = app.a
			fmt.Printf("x0 = a = %10.5f\n", x0)
		} else {
			x0 = app.b
			fmt.Printf("x0 = %10.5f\n", x0)
		}

		s := IterationMethod{
			a: app.a,
			b: app.b,
			x0: x0,
			e:  app.e,
			eq: app.eq,
		}
		steps, err := s.Solve()
		fmt.Println("Поиск корня методом простой итерации: ")

		fmt.Printf("𝜑(x) = x + %10.5f * (%s)\n", s.lambda, s.eq.notation)
		fmt.Printf("𝜑'(a) = 𝜑'(%10.5f) = %10.5f; 𝜑'(b) = 𝜑'(%10.5f) = %10.5f\n",
			app.a, 1 + s.lambda * app.eq.d(app.a),
			app.b, 1 + s.lambda * app.eq.d(app.b),
		)
		if err != nil {
			fmt.Println("На этом интервале метод простой итерации не сходится.")
			return
		}

		header := fmt.Sprintf(
			"%6s | %10s | %10s | %10s | %10s | %10s |\n",
			"№ шага", "x_k", "f(x_k)", "x_(k+1)", "𝜑(x_k)", "|x_k - x_(k+1)|",
		)
		separator := strings.Repeat("-", len(header))
		fmt.Println(header, separator)

		format := fmt.Sprintf("%%6d | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df | %% 10.%df |\n", nac, nac, nac, nac, nac)
		for i, s := range steps {
			fmt.Printf(format, i, s.x, s.fx, s.newX, s.phiX, s.increment)
		}

		fmt.Printf(fmt.Sprintf("Ответ: x ≈ %%10.%df\n", nac), steps[len(steps)-1].newX)
	}
	}
}

var methods = []string {
	"Половинного деления",
	"Секущих",
	"Простой итерации",
}

var equations = []equation {
	{
		notation: "x^3 - x + 4",
		exec: func(x float64) float64 { return m.Pow(x, 3) - x + 4 },
		d: func(x float64) float64 { return 3*m.Pow(x, 2) - 1 },
		dd: func(x float64) float64 { return 6*x},
	},
	{
		notation: "sin(x)",
		exec: func(x float64) float64 { return m.Sin(x) },
		d: func(x float64) float64 { return m.Cos(x) },
		dd: func(x float64) float64 { return -m.Sin(x) },
	},
	{
		notation: "cos^2(x-2)",
		exec: func(x float64) float64 { return m.Pow(m.Cos(x-2), 2) },
		d: func(x float64) float64 { return m.Sin(4-2*x)},
		dd: func(x float64) float64 { return -2*m.Cos(4-2*x)},
	},
	{
		notation: "-2.7*x^3 - 1.48*x^2 + 19.23*x + 6.35",
		exec: func(x float64) float64 { return -2.7*m.Pow(x,3) - 1.48*m.Pow(x,2) + 19.23*x + 6.35},
		d: func(x float64) float64 {return -8.1*m.Pow(x, 2) - 2.96*x + 19.23},
		dd: func(x float64) float64 {return -16.2*x - 2.96},
	},
}

type equation struct {
	notation string
	exec func(x float64) (y float64)
	d func(x float64) float64
	dd func(x float64) float64
}

func promptMethod(r *bufio.Reader) int {
	fmt.Println("Пожалуйста, выберите метод решения вашего нелинейного уравнения...")
	for i, me := range methods {
		fmt.Printf("%d) %s\n", i+1, me)
	}

	for {
		fmt.Print("> ")
		method, err := readMethod(r)
		switch err {
		case ErrUnknown, ErrNotAnInt:
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
		case ErrNotInList:
			fmt.Println("Такого метода не предлагалось. Попробуйте снова.")
		case nil:
			return method
		}
	}
}

var powOfTenRegexp = regexp.MustCompile("0\\.0*1$")
func promptPrecision(r *bufio.Reader) float64 {
	fmt.Println("Пожалуйста введите точность вычислений: ")

	for {
		fmt.Print("> ")
		precision, err := readPrecision(r)
		switch err {
		case ErrUnknown, ErrNotANumber:
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
		case ErrNotNormal:
			fmt.Println("Точность должна принадлежать интервалу (0, 1)")
		case ErrNotPowOfTen:
			fmt.Println("Точность должна выражаться десяткой в отрицательной степени.")
		case nil:
			return precision
		}
	}
}

func promptEquation(r *bufio.Reader) (eq *equation) {
	fmt.Println("Пожалуйста, выберите ваше уравнение: ")
	for i, e := range equations {
		fmt.Printf("%d) %s = 0\n", i+1, e.notation)
	}

	for {
		fmt.Print("> ")
		eq, err := readEquation(r)
		switch err {
		case ErrUnknown, ErrNotAnInt:
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
		case ErrNotInList:
			fmt.Println("Такого уравнения не предлагалось. Попробуйте снова.")
		case nil:
			return eq
		}
	}
}

func promptInterval(r *bufio.Reader, app *app) (a, b float64) {
	var err error

	fmt.Println("Введите интервал для поиска корня через запятую:")
	for {
		fmt.Print("> ")
		a, b, err = readInterval(r, app)
		switch err {
		case ErrUnknown: fmt.Println("Ошибка ввода. Попробуйте снова")
		case ErrExpectedTwoNumbers: fmt.Println("Некорретный интервал. Попробуйте снова")
		case ErrANotANumber: fmt.Println("Ошибка ввода левой границы интервала. Попробуйте снова")
		case ErrBNotANumber: fmt.Println("Ошибка ввода правой границы интервала. Попробуйте снова")
		case ErrBadInterval: fmt.Println("Ну вы чего. Правая граница интервала должна быть строго больше левой. Попробуйте снова.")
		case ErrIntervalIsPoint: fmt.Println("Ваш интервал должен быть больше, чем точность вычислений. Иначе это все бессмысленно")
		case ErrNoConvergence: fmt.Println("На данном интервале либо совсем нету корней, либо несколько. Советую выбрать интервал получше.")
		case nil: return a, b
		}
	}
}

func promptSecantStartXs(r *bufio.Reader, app *app) (x0, x1 float64) {
	fmt.Println("Введите начальные приближения - x0 и x1 для поиска корня через запятую:")
	for {
		fmt.Print("> ")
		x0, x1, err := readSecantStartXs(r, app)
		switch err {
		case ErrUnknown: fmt.Println("Ошибка ввода. Попробуйте снова")
		case ErrExpectedTwoNumbers: fmt.Println("Некорретный интервал. Попробуйте снова")
		case ErrX0NotANumber: fmt.Println("Ошибка ввода x0. Попробуйте снова.")
		case ErrX0NotInInterval: fmt.Println("Выберите x0 в пределах заданного интервала.")
		case ErrX1NotANumber:fmt.Println("Ошибка ввода x1. Попробуйте снова")
		case ErrX1NotInInterval: fmt.Println("Выберите x1 в пределах заданного интервала.")
		case ErrTwoNumbersAreOne: fmt.Println("Выберите x0 и x1 отстоящия друг от друга на расстояние большее, чем точность вычислений.")
		case nil: return x0, x1
		}
	}
}

func promptIterationStartX(r *bufio.Reader, app *app) (x0 float64) {
	fmt.Println("Введите начальные приближения - x0 для поиска корня:")
	for {
		fmt.Print("> ")
		x0, err :=  readIterationStartX(r, app)
		switch err {
		case ErrUnknown: fmt.Println("Ошибка ввода. Попробуйте снова")
		case ErrNotANumber: fmt.Println("Ошибка ввода x0. Попробуйте снова")
		case ErrNotInInterval: fmt.Println("Выберите x0 в пределах заданного интервала.")
		case nil: return x0
		}
	}
}