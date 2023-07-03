package run

import "github.com/urfave/cli/v2"

var flags = []cli.Flag{}

var Command = &cli.Command{
	Name:      "run",
	Usage:     "Run an HPL-AI benchmark.",
	Flags:     flags,
	ArgsUsage: "<hostnames>",
	Action: func(cCtx *cli.Context) error {
		return nil
	},
}
