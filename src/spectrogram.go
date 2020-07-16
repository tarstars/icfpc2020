package main

import (
	// "fmt"
	"github.com/mjibson/go-dsp/wav"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"image/color"
	"os"
)

type Filter struct {
	acoeff []float64
	bcoeff []float64
	gain   float64
	xv, yv []float64
}

func New500Filter() *Filter {
	var filter Filter

	filter.acoeff = []float64{0.9928926106217069, -1.9878505507862085, 1}
	filter.bcoeff = []float64{-1, 0, 1}
	filter.gain = 188.89635546457632
	filter.xv = []float64{0, 0, 0}
	filter.yv = []float64{0, 0, 0}

	return &filter
}

func New600Filter() *Filter {
	var filter Filter

	filter.acoeff = []float64{0.9928953208350709, -1.9856305367036893, 1}
	filter.bcoeff = []float64{-1, 0, 1}
	filter.gain = 181.08092138940904
	filter.xv = []float64{0, 0, 0}
	filter.yv = []float64{0, 0, 0}

	return &filter
}

func (flt *Filter) apply(v float64) float64 {
	var i int
	var out float64
	out = 0
	for i = 0; i < 2; i++ {
		flt.xv[i] = flt.xv[i+1]
	}
	flt.xv[2] = v / flt.gain
	for i = 0; i < 2; i++ {
		flt.yv[i] = flt.yv[i+1]
	}
	for i = 0; i <= 2; i++ {
		out += flt.xv[i] * flt.bcoeff[i]
	}

	for i = 0; i < 2; i++ {
		out -= flt.yv[i] * flt.acoeff[i]
	}
	flt.yv[2] = out

	return out
}

func plotGraph(dataX,
	dataY1 []float64, caption1 string,
	dataY2 []float64, caption2 string,
	dataY3 []float64, caption3 string,
	fileName string) {
	p, _ := plot.New()

	p.Title.Text = "p(t)"
	p.X.Label.Text = "t"
	p.Y.Label.Text = "p"

	n := len(dataX)
	xys1 := make(plotter.XYs, n)
	xys2 := make(plotter.XYs, n)
	xys3 := make(plotter.XYs, n)
	for ind := 0; ind < n; ind += 1 {
		xys1[ind].X = dataX[ind]
		xys1[ind].Y = dataY1[ind]
		xys2[ind].X = dataX[ind]
		xys2[ind].Y = dataY2[ind]
		xys3[ind].X = dataX[ind]
		xys3[ind].Y = dataY3[ind]
	}

	lines2, _, _ := plotter.NewLinePoints(xys2)
	lines2.Color = color.RGBA{B:255, A:255}

	_ = plotutil.AddLines(p, caption1, xys1)
	//_ = plotutil.AddLines(p, caption2, xys2)
	p.Add(lines2)
	_ = plotutil.AddLines(p, caption3, xys3)
	_ = p.Save(800*vg.Inch, 3*vg.Inch, fileName)
}

func plotGraphics() {
	source_wav_fd, _ := os.Open("/home/tass/database/icfpc2020/messages/radio-transmission-recording.wav")
	source_wav, _ := wav.New(source_wav_fd)

	a, _ := source_wav.ReadFloats(source_wav.Samples)
	// f_d := float64(source_wav.SampleRate)
	// f_s := 500

	startInd := 0
	stopInd := 1000000
	n := stopInd - startInd

	f500 := New500Filter()
	f600 := New600Filter()

	dataX := make([]float64, n)
	dataRaw := make([]float64, n)
	data500 := make([]float64, n)
	data600 := make([]float64, n)

	for ind := startInd; ind < stopInd; ind += 1 {
		dataX[ind-startInd] = float64(ind)
		dataRaw[ind-startInd] = float64(a[ind])
		data500[ind-startInd] = f500.apply(float64(a[ind]))
		data600[ind-startInd] = f600.apply(float64(a[ind]))
	}

	plotGraph(dataX, dataRaw,"raw", data500, "500", data600, "600", "w.svg")
}

func main() {
	plotGraphics()
}
