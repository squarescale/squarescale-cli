package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

// RedisCommand is a cli.Command implementation for top level `sqsc Redis` command.
type RedisCommand struct {
}

// Run is part of cli.Command implementation.
func (c *RedisCommand) Run(args []string) int {
	return cli.RunResultHelp
}

// Synopsis is part of cli.Command implementation.
func (c *RedisCommand) Synopsis() string {
	return "Commands related to Redis"
}

// Help is part of cli.Command implementation.
func (c *RedisCommand) Help() string {
	helpText := `
usage: sqsc Redis <subcommand>
  Run a project Redis related command.

  List of supported subcommands is available below.

`
	return strings.TrimSpace(helpText)
}
