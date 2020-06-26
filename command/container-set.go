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
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	serviceArg := c.flagSet.String("service", "", "select the service")
	nInstancesArg := containerInstancesFlag(c.flagSet)
	runCmdArg := containerRunCmdFlag(c.flagSet)
	limitMemoryArg := containerLimitMemoryFlag(c.flagSet)
	limitCPUArg := containerLimitCPUFlag(c.flagSet)
	limitNetArg := containerLimitNetFlag(c.flagSet)
	noRunCmdArg := containerNoRunCmdFlag(c.flagSet)
	envCmdArg := envFileFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	if *serviceArg == "" {
		return c.errorWithUsage(errors.New("Service name cannot be empty."))
	}

	if *noRunCmdArg && *runCmdArg != "" {
		return c.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *nInstancesArg <= 0 {
		if *nInstancesArg != -1 {
			c.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set.")
		}

		if *runCmdArg == "" && !*noRunCmdArg && *limitCPUArg < 0 && *limitMemoryArg < 0 && *limitNetArg < 0 {
			err := errors.New("Invalid values provided for instance number.")
			return c.errorWithUsage(err)
		}
	}

	return c.runWithSpinner("configure service", endpoint.String(), func(client *squarescale.Client) (string, error) {
		container, err := client.GetServicesInfo(*projectUUID, *serviceArg)
		if err != nil {
			return "", err
		}

		c.Meta.spin.Stop()
		c.Meta.info("")

		if *nInstancesArg > 0 {
			container.Size = *nInstancesArg
			c.info("Configure service with %d instances", *nInstancesArg)
		}

		if *runCmdArg != "" {
			c.info("Configure service with run command: %s", *runCmdArg)
			container.RunCommand, err = shellquote.Split(*runCmdArg)
			if err != nil {
				return "", err
			}
			if container.RunCommand == nil {
				container.RunCommand = []string{}
			}
		}
		if *noRunCmdArg {
			c.info("Configure service without run command")
			container.RunCommand = []string{}
		}

		if *limitMemoryArg >= 0 {
			c.info("Configure service with memory limit of %d MB", *limitMemoryArg)
			container.Limits.Memory = *limitMemoryArg
		}
		if *limitCPUArg >= 0 {
			c.info("Configure service with CPU limit of %d Mhz", *limitCPUArg)
			container.Limits.CPU = *limitCPUArg
		}
		if *limitNetArg >= 0 {
			c.info("Configure service with network bandwidth limit of %d Mbps", *limitNetArg)
			container.Limits.Net = *limitNetArg
		}
		if *envCmdArg != "" {
			c.info("Configure service with some env")
			err := container.SetEnv(*envCmdArg)
			if err != nil {
				c.error(err)
			}
		}

		msg := fmt.Sprintf(
			"Successfully configured container '%s' for project '%s'",
			*serviceArg, *projectUUID)

		c.Meta.spin.Start()
		return msg, client.ConfigService(container)
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
  Containers are specified using their given name.

Example:
  sqsc container set                \
      -project="my-rails-project"   \
      -container="my-name/my-repo"  \
      -instances=42                 \
      -e env.json

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
