package main

import (
	"os"

	"github.com/squarefactory/benchmark-api/cmd/run"
	"github.com/squarefactory/benchmark-api/logger"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var version = "dev"

var flags = []cli.Flag{}

var app = &cli.App{
	Name:    "hpl-ai",
	Usage:   "",
	Version: version,
	Flags:   flags,
	Commands: []*cli.Command{
		run.Command,
	},
	Suggest: true,
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.I.Fatal("app crashed", zap.Error(err))
	}
}
