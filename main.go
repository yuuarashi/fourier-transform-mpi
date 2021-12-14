package main

import (
	"fmt"
	"time"

	mpi "github.com/sbromberger/gompi"
)

const MASTER_RANK int = 0

func main() {
	mpi.Start()
	rank := mpi.WorldRank()
	worldComm := mpi.NewCommunicator(nil)

	if rank == MASTER_RANK {
		wavFilename, plotFilename, denoisedFilename, err := parseCmdArgs()
		if err != nil {
			fmt.Println(err)
			worldComm.Abort(-1)
			return
		}

		samples, sampleRate, err := readWav(wavFilename)
		if err != nil {
			fmt.Println(err)
			worldComm.Abort(-1)
			return
		}

		fmt.Printf("Sample rate: %d\nNumber of samples: %d\nDuration: %.02fs\n",
			sampleRate, len(samples), float64(len(samples))/float64(sampleRate))

		startTimestamp := time.Now()
		fourier := dft(worldComm, samples)
		finishTimestamp := time.Now()
		fmt.Println("DFT computation time:", finishTimestamp.Sub(startTimestamp))

		drawFTPlot(fourier, sampleRate, plotFilename)
		fourier = removeHighMagnitudeFreq(fourier, sampleRate)

		startTimestamp = time.Now()
		invFourier := idft(worldComm, fourier)
		finishTimestamp = time.Now()
		fmt.Println("IDFT computation time:", finishTimestamp.Sub(startTimestamp))

		if err := writeWav(denoisedFilename, invFourier, sampleRate); err != nil {
			fmt.Println(err)
		}
	} else {
		dft(worldComm, nil)
		idft(worldComm, nil)
	}

	mpi.Stop()
}
