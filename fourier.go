package main

import (
	"math"
	"math/cmplx"
)

func dft(x []complex128) []complex128 {
	X := make([]complex128, len(x))
	for k := range x {
		X[k] = 0
		for i := range x {
			pow := complex(0, -2.0*math.Pi*float64(k)*float64(i)/float64(len(x)))
			X[k] += x[i] * cmplx.Exp(pow)
		}
	}
	return X
}

func findPeak(fourier []complex128, sampleRate int) (int, int) {
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

	return startIdx, endIdx
}

func idft(X []complex128) []complex128 {
	x := make([]complex128, len(X))
	for k := range X {
		x[k] = 0
		for i := range X {
			pow := complex(0, 2.0*math.Pi*float64(k)*float64(i)/float64(len(X)))
			x[k] += X[i] * cmplx.Exp(pow) / complex(float64(len(X)), 0)
		}
	}
	return x
}
