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
func (cmd *BatchSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	batchName := batchNameFlag(cmd.flagSet)
	runCmdArg := batchRunCmdFlag(cmd.flagSet)
	entrypoint := entrypointFlag(cmd.flagSet)
	limitMemoryArg := batchLimitMemoryFlag(cmd.flagSet)
	limitCPUArg := batchLimitCPUFlag(cmd.flagSet)
	noRunCmdArg := batchNoRunCmdFlag(cmd.flagSet)
	dockerCapabilities := dockerCapabilitiesFlag(cmd.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(cmd.flagSet)
	dockerDevices := dockerDevicesFlag(cmd.flagSet)
	envCmdArg := envFileFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *batchName == "" {
		return cmd.errorWithUsage(errors.New("Batch name cannot be empty."))
	}

	if *noRunCmdArg && *runCmdArg != "" {
		return cmd.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	return cmd.runWithSpinner("configure batch", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		cmd.Meta.spin.Stop()
		cmd.Meta.info("")
		if *runCmdArg != "" {
			cmd.info("Configure batch with run command: %s", *runCmdArg)
			batch.RunCommand = *runCmdArg
		}
		if *noRunCmdArg {
			cmd.info("Configure batch without run command")
			batch.RunCommand = ""
		}

		if *entrypoint != "" {
			cmd.info("Configure batch with entrypoint: %s", *entrypoint)
			batch.Entrypoint = *entrypoint
		}

		if *limitMemoryArg >= 0 {
			cmd.info("Configure batch with memory limit of %d MB", *limitMemoryArg)
			batch.Limits.Memory = *limitMemoryArg
		}
		if *limitCPUArg >= 0 {
			cmd.info("Configure batch with CPU limit of %d Mhz", *limitCPUArg)
			batch.Limits.CPU = *limitCPUArg
		}

		if *noDockerCapabilities {
			batch.DockerCapabilities = []string{"NONE"}
			cmd.info("Configure batch with all capabilities disabled")
		} else if *dockerCapabilities != "" {
			batch.DockerCapabilities = getDockerCapabilitiesArray(*dockerCapabilities)
			cmd.info("Configure batch with those capabilities : %v", strings.Join(batch.DockerCapabilities, ","))
		}

		if isFlagPassed("docker-devices", cmd.flagSet) {
			dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
			if err != nil {
				return "", err
			}
			batch.DockerDevices = dockerDevicesArray
			cmd.info("Configure batch with custom mapping")
		}

		if *envCmdArg != "" {
			cmd.info("Configure batch with some env")
			err := batch.SetEnv(*envCmdArg)
			if err != nil {
				cmd.error(err)
			}
		}

		msg := fmt.Sprintf(
			"Successfully configured batch '%s' for project '%s'",
			*batchName, *projectUUID)

		cmd.Meta.spin.Start()
		return msg, client.ConfigBatch(batch, UUID)
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *BatchSetCommand) Synopsis() string {
	return "Set batch runtime parameters for project"
}

// Help is part of cli.Command implementation.
func (cmd *BatchSetCommand) Help() string {
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
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
