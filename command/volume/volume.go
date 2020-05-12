package volume

import (
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command/flags"
)

func New() *cmd {
	return &cmd{}
}

type cmd struct{}

func (c *cmd) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return flags.Usage(help, nil)
}

const synopsis = "Interact with the key-value store"
const help = `
Usage: consul volume <subcommand> [options] [args]

`
