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
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	container := containerFlag(c.flagSet)
	remove := c.flagSet.Bool("remove", false, "Remove the key from environment variables")

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	var key, value string
	args = c.flagSet.Args()
	switch len(args) {
	case 0:
		return c.errorWithUsage(fmt.Errorf("No argument provided"))
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

	return c.runWithSpinner("set environment variable", endpoint.String(), func(client *squarescale.Client) (string, error) {
		env, err := squarescale.NewEnvironment(client, *projectUUID)
		if err != nil {
			return "", err
		}

		var msg string

		if *container != "" {
			service, err := env.GetServiceGroup(*container)
			if err != nil {
				return "", err
			}

			if *remove {
				if err = service.RemoveVariable(key); err != nil {
					return "", err
				}
				msg = fmt.Sprintf("Successfully removed variable '%s' for container '%s'",
					key, *container)
			} else {
				service.SetVariable(key, value)
				msg = fmt.Sprintf(
					"Successfully set variable '%s' to value '%s' for container '%s'",
					key, value, *container)
			}
		} else if *remove {
			if err = env.Project.RemoveVariable(key); err != nil {
				return "", err
			}
			msg = fmt.Sprintf("Successfully removed global variable '%s'", key)
		} else {
			env.Project.SetVariable(key, value)
			msg = fmt.Sprintf("Successfully set global variable '%s' to value '%s'", key, value)
		}

		return msg, env.CommitEnvironment(client, *projectUUID)
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
