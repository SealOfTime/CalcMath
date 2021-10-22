package main

import (
	"bufio"
	"fmt"
	"log"
)

func readMethodFromFile(r *bufio.Reader) int {
	method, err := readMethod(r)
	switch err {
	case ErrUnknown, ErrNotAnInt:
		log.Fatalln("Ошибка ввода метода. Это должно быть целое число")
	case ErrNotInList:
		log.Fatalln("Такого метода не предлагалось. Попробуйте снова.")
	}
	return method
}

func readPrecisionFromFile(r *bufio.Reader) float64 {
	precision, err := readPrecision(r)
	fmt.Printf("Необходимая точность: %f\n", precision)
	switch err {
	case ErrUnknown, ErrNotANumber:
		log.Fatalln("Ошибка ввода точности. Это должно быть число")
	case ErrNotNormal:
		log.Fatalln("Точность должна принадлежать интервалу (0, 1)")
	case ErrNotPowOfTen:
		log.Fatalln("Точность должна выражаться десяткой в отрицательной степени.")
	}
	
	return precision
}

func readEquationFromFile(r *bufio.Reader) (eq *equation) {
	eq, err := readEquation(r)
	fmt.Printf("Исследуемое уравнение: %s = 0\n", eq.notation)
	switch err {
	case ErrUnknown, ErrNotAnInt:
		log.Fatalln("Ошибка ввода уравнения. Это должно было быть целое число.")
	case ErrNotInList:
		log.Fatalln("Такого уравнения не предлагалось. Попробуйте снова.")
	}
	return eq
}

func readIntervalFromFile(r *bufio.Reader, app *app) (a, b float64) {
	var err error
	a, b, err = readInterval(r, app)
	fmt.Printf("Исследуемый интервал a: %f; b: %f\n", a, b)
	switch err {
	case ErrUnknown: log.Fatalln("Ошибка ввода. Попробуйте снова")
	case ErrExpectedTwoNumbers: log.Fatalln("Некорретный интервал. Попробуйте снова")
	case ErrANotANumber: log.Fatalln("Ошибка ввода левой границы интервала. Попробуйте снова")
	case ErrBNotANumber: log.Fatalln("Ошибка ввода правой границы интервала. Попробуйте снова")
	case ErrBadInterval: log.Fatalln("Ну вы чего. Правая граница интервала должна быть строго больше левой. Попробуйте снова.")
	case ErrIntervalIsPoint: log.Fatalln("Ваш интервал должен быть больше, чем точность вычислений. Иначе это все бессмысленно")
	case ErrNoConvergence: log.Fatalln("На данном интервале либо совсем нету корней, либо несколько. Советую выбрать интервал получше.")
	}
	return a, b
}

func readSecantStartXsFromFile(r *bufio.Reader, app *app) (x0, x1 float64) {
	x0, x1, err := readSecantStartXs(r, app)
	fmt.Printf("Начальные приближения X0: %f; X1: %f\n", x0, x1)
	switch err {
	case ErrUnknown: log.Fatalln("Ошибка ввода. Попробуйте снова")
	case ErrExpectedTwoNumbers: log.Fatalln("Некорретный интервал. Попробуйте снова")
	case ErrX0NotANumber: log.Fatalln("Ошибка ввода x0. Попробуйте снова.")
	case ErrX0NotInInterval: log.Fatalln("Выберите x0 в пределах заданного интервала.")
	case ErrX1NotANumber:log.Fatalln("Ошибка ввода x1. Попробуйте снова")
	case ErrX1NotInInterval: log.Fatalln("Выберите x1 в пределах заданного интервала.")
	case ErrTwoNumbersAreOne: log.Fatalln("Выберите x0 и x1 отстоящия друг от друга на расстояние большее, чем точность вычислений.")
	}
	return x0, x1
}

func readIterationStartXFromFile(r *bufio.Reader, app *app) (x0 float64) {
	x0, err :=  readIterationStartX(r, app)
	fmt.Printf("Начальное приближение X0: %f\n", x0)
	switch err {
	case ErrUnknown: log.Fatalln("Ошибка ввода начального приближения X0 для метода итераций")
	case ErrNotANumber: log.Fatalln("Ошибка ввода x0.")
	case ErrNotInInterval: log.Fatalln("Выберите x0 в пределах заданного интервала.")
	}
	return x0
}
