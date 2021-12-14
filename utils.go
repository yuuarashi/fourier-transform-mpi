package main

import (
	"errors"
	"flag"
	"math/cmplx"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func parseCmdArgs() (string, string, string, error) {
	var wavFilename, plotFilename, denoisedFilename string
	flag.StringVar(&wavFilename, "f", "", "Path to the WAV audio sample.")
	flag.StringVar(&plotFilename, "p", "fourier.png", "Output file for the DFT plot (PNG).")
	flag.StringVar(&denoisedFilename, "d", "denoised.wav", "Output file for IDFT (WAV).")
	flag.Parse()
	if wavFilename == "" {
		return "", "", "", errors.New("No audio sample provided.")
	}

	return wavFilename, plotFilename, denoisedFilename, nil
}

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

func removeHighMagnitudeFreq(fourier []complex128, sampleRate int) []complex128 {
	maxMgn, maxIdx := 0.0, 0
	for i := range fourier {
		if i > len(fourier)/2 {
			break
		}

		if i > 0 && cmplx.Abs(fourier[i]) > maxMgn {
			maxMgn = cmplx.Abs(fourier[i])
			maxIdx = i
		}
	}

	windowSize := int(50.0 * float64(len(fourier)) / float64(sampleRate))
	var startIdx, endIdx int
	if maxIdx > windowSize {
		startIdx = maxIdx - windowSize
	} else {
		startIdx = 1
	}
	endIdx = maxIdx + windowSize

	for i := startIdx; i <= endIdx; i++ {
		fourier[i] = 0
		fourier[len(fourier)-i-1] = 0
	}

	return fourier
}
