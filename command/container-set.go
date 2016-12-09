package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ContainerSetCommand allows to configure a project container.
type ContainerSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ContainerSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectArg := projectFlag(c.flagSet)
	containerArg := containerFlag(c.flagSet)
	nInstancesArg := containerInstancesFlag(c.flagSet)
	updateCmdArg := containerUpdateCmdFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
	}

	if *containerArg == "" {
		return c.errorWithUsage(errors.New("Container short name cannot be empty."))
	}

	if *nInstancesArg <= 0 {
		if *nInstancesArg != -1 {
			c.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set.")
		}

		if *updateCmdArg == "" {
			err := errors.New("Invalid values provided for instance number and update command.")
			return c.errorWithUsage(err)
		}
	}

	var msg string
	err := c.runWithSpinner("configure container", *endpoint, func(token string) error {
		container, err := squarescale.GetContainerInfo(*endpoint, token, *projectArg, *containerArg)
		if err != nil {
			return err
		}

		if *nInstancesArg > 0 {
			container.Size = *nInstancesArg
		}

		if *updateCmdArg != "" {
			container.Command = *updateCmdArg
		}

		msg = fmt.Sprintf(
			"Successfully configured container (instances = '%d', command = '%s') '%s' for project '%s'",
			container.Size, container.Command, *containerArg, *projectArg)

		return squarescale.ConfigContainer(*endpoint, token, container)
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
}

// Synopsis is part of cli.Command implementation.
func (c *ContainerSetCommand) Synopsis() string {
	return "Set container runtime parameters for a Squarescale project"
}

// Help is part of cli.Command implementation.
func (c *ContainerSetCommand) Help() string {
	helpText := `
usage: sqsc container set [options]

  Set container runtime parameters for a Squarescale project.
  '-update' options is an update command which will be run at
  each build for the specified container, eg "rake db:migrate".
  Containers are specified using the form '${USER}/${REPOSITORY}'

Example:
  sqsc container set                \
      --project="my-rails-project"  \
      --container="my-name/my-repo" \
      --instances=42                \
      --update="rake db:migrate"

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}