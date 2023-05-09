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
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	serviceName := serviceFlag(c.flagSet)
	instances := containerInstancesFlag(c.flagSet)
	runCommand := containerRunCmdFlag(c.flagSet)
	entrypoint := entrypointFlag(c.flagSet)
	limitMemory := containerLimitMemoryFlag(c.flagSet)
	limitCPU := containerLimitCPUFlag(c.flagSet)
	noRunCommand := containerNoRunCmdFlag(c.flagSet)
	envVariables := envFileFlag(c.flagSet)
	var envParams []string
	envParameterFlag(c.flagSet, &envParams)
	schedulingGroups := containerSchedulingGroupsFlag(c.flagSet)
	dockerCapabilities := dockerCapabilitiesFlag(c.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(c.flagSet)
	autostart := autostartFlag(c.flagSet)
	dockerDevices := dockerDevicesFlag(c.flagSet)
	maxClientDisconnect := maxclientdisconnect(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *serviceName == "" {
		return c.errorWithUsage(errors.New("Service name cannot be empty."))
	}

	if *noRunCommand && *runCommand != "" {
		return c.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *instances <= 0 {
		if *instances != -1 {
			c.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set.")
		}

		if *runCommand == "" && !*noRunCommand && *limitCPU < 0 && *limitMemory < 0 {
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

		container, err := client.GetServicesInfo(UUID, *serviceName)
		if err != nil {
			return "", err
		}

		c.Meta.spin.Stop()
		c.Meta.info("")

		if *instances > 0 {
			container.Size = *instances
			c.info("Configure service with %d instances", *instances)
		}

		if *runCommand != "" {
			c.info("Configure service with run command: %s", *runCommand)
			container.RunCommand = *runCommand
			if len(container.RunCommand) == 0 {
				container.RunCommand = ""
			}
		}
		if *noRunCommand {
			c.info("Configure service without run command")
			container.RunCommand = ""
		}

		if *entrypoint != "" {
			c.info("Configure service with entrypoint: %s", *entrypoint)
			container.Entrypoint = *entrypoint
		}

		if *limitMemory >= 0 {
			c.info("Configure service with memory limit of %d MB", *limitMemory)
			container.Limits.Memory = *limitMemory
		}
		if *limitCPU >= 0 {
			c.info("Configure service with CPU limit of %d Mhz", *limitCPU)
			container.Limits.CPU = *limitCPU
		}

		if *envVariables != "" {
			c.info("Configure service with some environment JSON file")
			err := container.SetEnv(*envVariables)
			if err != nil {
				c.error(err)
			}
		}
		if len(envParams) > 0 {
			c.info("Configure service with some environment parameters")
			err := container.SetEnvParams(envParams)
			if err != nil {
				c.error(err)
			}
		}

		if isFlagPassed("auto-start", c.flagSet) {
			c.info("Configure service with CPU autostartFlag %v", *autostart)
			container.AutoStart = *autostart
		}

		if isFlagPassed("docker-devices", c.flagSet) {
			c.info("Configure service with custom mapping")
			dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
			if err != nil {
				c.error(err)
			} else {
				container.DockerDevices = dockerDevicesArray
			}
		}

		if *noDockerCapabilities {
			container.DockerCapabilities = []string{"NONE"}
			c.info("Configure batch with all capabilities disabled")
		} else if *dockerCapabilities != "" {
			container.DockerCapabilities = getDockerCapabilitiesArray(*dockerCapabilities)
			c.info("Configure batch with those capabilities : %v", strings.Join(container.DockerCapabilities, ","))
		}

		schedulingGroupsToAdd := getSchedulingGroupsArray(UUID, client, *schedulingGroups)
		if len(schedulingGroupsToAdd) != 0 {
			c.info("Configure scheduling groups")
			container.SchedulingGroups = schedulingGroupsToAdd
		}

		if isFlagPassed("max-client-disconnect", c.flagSet) {
			c.info("Configure max-client-disconnect")
			container.MaxClientDisconnect = *maxClientDisconnect
		}

		msg := fmt.Sprintf(
			"Successfully configured service '%s' for project '%s'",
			*serviceName, projectToShow)

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
  sqsc service set          \
      -project="my-project" \
      -service="my-service" \
      -instances=42         \
      -env=env.json
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
