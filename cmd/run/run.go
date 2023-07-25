package run

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/squarefactory/benchmark-api/benchmark"
	"github.com/squarefactory/benchmark-api/executor"
	"github.com/squarefactory/benchmark-api/resultparser"
	"github.com/squarefactory/benchmark-api/scheduler"
	"github.com/squarefactory/benchmark-api/try"
	"github.com/urfave/cli/v2"
)

const (
	user = "root"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "container.path",
		Value: "/etc/hpl-benchmark/hpc-benchmarks:hpl.sqsh",
		EnvVars: []string{
			"CONTAINER_PATH",
		},
		Aliases: []string{"c"},
		Action: func(ctx *cli.Context, s string) error {
			info, err := os.Stat(s)
			if err != nil {
				return err
			}
			perms := info.Mode().Perm()
			if perms&0o077 != 0 {
				log.Fatal(
					"incorrect permissions for container .sqsh, must be user-only",
				)
			}
			return nil
		},
	},
}

var Command = &cli.Command{
	Name:      "run",
	Usage:     "Run an HPL-AI benchmark.",
	Flags:     flags,
	ArgsUsage: "<node_number>",
	Action: func(cCtx *cli.Context) error {

		ctx := cCtx.Context
		if cCtx.NArg() < 1 {
			return errors.New("not enough arguments")
		}

		arg := cCtx.Args().Get(0)
		node, err := strconv.Atoi(arg)
		if err != nil {
			log.Printf("Failed to convert %s to integer: %s", arg, err)
			return err
		}

		containerPath := os.Getenv("CONTAINER_PATH")
		workspace := filepath.Dir(containerPath)
		firstSet := benchmark.NewBenchmark(
			benchmark.DATParams{},
			benchmark.SBATCHParams{
				Node:          node,
				ContainerPath: containerPath,
				Workspace:     workspace,
			},
			scheduler.NewSlurm(&executor.Shell{}, user),
		)

		log.Printf("running first set, with general parameters")
		if err := RunFirstSet(firstSet, ctx); err != nil {
			log.Printf("failed to run first set of benchmark: %s", err)
			return err
		}

		log.Printf("first set finished running, processing results")

		// Get optimal benchmark DAT parameters
		optimalParams, err := ProcessFirstSet()
		if err != nil {
			log.Printf("failed to process first set: %s", err)
			return err
		}

		optimalSet := benchmark.NewBenchmark(
			benchmark.DATParams{
				ProblemSize: optimalParams.ProblemSize,
				P:           optimalParams.P,
				Q:           optimalParams.Q,
			},
			benchmark.SBATCHParams{
				Node:          node,
				ContainerPath: containerPath,
				Workspace:     workspace,
			},
			scheduler.NewSlurm(&executor.Shell{}, user),
		)

		log.Printf("running second set, with optimal parameters")
		if err := RunOptimalSet(optimalSet, ctx); err != nil {
			log.Printf("failed to run second set of benchmark: %s", err)
			return err
		}

		return nil
	},
}

func RunFirstSet(b *benchmark.Benchmark, ctx context.Context) error {

	if err := b.CalculateBenchmarkParams(ctx); err != nil {
		log.Printf("failed to calculate first set parameters")
		return err
	}

	files, err := b.GenerateFiles(ctx)
	if err != nil {
		log.Printf("Failed to generate benchmark files: %s", err)
		return err
	}

	output, err := os.Create(scheduler.JobOutput)
	if err != nil {
		log.Printf("failed to create output file: %s", err)
		return err
	}
	defer output.Close()

	if err := b.Run(ctx, &files); err != nil {
		log.Printf("Failed to run benchmark: %s", err)
		return err
	}

	_, err = try.Do(func() (int, error) {
		_, err := b.SlurmClient.FindRunningJobByName(
			ctx,
			&scheduler.FindRunningJobByNameRequest{
				Name: benchmark.JobName,
				User: benchmark.User,
			},
		)
		if err == nil {
			log.Print("benchmark is still running, unable to process results")
			return 0, errors.New("benchmark is still running")
		}

		return 0, nil
	}, 60, 5*time.Minute)

	if err != nil {
		log.Printf("Benchmark is still running, unable to process results")
		return err
	}

	return nil
}

func ProcessFirstSet() (benchmark.DATParams, error) {

	if err := resultparser.WriteResultsToCSV(scheduler.JobOutput); err != nil {
		log.Printf("Failed to process results: %s", err)
		return benchmark.DATParams{}, err
	}

	optimalRow, err := resultparser.FindMaxGflopsRow(resultparser.CsvFile)
	if err != nil {
		log.Printf("Failed to find row containing max gflops score: %s", err)
		return benchmark.DATParams{}, err
	}

	p, err := strconv.Atoi(optimalRow[2])
	if err != nil {
		log.Printf("failed to convert %s as integer: %s", optimalRow[2], err)
		return benchmark.DATParams{}, err
	}

	q, err := strconv.Atoi(optimalRow[3])
	if err != nil {
		log.Printf("failed to convert %s as integer: %s", optimalRow[2], err)
		return benchmark.DATParams{}, err
	}

	return benchmark.DATParams{
		NProblemSize: 1,
		ProblemSize:  optimalRow[0],
		NBlockSize:   1,
		BlockSize:    optimalRow[1],
		P:            p,
		Q:            q,
	}, nil
}

func RunOptimalSet(b *benchmark.Benchmark, ctx context.Context) error {

	if err := b.CalculateSBATCHParams(ctx); err != nil {
		log.Printf("failed to calculate sbatch params for optimal set: %s", err)
		return err
	}

	files, err := b.GenerateFiles(ctx)
	if err != nil {
		log.Printf("Failed to generate benchmark files: %s", err)
		return err
	}

	output, err := os.Create(scheduler.JobOutput)
	if err != nil {
		log.Printf("failed to create output file: %s", err)
		return err
	}
	defer output.Close()

	if err := b.Run(ctx, &files); err != nil {
		log.Printf("Failed to run benchmark: %s", err)
		return err
	}

	return nil
}
