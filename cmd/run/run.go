package run

import (
	"errors"
	"log"
	"strconv"

	"github.com/squarefactory/benchmark-api/benchmark"
	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{}

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

		b := benchmark.NewBenchmark(benchmark.DATParams{}, benchmark.SBATCHParams{Node: node})
		files, err := b.GenerateFiles(ctx)
		if err != nil {
			log.Printf("Failed to generate benchmark files: %s", err)
			return err
		}

		if err := b.Run(ctx, &files); err != nil {
			log.Printf("Failed to run benchmark: %s", err)
			return err
		}
		return nil
	},
}
