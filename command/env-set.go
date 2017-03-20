package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// EnvSetCommand sets or remove an environment variable from a project.
type EnvSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *EnvSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	project := projectFlag(c.flagSet)
	remove := c.flagSet.Bool("remove", false, "Remove the key from environment variables")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	var key, value string
	args = c.flagSet.Args()
	switch len(args) {
	case 1:
		key = args[0]
	case 2:
		key, value = args[0], args[1]
	default:
		return c.errorWithUsage(fmt.Errorf("Extra arguments on the command line: %v", args[2:]))
	}

	if err := validateEnvVariable(key, value, *remove); err != nil {
		return c.errorWithUsage(err)
	}

	return c.runWithSpinner("set environment variable", *endpoint, func(client *squarescale.Client) (string, error) {
		vars, err := client.CustomEnvironmentVariables(*project)
		if err != nil {
			return "", err
		}

		var msg string
		if *remove {
			msg = fmt.Sprintf("Successfully removed variable '%s'", key)
			delete(vars, key)
		} else {
			msg = fmt.Sprintf("Successfully set variable '%s' to value '%s'", key, value)
			vars[key] = value
		}

		return msg, client.SetEnvironmentVariables(*project, vars)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *EnvSetCommand) Synopsis() string {
	return "Set an environment variable for a project"
}

// Help is part of cli.Command implementation.
func (c *EnvSetCommand) Help() string {
	helpText := `
usage: sqsc env set [options] <key> <value>

  Set an environment variable for a project. The variable
  must be of the form "<key> <value>" where <key> and <value>
  are both strings. When "--remove" is specified, only the
  key to remove is required.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func validateEnvVariable(key, value string, remove bool) error {
	if key == "" {
		return errors.New("Environment variable key cannot be empty")
	}

	if !remove && value == "" {
		return errors.New("Environment variable value cannot be empty")
	}

	return nil
}
