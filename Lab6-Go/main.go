package main

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"log"
	m "math"
	"os"
	"strconv"
	"strings"
)

var equations = [3]equation{
	{
		notation:        "y' = y + (1 + x) * y^2; y(1) = -1",
		exec:            func(x, y float64) float64 { return y + (1+x)*m.Pow(y, 2.0) },
		preciseSolution: func(x float64) float64 { return -1/x },

		x0: 1,
		y0: -1,
	},
	{
		notation:        "y' = (x-y)^2 + 1; y(0) = 0",
		exec:            func(x, y float64) float64 { return m.Pow(x-y, 2) + 1 },
		preciseSolution: func(x float64) float64 { return x },

		x0: 0,
		y0: 0,
	},
	{
		notation:        "y' = xe^(-x^2) - 2 * xy; y(-1) = 1/2e",
		exec:            func(x, y float64) float64 { return x * m.Exp(-m.Pow(x, 2)) - 2*x*y },
		preciseSolution: func(x float64) float64 { return m.Pow(x, 2)/2 * m.Exp(-m.Pow(x, 2))},

		x0: -1,
		y0: 1/(2*m.E),
	},
}

type exec func(x, y float64) float64

type equation struct {
	notation string
	x0       float64
	y0       float64

	preciseSolution func(x float64) float64
	exec
}

type Point struct {
	X float64
	Y float64
}

type DifferentialSolver interface {
	solve(eq *equation) []Point
}

func main() {
	r := bufio.NewReader(os.Stdin)
	eq  := promptFunc(r)
	b := promptRightBorder(r)
	h := promptStep(r)
	precision := promptPrecision(r)

	rk := rungeKuttaDiffSolver{
		rightBorder: b,
		step:        h,
		precision:   precision,
	}

	ad := adamsDiffSolver{
		rightBorder: b,
		step:        h,
		precision:   precision,
	}

	p := preciseSolver{
		rightBorder: b,
		step: h,
		precision: precision,
	}

	ideal := p.solve(eq)
	rkAns := rk.solve(eq)
	adAns := ad.solve(eq)
	fmt.Println(" i |     xi     | yi(Рун-Кут) | yi(Адамса) | yi(Точн.)  | ")
	fmt.Println("-------------------------------------------")
	for i, p := range ideal {
		fmt.Printf("%2d | % 10.5f | % 10.5f |  % 10.5f |  % 10.5f\n", i, p.X, rkAns[i].Y, adAns[i].Y, p.Y)
	}
	plot(eq, ideal, rkAns, adAns)
	fmt.Println("Вы можете найти график в рабочей директории.")
}

func promptFunc(r *bufio.Reader) *equation {
	fmt.Println("Выберите, пожалуйста, функцию, для которой вы собираетесь решить задачу Коши: ")
	for i, v := range equations {
		fmt.Printf("%d) %s\n", i+1, v.notation)
	}

	var input int64
	for {
		fmt.Print("> ")
		raw, err := r.ReadString('\n')
		if err == nil {
			input, err = strconv.ParseInt(strings.Trim(raw, "\n"), 10, 32)
		}
		if err != nil {
			fmt.Println("Вы ввели что-то не то. Быть может, это было не целое число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
			continue
		}

		if input <= 0 || int(input) > len(equations) {
			fmt.Println("Такого варианта нету, внимательнее прочитайте предоставленные варианты. Иного не дано, даже не пытайтесь. Повторите снова свой ввод.")
			continue
		}

		return &equations[input - 1]
	}
}

func promptRightBorder(r *bufio.Reader) float64 {
	fmt.Println("Настоятельно вас прошу ввести правую границу для решения задачи Коши: ")
	var input float64
	for {
		fmt.Print("> ")
		raw, err := r.ReadString('\n')
		if err == nil {
			input, err = strconv.ParseFloat(strings.Trim(raw, "\n"), 64)
		}
		if err != nil {
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
			continue
		}

		return input
	}
}

func promptStep(r *bufio.Reader) float64 {
	fmt.Println("Настоятельно вас прошу ввести шаг для решения задач Коши численными методами: ")
	var input float64
	for {
		fmt.Print("> ")
		raw, err := r.ReadString('\n')
		if err == nil {
			input, err = strconv.ParseFloat(strings.Trim(raw, "\n"), 64)
		}
		if err != nil {
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
			continue
		}

		return input
	}
}

func promptPrecision(r *bufio.Reader) float64 {
	fmt.Println("Настоятельно вас прошу ввести точность (количество знаков после запятой) для применения численных методов: ")
	var input int64
	for {
		fmt.Print("> ")
		raw, err := r.ReadString('\n')
		if err == nil {
			input, err = strconv.ParseInt(strings.Trim(raw, "\n"), 10, 32)
		}
		if err != nil {
			fmt.Println("Вы ввели что-то не то. Быть может, это было не число. Быть может, какой-то странный символ. Итог один - пробуйте пока не получится.")
			continue
		}

		return m.Pow(10, -float64(input))
	}
}

func plot(eq *equation, precise []Point, runge []Point, adams []Point) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Plot",
			Subtitle: eq.notation,
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:          true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:        "value",
			Show:        true,
			SplitNumber: 10,
			Scale:       true,
			Min:         precise[0].X,
			Max:         precise[len(precise)-1].X,
			SplitArea:   &opts.SplitArea{
				Show:      true,
			},
			SplitLine:   &opts.SplitLine{
				Show:      true,
			},
			AxisLabel:   nil,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
		}),
	)
	line.AddSeries(
		"Точное решение",
		convertToLineData(precise),
		charts.WithLineChartOpts(opts.LineChart{
			Smooth:       true,
		}),
	)
	line.AddSeries(
		"Рунге-Кутта 4 порядка",
		convertToLineData(runge),
	)
	line.AddSeries(
		"Адамса",
		convertToLineData(adams),
	)

	f, err := os.Create("plot.html")
	if err != nil {
		log.Fatalf("Ошибка при создании plot.html: %+v", err)
	}

	page := components.NewPage()
	page.AddCharts(line)
	err = page.Render(bufio.NewWriter(f))
	if err != nil {
		log.Fatalf("Ошибка при отрисовывании plot.html: %+v", err)
	}
}

func convertToLineData(points []Point) []opts.LineData {
	res := make([]opts.LineData, len(points))
	for _, p := range points {
		res = append(res, opts.LineData{Value: []float64{p.X, p.Y}})
	}
	return res
}

func convertToXAxis(points []Point) []string {
	res := make([]string, len(points))
	for _, p := range points {
		res = append(res, fmt.Sprintf("%10.5f", p.Y))
	}
	return res
}