package main

import (
	"flag"
	"fmt"
	"math/cmplx"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func drawFTPlot(fourier []complex128, sampleRate int, filename string) error {
	p := plot.New()
	p.Title.Text = "Fourier transform"
	p.X.Label.Text = "Frequency, Hz"
	p.Y.Label.Text = "Magnitude"

	pts := make(plotter.XYs, len(fourier)/2)
	for i := range pts {
		pts[i].X = float64(i+1) * float64(sampleRate) / float64(len(fourier))
		pts[i].Y = cmplx.Abs(fourier[i+1]) / float64(len(fourier))
	}
	l, _ := plotter.NewLine(pts)
	p.Add(l)
	if err := p.Save(30*vg.Centimeter, 10*vg.Centimeter, filename); err != nil {
		return err
	}
	return nil
}

func main() {
	var wavFilename, plotFilename, denoisedFilename string
	flag.StringVar(&wavFilename, "f", "", "Path to the WAV audio sample.")
	flag.StringVar(&plotFilename, "p", "fourier.png", "Output file for the DFT plot (PNG).")
	flag.StringVar(&denoisedFilename, "d", "denoised.wav", "Output file for IDFT (WAV).")
	flag.Parse()
	if wavFilename == "" {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	samples, sampleRate, err := readWav(wavFilename)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Sample rate: %d\nNumber of samples: %d\nDuration: %.02fs\n",
		sampleRate, len(samples), float64(len(samples))/float64(sampleRate))

	startTimestamp := time.Now()
	fourier := dft(samples)
	finishTimestamp := time.Now()
	fmt.Println("DFT computation time:", finishTimestamp.Sub(startTimestamp))

	drawFTPlot(fourier, sampleRate, plotFilename)

	peakStart, peakEnd := findPeak(fourier, sampleRate)
	for i := range fourier {
		if peakStart <= i && i <= peakEnd {
			fourier[i] = 0
		}
	}
	startTimestamp = time.Now()
	invFourier := idft(fourier)
	finishTimestamp = time.Now()
	fmt.Println("IDFT computation time:", finishTimestamp.Sub(startTimestamp))

	if err := writeWav(denoisedFilename, invFourier, sampleRate); err != nil {
		fmt.Println(err)
	}
}
