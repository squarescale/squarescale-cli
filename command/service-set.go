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
func (cmd *ServiceSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	serviceName := serviceFlag(cmd.flagSet)
	instances := containerInstancesFlag(cmd.flagSet)
	runCommand := containerRunCmdFlag(cmd.flagSet)
	noRunCommand := containerNoRunCmdFlag(cmd.flagSet)
	entrypoint := entrypointFlag(cmd.flagSet)
	autostart := autostartFlag(cmd.flagSet)
	maxClientDisconnect := maxclientdisconnect(cmd.flagSet)

	schedulingGroups := containerSchedulingGroupsFlag(cmd.flagSet)
	limitMemory := containerLimitMemoryFlag(cmd.flagSet, 0)
	limitCPU := containerLimitCPUFlag(cmd.flagSet, 0)
	dockerCapabilities := dockerCapabilitiesFlag(cmd.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(cmd.flagSet)
	dockerDevices := dockerDevicesFlag(cmd.flagSet)

	envVariables := envFileFlag(cmd.flagSet)
	var envParams []string
	envParameterFlag(cmd.flagSet, &envParams)
	volumes := volumeFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *serviceName == "" {
		return cmd.errorWithUsage(errors.New("Service name cannot be empty."))
	}

	if *noRunCommand && *runCommand != "" {
		return cmd.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *instances <= 0 {
		if *instances != -1 {
			cmd.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set and default 1 will be used.")
		}
	}

	if *limitCPU < 0 || *limitMemory < 0 {
		return cmd.errorWithUsage(errors.New("Invalid values provided for limits."))
	}

	return cmd.runWithSpinner("configure service", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		container, err := client.GetServiceInfo(UUID, *serviceName)
		if err != nil {
			return "", err
		}

		cmd.Meta.spin.Stop()
		cmd.Meta.info("")

		if *instances > 0 && container.Size != *instances {
			container.Size = *instances
			cmd.info("Configure service with %d instances", *instances)
		}

		if *runCommand != "" && container.RunCommand != *runCommand {
			cmd.info("Configure service with run command: %s", *runCommand)
			container.RunCommand = *runCommand
			if len(container.RunCommand) == 0 {
				container.RunCommand = ""
			}
		}
		if *noRunCommand {
			cmd.info("Configure service without run command")
			container.RunCommand = ""
		}

		if *entrypoint != "" && container.Entrypoint != *entrypoint {
			cmd.info("Configure service with entrypoint: %s", *entrypoint)
			container.Entrypoint = *entrypoint
		}

		if *limitMemory > 0 && container.Limits.Memory != *limitMemory {
			cmd.info("Configure service with memory limit of %d MB", *limitMemory)
			container.Limits.Memory = *limitMemory
		}
		if *limitCPU > 0 && container.Limits.CPU != *limitCPU {
			cmd.info("Configure service with CPU limit of %d Mhz", *limitCPU)
			container.Limits.CPU = *limitCPU
		}

		if *envVariables != "" {
			cmd.info("Configure service with some environment JSON file")
			err := container.SetEnv(*envVariables)
			if err != nil {
				cmd.error(err)
			}
		}
		if len(envParams) > 0 {
			cmd.info("Configure service with some environment parameters")
			err := container.SetEnvParams(envParams)
			if err != nil {
				cmd.error(err)
			}
		}

		if isFlagPassed("auto-start", cmd.flagSet) {
			cmd.info("Configure service with CPU autostartFlag %v", *autostart)
			container.AutoStart = *autostart
		}

		if isFlagPassed("docker-devices", cmd.flagSet) {
			cmd.info("Configure service with custom Docker devices mapping")
			dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
			if err != nil {
				cmd.error(err)
			} else {
				container.DockerDevices = dockerDevicesArray
			}
		}

		if *noDockerCapabilities {
			container.DockerCapabilities = []string{"NONE"}
			cmd.info("Configure service with all capabilities disabled")
		} else if *dockerCapabilities != "" {
			container.DockerCapabilities = getDockerCapabilitiesArray(*dockerCapabilities)
			cmd.info("Configure service with those capabilities : %v", strings.Join(container.DockerCapabilities, ","))
		}

		schedulingGroupsToAdd := getSchedulingGroupsArray(UUID, client, *schedulingGroups)
		if len(schedulingGroupsToAdd) != 0 {
			cmd.info("Configure service with some scheduling groups")
			container.SchedulingGroups = schedulingGroupsToAdd
		}

		if isFlagPassed("max-client-disconnect", cmd.flagSet) {
			cmd.info("Configure max-client-disconnect")
			container.MaxClientDisconnect = *maxClientDisconnect
		}

		volumesToBind := parseVolumesToBind(*volumes)
		if len(volumesToBind) != 0 {
			cmd.info("Configure service with some volumes")
			container.Volumes = volumesToBind
		}

		msg := fmt.Sprintf(
			"Successfully configured service '%s' for project '%s'",
			*serviceName, projectToShow)

		cmd.Meta.spin.Start()
		return msg, client.ConfigService(container)
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ServiceSetCommand) Synopsis() string {
	return "Set service aka Docker container runtime parameters for project"
}

// Help is part of cli.Command implementation.
func (cmd *ServiceSetCommand) Help() string {
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
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
