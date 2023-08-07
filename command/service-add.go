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
func (c *ServiceAddCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)

	image := dockerImageNameFlag(c.flagSet)
	username := dockerImageUsernameFlag(c.flagSet)
	password := dockerImagePasswordFlag(c.flagSet)

	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	serviceName := serviceFlag(c.flagSet)
	instances := containerInstancesFlag(c.flagSet)
	runCommand := containerRunCmdFlag(c.flagSet)
	noRunCommand := containerNoRunCmdFlag(c.flagSet)
	entrypoint := entrypointFlag(c.flagSet)
	autostart := autostartFlag(c.flagSet)
	maxClientDisconnect := maxclientdisconnect(c.flagSet)

	schedulingGroups := containerSchedulingGroupsFlag(c.flagSet)
	limitMemory := containerLimitMemoryFlag(c.flagSet, 256)
	limitCPU := containerLimitCPUFlag(c.flagSet, 100)
	dockerCapabilities := dockerCapabilitiesFlag(c.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(c.flagSet)
	dockerDevices := dockerDevicesFlag(c.flagSet)

	envVariables := envFileFlag(c.flagSet)
	var envParams []string
	envParameterFlag(c.flagSet, &envParams)
	volumes := volumeFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *noRunCommand && *runCommand != "" {
		return c.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *instances <= 0 {
		if *instances != -1 {
			c.Ui.Warn("Number of instances cannot be 0 or negative. This value won't be set and default 1 will be used.")
		}
	}

	if *limitCPU <= 0 || *limitMemory <= 0 {
		return c.errorWithUsage(errors.New("Invalid values provided for limits."))
	}

	if err := c.validateArgs(*image, *instances); err != nil {
		return c.errorWithUsage(err)
	}

	dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
	if isFlagPassed("docker-devices", c.flagSet) {
		if err != nil {
			return c.errorWithUsage(err)
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
			c.error(err)
		}
	}
	if len(envParams) > 0 {
		err := container.SetEnvParams(envParams)
		if err != nil {
			c.error(err)
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

	return c.runWithSpinner("adding Docker image", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		schedulingGroupsToAdd := getSchedulingGroupsArray(UUID, client, *schedulingGroups)
		if len(schedulingGroupsToAdd) != 0 {
			payload["scheduling_groups"] = schedulingGroupsToAdd
		}
		msg := fmt.Sprintf("Successfully added Docker image '%s' to project '%s' (%v instance(s))", *image, projectToShow, *instances)
		return msg, client.AddService(UUID, payload)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ServiceAddCommand) Synopsis() string {
	return "Add service aka Docker container to project"
}

// Help is part of cli.Command implementation.
func (c *ServiceAddCommand) Help() string {
	helpText := `
usage: sqsc service add [options]

  Add service aka Docker container to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ServiceAddCommand) validateArgs(image string, instances int) error {
	if image == "" {
		return errors.New("Docker image name cannot be empty")
	}

	if instances <= 0 {
		return errors.New("Invalid value provided for instances count: it cannot be 0 or negative")
	}
	return nil
}
