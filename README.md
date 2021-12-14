# Parallel Fourier transform in Go using MPI

### Launching the program
1. Install [OpenMPI](https://www.open-mpi.org/) on your system (tested with OpenMPI v4.1.2).
2. Create an executable with `go build .`
3. Run the executable with `mpirun` (Consult the [OpenMPI documentation](https://www.open-mpi.org/doc/v4.1/man1/mpirun.1.php) for detail).

### Usage
`fourier-transform-mpi -f INPUT_WAV_FILE -p OUTPUT_PLOT_FILE -d DENOISED_WAV_FILE`

- `-f` : Path to the input WAV file.
- `-p` : Path to the output PNG file for the frequency domain plot.
- `-d` : Path to the denoised output WAV file where the highest magnitude frequency is removed with its ±50Hz neighbourhood.

### Generating toy examples
A simple additive generator is provided for creating samples by summing sine waves at different frequencies with no phase shift.
Install the dependencies with `pip install -r requirements.txt`.

Usage: `additive_synth.py -o OUTPUT_WAV_FILE -r SAMPLING_RATE -d DURATION_SEC -f FREQ1 FREQ2 ...`
