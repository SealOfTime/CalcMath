package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readDataSetFromStd() DataSet {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("Добро пожаловать в программу интерполяции каких-то функций какими-то методами!")
	fmt.Println("Выберите: ")
	fmt.Println("1. cos(x)")
	fmt.Println("2. 2 - x")
	fmt.Println("3. x^3 + 5x - 120")
	fmt.Println("4. Табличная функция")
	for {
		fmt.Print("Ваш выбор: ")
		input, _, err := r.ReadRune()
		if err != nil {
			fmt.Println("Ошибка ввода. Попробуйте снова")
			continue
		}

		switch input {
		case '1':
			fmt.Println("Вы выбрали cos(x).")
			a, b := readInterval(r)
			return generateInterval(a, b, func(x float64) float64 {
				return math.Cos(x)
			})
		case '2':
			fmt.Println("Вы выбрали 2 - x.")
			a, b := readInterval(r)
			return generateInterval(a, b, func(x float64) float64 {
				return 2 - x
			})
		case '3':
			fmt.Println("Вы выбрали x^3 + 5x - 120.")
			a, b := readInterval(r)
			return generateInterval(a, b, func(x float64) float64 {
				return math.Pow(a, 3.0)+5*a-120
			})
		case '4':
			fmt.Println("Вы выбрали табличную функцию.")
			return DataSet{
				X: []float64{0.25, 0.30, 0.35, 0.40, 0.45, 0.50, 0.55},
				Y: []float64{1.2557, 2.1764, 3.1218, 4.0482, 5.9875, 6.9195, 7.8359},
			}
		default:
			fmt.Println("Ошибка ввода. Попробуйте снова")
			continue
		}
	}
}

func generateInterval(a, b float64, f func(float64) float64) DataSet {
	x, y := make([]float64, 0, 6), make([]float64, 0, 6)
	fmt.Println(a, b)
	for step := (b - a) / 6; a <= b; a += step {
		x = append(x, a)
		y = append(y, f(a))
	}
	return DataSet{x, y}
}

func readInterval(r *bufio.Reader) (a, b float64) {
	fmt.Println("Введите интервал для интерполяции через запятую:")
	for {
		fmt.Print("Интервал: ")
		input, err := r.ReadString(byte('\n'))
		if err != nil {
			fmt.Printf("%+v", err)
			fmt.Println("Ошибка ввода. Попробуйте снова")
			continue
		}

		interval := strings.Split(input, ",")
		if len(interval) != 2 {
			fmt.Println("Некорретный интервал. Попробуйте снова")
			continue
		}

		a, err = strconv.ParseFloat(strings.TrimSpace(interval[0]), 64)
		if err != nil {
			fmt.Println("Ошибка ввода левой границы интервала. Попробуйте снова")
			continue
		}

		b, err = strconv.ParseFloat(strings.TrimSpace(interval[1]), 64)
		if err != nil {
			fmt.Println("Ошибка ввода правой границы интервала. Попробуйте снова")
			continue
		}
		return
	}
}
