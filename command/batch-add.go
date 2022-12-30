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
	DbDisable bool
	Db        squarescale.DbConfig
	Cluster   squarescale.ClusterConfig
}

// Run is part of cli.Command implementation.
func (b *BatchAddCommand) Run(args []string) int {
	// Parse flags
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := b.flagSet.String("project-uuid", "", "set the uuid of the project")
	projectName := b.flagSet.String("project-name", "", "set the name of the project")

	wantedBatchName := b.flagSet.String("name", "", "Batch name")
	runCommand := b.flagSet.String("run-command", "", "command / arguments that are used for execution")
	entrypoint := b.flagSet.String("entrypoint", "", "This is the script / program that will be executed")
	dockerImageName := b.flagSet.String("imageName", "", "docker image name")
	dockerImagePrivate := b.flagSet.Bool("imagePrivate", false, "docker image is private")
	dockerImageUsername := b.flagSet.String("imageUser", "", "docker image user")
	dockerImagePassword := b.flagSet.String("imagePwd", "", "docker image user")
	periodicBatch := b.flagSet.Bool("periodic", false, "batch periodicity")
	cronExpression := b.flagSet.String("cron", "", "cron expression")
	timeZoneName := b.flagSet.String("time", "", "time zone name")
	limitNet := b.flagSet.Int("net", 1, "This is an indicative limit of how much network bandwidth your service requires.")
	limitMemory := b.flagSet.Int("mem", 256, "This is the maximum amount of memory your service will be able to use until it is killed and restarted automatically.")
	limitIOPS := b.flagSet.Int("iops", 0, "This is an indicative limit of how many I/O operation per second your service requires.")
	limitCPU := b.flagSet.Int("cpu", 100, "This is an indicative limit of how much CPU your service requires.")
	dockerCapabilities, _ := dockerCapabilitiesFlag(b.flagSet)
	volumes := b.flagSet.String("volumes", "", "Volumes")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *wantedBatchName == "" && b.flagSet.Arg(0) != "" {
		name := b.flagSet.Arg(0)
		wantedBatchName = &name
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()[1:]))
	}

	//check rules
	if *wantedBatchName == "" {
		return b.errorWithUsage(fmt.Errorf(("Batch name is mandatory. Please, chose a batch name.")))
	}

	if *dockerImageName == "" {
		return b.errorWithUsage(fmt.Errorf(("DockerImage name is mandatory. Please, chose a dockerImage name.")))
	}

	if *dockerImagePrivate == true && *dockerImageUsername == "" {
		return b.errorWithUsage(fmt.Errorf(("Username is mandatory when the dockerImage is private. Please, complete the username.")))
	}

	if *dockerImagePrivate == true && *dockerImagePassword == "" {
		return b.errorWithUsage(fmt.Errorf(("Password is mandatory when the dockerImage is private. Please, complete the password.")))
	}

	if *periodicBatch == true && *cronExpression == "" {
		return b.errorWithUsage(fmt.Errorf(("Cron_expression is mandatory when the batch is periodic. Please, complete the cron_expression.")))
	}

	if *periodicBatch == true && *timeZoneName == "" {
		return b.errorWithUsage(fmt.Errorf(("Time_zone_name is mandatory when the batch is periodic. Please, complete the time_zone_name.")))
	}

	if _, err := time.LoadLocation(*timeZoneName); err != nil {
		return b.errorWithUsage(fmt.Errorf(("Time_zone_name is not in the good format (IANA Time Zone name)")))
	}

	if *limitNet < 0 {
		return b.errorWithUsage(fmt.Errorf(("NET must be strictly greater than 0")))
	}

	if *limitMemory <= 10 {
		return b.errorWithUsage(fmt.Errorf(("Memory must be greater than or equal to 10MB")))
	}

	if *limitCPU < 99 {
		return b.errorWithUsage(fmt.Errorf(("CPU must be greater than or equal to 100MHz")))
	}

	dockerCapabilitiesArray, err := getDockerCapabilitiesArray(*dockerCapabilities)
	if err != nil {
		return b.errorWithUsage(err)
	}

	volumesToBind := parseVolumesToBind(*volumes)

	//payload
	dockerImageContent := squarescale.DockerImageInfos{
		Name:     *dockerImageName,
		Private:  *dockerImagePrivate,
		Username: *dockerImageUsername,
		Password: *dockerImagePassword,
	}

	batchLimitContent := squarescale.BatchLimits{
		NET:    *limitNet,
		CPU:    *limitCPU,
		Memory: *limitMemory,
		IOPS:   *limitIOPS,
	}

	batchCommonContent := squarescale.BatchCommon{
		Name:           *wantedBatchName,
		Periodic:       *periodicBatch,
		CronExpression: *cronExpression,
		TimeZoneName:   *timeZoneName,
		RunCommand:     *runCommand,
		Entrypoint:     *entrypoint,
		Limits:         batchLimitContent,
		DockerCapabilities: dockerCapabilitiesArray,
	}

	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	res := b.runWithSpinner("add batch", endpoint.String(), func(client *squarescale.Client) (string, error) {

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

		//create function
		_, err = client.CreateBatch(UUID, batchOrderContent)

		return fmt.Sprintf("Successfully added batch '%s'", *wantedBatchName), err
	})

	if res != 0 {
		return res
	}

	return res

}

// Synopsis is part of cli.Command implementation.
func (b *BatchAddCommand) Synopsis() string {
	return "Add batch to project"
}

// Help is part of cli.Command implementation.
func (b *BatchAddCommand) Help() string {
	helpText := `
usage: sqsc batch add [options] <batch_name>

  Add batch to project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
