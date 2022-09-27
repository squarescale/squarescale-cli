package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceSetCommand allows to configure a project service.
type ServiceSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ServiceSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	serviceArg := serviceFlag(c.flagSet)
	nInstancesArg := containerInstancesFlag(c.flagSet)
	runCmdArg := containerRunCmdFlag(c.flagSet)
	entrypoint := c.flagSet.String("entrypoint", "", "This is the script / program that will be executed")
	limitMemoryArg := containerLimitMemoryFlag(c.flagSet)
	limitCPUArg := containerLimitCPUFlag(c.flagSet)
	limitNetArg := containerLimitNetFlag(c.flagSet)
	noRunCmdArg := containerNoRunCmdFlag(c.flagSet)
	envCmdArg := envFileFlag(c.flagSet)
	schedulingGroupsArg := containerSchedulingGroupsFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
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
		var UUID string
		var err error
		var projectToShow string
		if *projectUUID == "" {
			projectToShow = *projectName
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			projectToShow = *projectUUID
			UUID = *projectUUID
		}

		container, err := client.GetServicesInfo(UUID, *serviceArg)
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
			container.RunCommand = *runCmdArg
			if len(container.RunCommand) == 0 {
				container.RunCommand = ""
			}
		}
		if *noRunCmdArg {
			c.info("Configure service without run command")
			container.RunCommand = ""
		}

		if *entrypoint != "" {
			c.info("Configure service with entrypoint: %s", *entrypoint)
			container.Entrypoint = *entrypoint
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

		schedulingGroupsToAdd := getSchedulingGroupsArray(UUID, client, *schedulingGroupsArg)
		if len(schedulingGroupsToAdd) != 0 {
			c.info("Configure scheduling groups")
			container.SchedulingGroups = schedulingGroupsToAdd
		}

		msg := fmt.Sprintf(
			"Successfully configured service '%s' for project '%s'",
			*serviceArg, projectToShow)

		c.Meta.spin.Start()
		return msg, client.ConfigService(container)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ServiceSetCommand) Synopsis() string {
	return "Set service aka Docker container runtime parameters for project"
}

// Help is part of cli.Command implementation.
func (c *ServiceSetCommand) Help() string {
	helpText := `
usage: sqsc service set [options]

  Set service aka Docker container runtime parameters for project.
  Services are specified using their given name.

Example:
  sqsc service set                \
      -project="my-rails-project" \
      -service="my-name/my-repo"  \
      -instances=42               \
      -env=env.json
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}