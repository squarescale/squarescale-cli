package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/kballard/go-shellquote"
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
	buildServiceArg := containerBuildServiceFlag(c.flagSet)
	updateCmdArg := containerUpdateCmdFlag(c.flagSet)
	noUpdateCmdArg := containerNoUpdateCmdFlag(c.flagSet)
	runCmdArg := containerRunCmdFlag(c.flagSet)
	limitMemoryArg := containerLimitMemoryFlag(c.flagSet)
	limitCPUArg := containerLimitCPUFlag(c.flagSet)
	limitNetArg := containerLimitNetFlag(c.flagSet)
	noRunCmdArg := containerNoRunCmdFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if err := validateProjectName(*projectArg); err != nil {
		return c.errorWithUsage(err)
	}

	if *containerArg == "" {
		return c.errorWithUsage(errors.New("Container short name cannot be empty."))
	}

	if *noUpdateCmdArg && *updateCmdArg != "" {
		return c.errorWithUsage(errors.New("Cannot specify an update command and disable it at the same time"))
	}

	if *noRunCmdArg && *runCmdArg != "" {
		return c.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *nInstancesArg <= 0 {
		if *nInstancesArg != -1 {
			c.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set.")
		}

		if *updateCmdArg == "" && *buildServiceArg == "" && *runCmdArg == "" && !*noRunCmdArg && !*noUpdateCmdArg && *limitCPUArg < 0 && *limitMemoryArg < 0 && *limitNetArg < 0 {
			err := errors.New("Invalid values provided for instance number and update command.")
			return c.errorWithUsage(err)
		}
	}

	return c.runWithSpinner("configure container", endpoint.String(), func(client *squarescale.Client) (string, error) {
		container, err := client.GetContainerInfo(*projectArg, *containerArg)
		if err != nil {
			return "", err
		}

		if *nInstancesArg > 0 {
			container.Size = *nInstancesArg
		}

		if *updateCmdArg != "" {
			container.PreCommand, err = shellquote.Split(*updateCmdArg)
			if err != nil {
				return "", err
			}
			if container.PreCommand == nil {
				container.PreCommand = []string{}
			}
		}
		if *noUpdateCmdArg {
			container.PreCommand = []string{}
		}
		if *runCmdArg != "" {
			container.RunCommand, err = shellquote.Split(*runCmdArg)
			if err != nil {
				return "", err
			}
			if container.RunCommand == nil {
				container.RunCommand = []string{}
			}
		}
		if *noRunCmdArg {
			container.RunCommand = []string{}
		}

		if *buildServiceArg != "" {
			container.BuildService = *buildServiceArg
		}

		if *limitMemoryArg >= 0 {
			container.Limits.Memory = *limitMemoryArg
		}
		if *limitCPUArg >= 0 {
			container.Limits.CPU = *limitCPUArg
		}
		if *limitNetArg >= 0 {
			container.Limits.Net = *limitNetArg
		}

		msg := fmt.Sprintf(
			"Successfully configured container (instances = '%d', pre_command = %v, run_command = %v, build_service = %s) '%s' for project '%s'",
			container.Size, container.PreCommand, container.RunCommand, container.BuildService, *containerArg, *projectArg)

		return msg, client.ConfigContainer(container)
	})
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
