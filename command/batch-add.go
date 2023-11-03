package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// BatchAddCommand is a struct to define a batch
type BatchAddCommand struct {
	Meta
	flagSet   *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *BatchAddCommand) Run(args []string) int {
	// Parse flags
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)

	image := dockerImageNameFlag(cmd.flagSet)
	username := dockerImageUsernameFlag(cmd.flagSet)
	password := dockerImagePasswordFlag(cmd.flagSet)

	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	batchName := serviceFlag(cmd.flagSet)
	runCommand := containerRunCmdFlag(cmd.flagSet)
	noRunCommand := containerNoRunCmdFlag(cmd.flagSet)
	entrypoint := entrypointFlag(cmd.flagSet)

	periodicBatch := periodicBatchFlag(cmd.flagSet)
	cronExpression := cronExpressionFlag(cmd.flagSet)
	timeZoneName := timeZoneNameFlag(cmd.flagSet)

	schedulingGroups := containerSchedulingGroupsFlag(cmd.flagSet)
	limitMemory := containerLimitMemoryFlag(cmd.flagSet, 256)
	limitCPU := containerLimitCPUFlag(cmd.flagSet, 100)
	dockerCapabilities := dockerCapabilitiesFlag(cmd.flagSet)
	noDockerCapabilities := noDockerCapabilitiesFlag(cmd.flagSet)
	dockerDevices := dockerDevicesFlag(cmd.flagSet)

	envVariables := envFileFlag(cmd.flagSet)
	var envParams []string
	envParameterFlag(cmd.flagSet, &envParams)
	volumes := cmd.flagSet.String("volumes", "", "Volumes")

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *noRunCommand && *runCommand != "" {
		return cmd.errorWithUsage(errors.New("Cannot specify an override command and disable it at the same time"))
	}

	if *batchName == "" && cmd.flagSet.Arg(0) != "" {
		name := cmd.flagSet.Arg(0)
		batchName = &name
	}

	if cmd.flagSet.NArg() > 1 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()[1:]))
	}

	//check rules
	if *batchName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Batch name is mandatory. Please, chose a batch name.")))
	}

	if *image == "" {
		return cmd.errorWithUsage(fmt.Errorf(("DockerImage name is mandatory. Please, chose a dockerImage name.")))
	}

	if *periodicBatch == true && *cronExpression == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Cron_expression is mandatory when the batch is periodic. Please, complete the cron_expression.")))
	}

	if *periodicBatch == true && *timeZoneName == "" {
		return cmd.errorWithUsage(fmt.Errorf(("Time_zone_name is mandatory when the batch is periodic. Please, complete the time_zone_name.")))
	}

	if _, err := time.LoadLocation(*timeZoneName); err != nil {
		return cmd.errorWithUsage(fmt.Errorf(("Time_zone_name is not in the good format (IANA Time Zone name)")))
	}

	if *limitMemory <= 10 {
		return cmd.errorWithUsage(fmt.Errorf(("Memory must be greater than or equal to 10MB")))
	}

	if *limitCPU < 99 {
		return cmd.errorWithUsage(fmt.Errorf(("CPU must be greater than or equal to 100MHz")))
	}

	dockerCapabilitiesArray := []string{}
	if *noDockerCapabilities {
		dockerCapabilitiesArray = []string{"NONE"}
	} else if *dockerCapabilities != "" {
		dockerCapabilitiesArray = getDockerCapabilitiesArray(*dockerCapabilities)
	} else {
		dockerCapabilitiesArray = getDefaultDockerCapabilitiesArray()
	}

	dockerDevicesArray, err := getDockerDevicesArray(*dockerDevices)
	if isFlagPassed("docker-devices", cmd.flagSet) {
		if err != nil {
			return cmd.errorWithUsage(err)
		}
	}

	volumesToBind := parseVolumesToBind(*volumes)

	//payload
	dockerImageContent := squarescale.DockerImageInfos{
		Name:     *image,
		Private:  (*username != "" && *password != ""),
		Username: *username,
		Password: *password,
	}

	batchLimitContent := squarescale.BatchLimits{
		CPU:    *limitCPU,
		Memory: *limitMemory,
	}

	batchCommonContent := squarescale.BatchCommon{
		Name:               *batchName,
		Periodic:           *periodicBatch,
		CronExpression:     *cronExpression,
		TimeZoneName:       *timeZoneName,
		RunCommand:         *runCommand,
		Entrypoint:         *entrypoint,
		Limits:             batchLimitContent,
		DockerCapabilities: dockerCapabilitiesArray,
	}

	if len(dockerDevicesArray) > 0 {
		batchCommonContent.DockerDevices = dockerDevicesArray
	}

	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
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
		fmt.Println("XXX")
		//payload["custom_environment"] = container.CustomEnv
	}

	res := cmd.runWithSpinner("add batch", endpoint.String(), func(client *squarescale.Client) (string, error) {

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
		schedulingGroupsToAdd := getSchedulingGroupsIntArray(UUID, client, *schedulingGroups)
/*		if len(schedulingGroupsToAdd) != 0 {
			payload["scheduling_groups"] = schedulingGroupsToAdd
		}*/
		fmt.Printf("GOT %s\nBATCH\n%+v\nSCHED\n%+v\n", UUID, batchOrderContent, schedulingGroupsToAdd)
		return fmt.Sprintf("Not creating batch '%s'", *batchName), fmt.Errorf("toto")
/*		//create function
		_, err = client.CreateBatch(UUID, batchOrderContent)

		return fmt.Sprintf("Successfully added batch '%s'", *batchName), err*/
	})

	if res != 0 {
		return res
	}

	return res

}

// Synopsis is part of cli.Command implementation.
func (cmd *BatchAddCommand) Synopsis() string {
	return "Add batch to project"
}

// Help is part of cli.Command implementation.
func (cmd *BatchAddCommand) Help() string {
	helpText := `
usage: sqsc batch add [options] <batch_name>

  Add batch to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
