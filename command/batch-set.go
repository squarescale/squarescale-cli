package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchSetCommand allows to configure a project container.
type BatchSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *BatchSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	batchName := batchNameFlag(c.flagSet)
	runCmdArg := batchRunCmdFlag(c.flagSet)
	entrypoint := entrypointFlag(c.flagSet)
	limitMemoryArg := batchLimitMemoryFlag(c.flagSet)
	limitCPUArg := batchLimitCPUFlag(c.flagSet)
	noRunCmdArg := batchNoRunCmdFlag(c.flagSet)
	dockerCapabilities := dockerCapabilitiesFlag(c.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(c.flagSet)
	dockerDevices := dockerDevicesFlag(c.flagSet)
	envCmdArg := envFileFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *batchName == "" {
		return c.errorWithUsage(errors.New("Batch name cannot be empty."))
	}

	if *noRunCmdArg && *runCmdArg != "" {
		return c.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	return c.runWithSpinner("configure batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}

		batch, err := client.GetBatchesInfo(UUID, *batchName)
		if err != nil {
			return "", err
		}

		c.Meta.spin.Stop()
		c.Meta.info("")
		if *runCmdArg != "" {
			c.info("Configure batch with run command: %s", *runCmdArg)
			batch.RunCommand = *runCmdArg
		}
		if *noRunCmdArg {
			c.info("Configure batch without run command")
			batch.RunCommand = ""
		}

		if *entrypoint != "" {
			c.info("Configure batch with entrypoint: %s", *entrypoint)
			batch.Entrypoint = *entrypoint
		}

		if *limitMemoryArg >= 0 {
			c.info("Configure batch with memory limit of %d MB", *limitMemoryArg)
			batch.Limits.Memory = *limitMemoryArg
		}
		if *limitCPUArg >= 0 {
			c.info("Configure batch with CPU limit of %d Mhz", *limitCPUArg)
			batch.Limits.CPU = *limitCPUArg
		}

		if *noDockerCapabilities {
			batch.DockerCapabilities = []string{"NONE"}
			c.info("Configure batch with all capabilities disabled")
		} else if *dockerCapabilities != "" {
			batch.DockerCapabilities = getDockerCapabilitiesArray(*dockerCapabilities)
			c.info("Configure batch with those capabilities : %v", strings.Join(batch.DockerCapabilities, ","))
		}

		if isFlagPassed("docker-devices", c.flagSet) {
			dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
			if err != nil {
				return "", err
			}
			batch.DockerDevices = dockerDevicesArray
			c.info("Configure batch with custom mapping")
		}

		if *envCmdArg != "" {
			c.info("Configure batch with some env")
			err := batch.SetEnv(*envCmdArg)
			if err != nil {
				c.error(err)
			}
		}

		msg := fmt.Sprintf(
			"Successfully configured batch '%s' for project '%s'",
			*batchName, *projectUUID)

		c.Meta.spin.Start()
		return msg, client.ConfigBatch(batch, UUID)
	})
}

// Synopsis is part of cli.Command implementation.
func (c *BatchSetCommand) Synopsis() string {
	return "Set batch runtime parameters for project"
}

// Help is part of cli.Command implementation.
func (c *BatchSetCommand) Help() string {
	helpText := `
usage: sqsc batch set [options]

  Set batch runtime parameters for project.
  Batch are specified using their given name.

Example:
  sqsc batch set                       \
      -project-name my-rails-project   \
			-batch-name my-batch             \
      -e env.json
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
