package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ServiceAddCommand is a cli.Command implementation for adding a service aka Docker container to project.
type ServiceAddCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ServiceAddCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)

	image := dockerImageNameFlag(cmd.flagSet)
	username := dockerImageUsernameFlag(cmd.flagSet)
	password := dockerImagePasswordFlag(cmd.flagSet)

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
	limitMemory := containerLimitMemoryFlag(cmd.flagSet, 256)
	limitCPU := containerLimitCPUFlag(cmd.flagSet, 100)
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

	if *noRunCommand && *runCommand != "" {
		return cmd.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *instances <= 0 {
		if *instances != -1 {
			cmd.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set and default 1 will be used.")
		}
	}

	if *limitCPU <= 0 || *limitMemory <= 0 {
		return cmd.errorWithUsage(errors.New("Invalid values provided for limits."))
	}

	if err := cmd.validateArgs(*image, *instances); err != nil {
		return cmd.errorWithUsage(err)
	}

	dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
	if isFlagPassed("docker-devices", cmd.flagSet) {
		if err != nil {
			return cmd.errorWithUsage(err)
		}
	}

	volumesToBind := parseVolumesToBind(*volumes)

	payload := squarescale.JSONObject{
		"docker_image": squarescale.DockerImage(
			*image,
			*username,
			*password,
		),
		"auto_start":		*autostart,
		"size":			*instances,
	}

	limits := squarescale.JSONObject{}
	if *limitMemory > 0 {
		limits["mem"] = *limitMemory
	}
	if *limitCPU > 0 {
		limits["cpu"] = *limitCPU
	}
	if len(limits) > 0 {
		payload["container"] = limits
	}

	if len(*maxClientDisconnect) > 0 {
		payload["max_client_disconnect"] = *maxClientDisconnect
	} else {
		payload["max_client_disconnect"] = "0"
	}

	if len(*serviceName) > 0 {
		payload["name"] = *serviceName
	}

	if len(*entrypoint) > 0 {
		payload["entrypoint"] = *entrypoint
	}

	if len(*runCommand) > 0 {
		payload["run_command"] = *runCommand
	}

	if len(*dockerDevices) > 0 {
		payload["docker_devices"] = dockerDevicesArray
	}

	if len(volumesToBind) > 0 {
		payload["volumes_to_bind"] = volumesToBind
	}

	container := squarescale.Service{}
	if *envVariables != "" {
		err := container.SetEnv(*envVariables)
		if err != nil {
			cmd.error(err)
		}
	}
	if len(envParams) > 0 {
		err := container.SetEnvParams(envParams)
		if err != nil {
			cmd.error(err)
		}
	}
	if len(container.CustomEnv) > 0 {
		payload["custom_environment"] = container.CustomEnv
	}

	dockerCapabilitiesArray := []string{}
	if *noDockerCapabilities {
		dockerCapabilitiesArray = []string{"NONE"}
	} else if *dockerCapabilities != "" {
		dockerCapabilitiesArray = getDockerCapabilitiesArray(*dockerCapabilities)
	} else {
		dockerCapabilitiesArray = getDefaultDockerCapabilitiesArray()
	}
	payload["docker_capabilities"] = dockerCapabilitiesArray

	return cmd.runWithSpinner("adding Docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		// IMPORTANT NOTE: if scheduling group does not exist
		// you only get non fatal error output like
		// adding Docker imageScheduling group 'dummy' not found for project '8944c0fd-98a5-4620-96a0-0d5603cb0e6d'
		// but service will still be created
		schedulingGroupsToAdd := getSchedulingGroupsIntArray(UUID, client, *schedulingGroups)
		if len(schedulingGroupsToAdd) != 0 {
			payload["scheduling_groups"] = schedulingGroupsToAdd
		}
		msg := fmt.Sprintf("Successfully added Docker image '%s' to project '%s' (%v instance(s))", *image, projectToShow, *instances)
		return msg, client.AddService(UUID, payload)
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ServiceAddCommand) Synopsis() string {
	return "Add service aka Docker container to project"
}

// Help is part of cli.Command implementation.
func (cmd *ServiceAddCommand) Help() string {
	helpText := `
usage: sqsc service add [options]

  Add service aka Docker container to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

func (cmd *ServiceAddCommand) validateArgs(image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}
	return nil
}
