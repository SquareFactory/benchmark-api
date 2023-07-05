package main

import (
	"log"
	"os"

	"github.com/squarefactory/benchmark-api/cmd/run"
	"github.com/urfave/cli/v2"
)

var version = "dev"

var flags = []cli.Flag{}

var app = &cli.App{
	Name:    "hpl-ai",
	Usage:   "Run an HPL-AI benchmark",
	Version: version,
	Flags:   flags,
	Commands: []*cli.Command{
		run.Command,
	},
	Suggest: true,
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("app crashed, err: %s", err)
	}
}
