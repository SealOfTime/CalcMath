package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	m "math"
	"strconv"
	"strings"
)

var (
	ErrUnknown = errors.New("вы ввели что-то странное")
	ErrNotAnInt = errors.New("ожидалось целое число")
	ErrNotANumber = errors.New("ожидалось число")
	ErrExpectedTwoNumbers = errors.New("ожидалось два числа, разделённых запятой")
)

var (
	ErrNotInList = errors.New("элемента с таким номером не представлено")
)

func readMethod(r *bufio.Reader) (method int, err error) {
	raw, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF){
		return method, ErrUnknown
	}

	input, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return method, ErrNotAnInt
	}

	if input < 1 || int(input) > len(methods) {
		return method, ErrNotInList
	}

	return int(input), nil
}

var (
	ErrNotNormal = errors.New("ожидалось число на интервале (0, 1)")
	ErrNotPowOfTen = errors.New("ожидалась десятка в отрицательной степени")
)
func readPrecision(r *bufio.Reader) (precision float64, err error) {
	raw, err := r.ReadString('\n')
	if err != nil {
		return precision, ErrUnknown
	}
	precision, err = strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return precision, ErrNotANumber
	}

	if precision < 0 || precision > 1 {
		return precision, ErrNotNormal
	}

	if !powOfTenRegexp.MatchString(strings.TrimSpace(raw)) {
		return precision, ErrNotPowOfTen
	}

	return precision, nil
}

func readEquation(r *bufio.Reader) (eq *equation, err error) {
	raw, err := r.ReadString('\n')
	if err != nil {
		return nil, ErrUnknown
	}

	input, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		return nil, ErrNotAnInt
	}

	if input < 1 || int(input) > len(equations) {
		return nil, ErrNotInList
	}

	return &equations[int(input)-1], nil
}

var (
	ErrANotANumber = fmt.Errorf("ошибка ввода левой границы: %w", ErrNotANumber)
	ErrBNotANumber = fmt.Errorf("ошибка ввода правой границы: %w", ErrNotANumber)
	ErrBadInterval = errors.New("правая граница интервала должна быть больше левой")
	ErrIntervalIsPoint = errors.New("правая граница не должна находиться в эпсилон-окрестности левой")
	ErrNoConvergence = errors.New("на данном интервале метод не сходится")
)
func readInterval(r *bufio.Reader, app *app) (a, b float64, err error) {
	input, err := r.ReadString('\n')
	if err != nil {
		return a, b, ErrUnknown
	}

	interval := strings.Split(input, ",")
	if len(interval) != 2 {
		return a, b, ErrExpectedTwoNumbers
	}

	a, err = strconv.ParseFloat(strings.TrimSpace(interval[0]), 64)
	if err != nil {
		return a, b, ErrANotANumber
	}

	b, err = strconv.ParseFloat(strings.TrimSpace(interval[1]), 64)
	if err != nil {
		return a, b, ErrBNotANumber
	}

	if b <= a {
		return a, b, ErrBadInterval
	}

	if m.Abs(b - a) <= app.e {
		return a, b, ErrIntervalIsPoint
	}

	if app.eq.exec(b)*app.eq.exec(a) > 0 {
		return a, b, ErrNoConvergence
	}

	return
}

var (
	ErrX0NotANumber = fmt.Errorf("ошибка ввода x0: %w", ErrNotANumber)
	ErrX1NotANumber = fmt.Errorf("ошибка ввода x1: %w", ErrNotANumber)
	ErrX0NotInInterval = fmt.Errorf("ошибка ввода x0: %w", ErrNotInInterval)
	ErrX1NotInInterval = fmt.Errorf("ошибка ввода x1: %w", ErrNotInInterval)
	ErrTwoNumbersAreOne = errors.New("x1 должно не находиться в эпсилон-окрестности x0")
)
func readSecantStartXs(r *bufio.Reader, app *app) (x0, x1 float64, err error) {
	input, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF){
		return x0, x1, ErrUnknown
	}

	interval := strings.Split(input, ",")
	if len(interval) != 2 {
		return x0, x1, ErrExpectedTwoNumbers
	}

	x0, err = strconv.ParseFloat(strings.TrimSpace(interval[0]), 64)
	if err != nil {
		return x0, x1, ErrX0NotANumber
	}
	if x0 < app.a || x0 > app.b {
		return x0, x1, ErrX0NotInInterval
	}

	x1, err = strconv.ParseFloat(strings.TrimSpace(interval[1]), 64)
	if err != nil {
		return x0, x1, ErrX1NotANumber
	}
	if x1 < app.a || x1 > app.b {
		return x0, x1, ErrX1NotInInterval
	}

	if m.Abs(x0 - x1) <= app.e {
		return x0, x1, ErrTwoNumbersAreOne
	}

	return x0, x1, nil
}

var (
	ErrNotInInterval = errors.New("число должно находиться на исследуемом интервале")
)
func readIterationStartX(r *bufio.Reader, app *app) (x0 float64, err error) {
	input, err := r.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF){
		return x0, ErrUnknown
	}

	x0, err = strconv.ParseFloat(strings.TrimSpace(input), 64)
	if err != nil {
		return x0, ErrNotANumber
	}

	if x0 < app.a || x0 > app.b {
		return x0, ErrNotInInterval
	}

	return x0, nil
}