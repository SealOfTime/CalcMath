package main

import (
	"bufio"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"log"
	"math/rand"
	"os"
)

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func main() {
	ds := readDataSetFromStd()
	step := ds.X[1] - ds.X[0]
	newt := NewtonPolynomialFinDiff{
		h:         step,
		DataSet:	ds,
	}
	lag := LagrangePolynomial{
		DataSet: ds,
	}

	fmt.Println("Формируем график...")
	savePlot(&newt, &lag, &ds)

	var(
		x float64
		err error
	)
	for {
		fmt.Print("Введите X: ")
		_, err = fmt.Scanf("%f\n", &x)
		if err != nil {
			fmt.Println("Ошибка ввода.")
			continue
		}
		fmt.Printf("Формулой Ньютона: F(%f)≈%f\n", x, newt.EvaluateAt(x))
		fmt.Printf("Многочленом Лагранжа: F(%f)≈%f\n", x, lag.EvaluateAt(x))
	}

	//servePlot()
}

type Solver interface {
	EvaluateAt(x float64) float64
}

func savePlot(newt *NewtonPolynomialFinDiff, lag *LagrangePolynomial, ds *DataSet) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Plot",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:          true,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:        "value",
			Show:        true,
			SplitNumber: 10,
			Scale:       true,
			Min:         ds.X[0],
			Max:         ds.X[len(ds.X)-1],
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

	st := (ds.X[len(ds.X)-1] - ds.X[0])/20
	newtSet := interpolatedPlot(ds.X[0], ds.X[len(ds.X)-1], st, newt)
	lagSet := interpolatedPlot(ds.X[0], ds.X[len(ds.X)-1], st, lag)

	line.AddSeries("Узлы интерполяции", nodesPlot(ds), charts.WithLineChartOpts(opts.LineChart{
		Smooth: false,
	}))
	line.AddSeries("Формула Ньютона", newtSet, charts.WithLineChartOpts(opts.LineChart{
		Smooth: false,
	}))
	line.AddSeries("Полином Лагранжа", lagSet, charts.WithLineChartOpts(opts.LineChart{
		Smooth: false,
	}))

	f, err := os.Create("plot.html")
	if err != nil {
		log.Fatalf("Ошибка при создании plot.html: %+v", err)
	}

	page := components.NewPage()
	page.AddCharts(line)
	r := bufio.NewWriter(f)
	err = page.Render(r)
	if err != nil {
		log.Fatalf("Ошибка при отрисовывании plot.html: %+v", err)
	}

	err = r.Flush()
	if err != nil {
		log.Fatalf("Ошибка при отрисовывании plot.html: %+v", err)
	}
}

func nodesPlot(ds *DataSet) []opts.LineData{
	res := make([]opts.LineData, len(ds.X))
	for i, x := range ds.X {
		res[i] = opts.LineData{
			Value: []float64{x, ds.Y[i]},
		}
	}
	return res
}
func interpolatedPlot(start, end, step float64, s Solver) []opts.LineData {
	res := make([]opts.LineData, 0, int((end-start)/step))
	for x := start; x <= end; x += step{
		res = append(res, opts.LineData{Value: []float64{x, s.EvaluateAt(x)}})
	}
	return res
}