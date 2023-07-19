# benchmark-cli

A CLI tool which automates the benchmarking of clusters with HPL-AI.

## Build

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o benchmark ./cmd
```

## Usage

The path to the .sqsh container image for HPL Benchmark is set as an environment variable:

```sh
export CONTAINER_PATH="$(pwd)/hpc-benchmarks:21.4-hpl.sqsh"
```

Then, you can launch a benchmark by using the run command. Example for a single node benchmark:

```sh
./benchmark run 1
```

The results are exported in the benchmark.csv file, in the same directory as the executable.
