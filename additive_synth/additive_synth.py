from argparse import ArgumentParser
import numpy as np
import os
from scipy.io import wavfile
from typing import List


def generate_sample(sampling_rate: int, freqs: List[float], duration: float) -> np.ndarray:
    num_samples = int(sampling_rate * duration)
    x = np.arange(0, num_samples, dtype=np.float64)
    wave = np.zeros(num_samples, dtype=np.float64)
    for f in freqs:
        wave += np.sin(2*np.pi*f * (x / sampling_rate))
    wave /= np.abs(wave).max()
    return wave


def write_to_wav(filename: str, sampling_rate: int, wave: np.ndarray):
    max_amplitude = np.iinfo(np.int16).max
    data = (max_amplitude * wave).astype(np.int16)
    wavfile.write(filename, sampling_rate, data)


def main():
    parser = ArgumentParser()
    parser.add_argument(
        '-o', '--out',
        dest='wav_filename', type=str, required=True,
        help='Output file name.'
    )
    parser.add_argument(
        '-r', '--rate',
        dest='sampling_rate', type=int, required=True,
        help='Sampling rate.'
    )
    parser.add_argument(
        '-d', '--duration',
        dest='duration', type=float, required=True,
        help='Duration (s).'
    )
    parser.add_argument(
        '-f', '--frequencies',
        dest='freqs', type=float, nargs='+', required=True,
        help='A space-separated list of frequencies (Hz).'
    )
    cmd_args = parser.parse_args()

    wave = generate_sample(cmd_args.sampling_rate, cmd_args.freqs, cmd_args.duration)
    write_to_wav(cmd_args.wav_filename, cmd_args.sampling_rate, wave)


if __name__ == '__main__':
    main()