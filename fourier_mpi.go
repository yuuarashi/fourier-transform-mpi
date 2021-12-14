package main

import (
	"math"
	"math/cmplx"

	mpi "github.com/sbromberger/gompi"
)

func dft(comm mpi.Communicator, x []complex128) []complex128 {
	return mpiDftProcedure(comm, x, false)
}

func idft(comm mpi.Communicator, x []complex128) []complex128 {
	return mpiDftProcedure(comm, x, true)
}

func mpiDftProcedure(comm mpi.Communicator, x []complex128, inverse bool) []complex128 {
	rank := comm.Rank()
	var startIdx, numSamples int
	var data []complex128
	if rank == MASTER_RANK {
		numSamples = len(x)
		samplesPerProc := int64(numSamples / comm.Size())

		indexInfo := make([]int64, 2)
		indexInfo[0] = int64(numSamples)
		indexInfo[1] = 0
		for pId := 1; pId < comm.Size(); pId++ {
			comm.SendInt64s(indexInfo, pId, 42)
			comm.SendComplex128s(x[indexInfo[1]:indexInfo[1]+samplesPerProc], pId, 42)
			indexInfo[1] += samplesPerProc
		}
		startIdx = int(indexInfo[1])
		data = x[startIdx:]
	} else {
		indexInfo := comm.RecvInt64s(MASTER_RANK, 42)
		numSamples = int(indexInfo[0])
		startIdx = int(indexInfo[1])
		data = comm.RecvComplex128s(MASTER_RANK, 42)
	}

	y := make([]complex128, numSamples)
	partialSums := make([]complex128, numSamples)
	scalingFactor, exponentSign := 1.0, -1.0
	if inverse {
		scalingFactor = 1.0 / float64(numSamples)
		exponentSign = 1.0
	}
	for k := range partialSums {
		partialSums[k] = 0
		for i := range data {
			pow := complex(0, exponentSign*2.0*math.Pi*float64(k)*float64(startIdx+i)/float64(numSamples))
			partialSums[k] += data[i] * cmplx.Exp(pow) * complex(scalingFactor, 0)
		}
	}

	comm.ReduceComplex128s(y, partialSums, mpi.OpSum, MASTER_RANK)

	return y
}
