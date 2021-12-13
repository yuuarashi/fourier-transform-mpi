package main

import (
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

func readWav(filename string) ([]complex128, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	dec := wav.NewDecoder(file)
	buf, err := dec.FullPCMBuffer()
	if err != nil {
		return nil, 0, err
	}
	file.Close()

	samples := make([]complex128, len(buf.Data))
	maxAmpl := 1<<(buf.SourceBitDepth-1) - 1
	for i := range samples {
		samples[i] = complex(float64(buf.Data[i])/float64(maxAmpl), 0)
	}

	return samples, int(dec.SampleRate), nil
}

func writeWav(filename string, samples []complex128, sampleRate int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	bitDepth, numChans, wavFormatPCM := 16, 1, 1
	format := audio.Format{NumChannels: numChans, SampleRate: sampleRate}
	bufData := make([]int, len(samples))
	maxAmpl := 1<<(bitDepth-1) - 1
	for i := range bufData {
		bufData[i] = int(real(samples[i]) * float64(maxAmpl))
	}
	buf := audio.IntBuffer{Format: &format, SourceBitDepth: bitDepth, Data: bufData}

	enc := wav.NewEncoder(file, sampleRate, bitDepth, numChans, wavFormatPCM)
	if err := enc.Write(&buf); err != nil {
		return err
	}
	if err := enc.Close(); err != nil {
		return err
	}
	file.Close()

	return nil
}
