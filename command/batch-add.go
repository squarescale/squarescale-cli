package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectAddCommand is a cli.Command implementation for creating a Squarescale project.
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
	wantedBatchName := batchNameFlag(b.flagSet)
	endpoint := endpointFlag(b.flagSet)
	project := projectFlag(b.flagSet)

	DockerImageName := batchDockerImageNameFlag(b.flagSet)
	DockerImagePrivate := batchDockerImagePrivateFlag(b.flagSet)   // FIXME
	DockerImageUsername := batchDockerImageUsernameFlag(b.flagSet) // FIXME
	DockerImagePassword := batchDockerImagePasswordFlag(b.flagSet) // FIXME
	PeriodicBatch := batchPeriodicFlag(b.flagSet)
	CronExpression := batchCronExpressionFlag(b.flagSet)
	TimeZoneName := batchTimeZoneNameFlag(b.flagSet)
	LimitNet := batchLimitNetFlag(b.flagSet)
	LimitMemory := batchLimitMemoryFlag(b.flagSet)
	LimitIOPS := batchLimitIopsFlag(b.flagSet)
	LimitCPU := batchLimitCpuFlag(b.flagSet)
	volumes := b.flagSet.String("volumes", "", "Volumes")

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *wantedBatchName == "" && b.flagSet.Arg(0) != "" {
		name := b.flagSet.Arg(0)
		wantedBatchName = &name
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()[1:]))
	}

	if err := validateProjectName(*project); err != nil {
		return b.errorWithUsage(err)
	}

	//check rules

	if *wantedBatchName == "" {
		return b.errorWithUsage(fmt.Errorf(("Batch name is mandatory. Please, chose a batch name.")))
	}

	if *DockerImageName == "" {
		return b.errorWithUsage(fmt.Errorf(("DockerImage name is mandatory. Please, chose a dockerImage name.")))
	}

	if *DockerImagePrivate == true && *DockerImageUsername == "" {
		return b.errorWithUsage(fmt.Errorf(("Username is mandatory when the dockerImage is private. Please, complete the username.")))
	}

	if *DockerImagePrivate == true && *DockerImagePassword == "" {
		return b.errorWithUsage(fmt.Errorf(("Password is mandatory when the dockerImage is private. Please, complete the password.")))
	}

	if *PeriodicBatch == true && *CronExpression == "" {
		return b.errorWithUsage(fmt.Errorf(("Cron_expression is mandatory when the batch is periodic. Please, complete the cron_expression.")))
	}

	if *PeriodicBatch == true && *TimeZoneName == "" {
		return b.errorWithUsage(fmt.Errorf(("Time_zone_name is mandatory when the batch is periodic. Please, complete the time_zone_name.")))
	}

	// _, err = time.LoadLocation(*TimeZoneName)
	// if err != nil {
	// 	return b.errorWithUsage(fmt.Errorf(("Time_zone_name is not in the good format (IANA Time Zone name)")))
	// }

	if *LimitNet < 0 {
		return b.errorWithUsage(fmt.Errorf(("NET must be strictly greater than 0")))
	}

	if *LimitMemory <= 10 {
		return b.errorWithUsage(fmt.Errorf(("Memory must be greater than or equal to 10MB")))
	}

	if *LimitCPU < 99 {
		return b.errorWithUsage(fmt.Errorf(("CPU must be greater than or equal to 100MHz")))
	}

	volumesToBind := parseVolumesToBind(*volumes)

	//payload
	dockerImageContent := squarescale.DockerImageInfos{
		Name:     *DockerImageName,
		Private:  *DockerImagePrivate,
		Username: *DockerImageUsername,
		Password: *DockerImagePassword,
	}

	batchLimitContent := squarescale.BatchLimits{
		NET:    *LimitNet,
		CPU:    *LimitCPU,
		Memory: *LimitMemory,
		IOPS:   *LimitIOPS,
	}

	batchCommonContent := squarescale.BatchCommon{
		Name:           *wantedBatchName,
		Periodic:       *PeriodicBatch,
		CronExpression: *CronExpression,
		TimeZoneName:   *TimeZoneName,
		Limits:         batchLimitContent,
	}

	batchOrderContent := squarescale.BatchOrder{
		BatchCommon: batchCommonContent,
		DockerImage: dockerImageContent,
		Volumes:     volumesToBind,
	}

	res := b.runWithSpinner("add batch", endpoint.String(), func(client *squarescale.Client) (string, error) {

		var err error

		//create function
		_, err = client.CreateBatch(*project, batchOrderContent)

		return fmt.Sprintf("Successfully added batch '%s'", *wantedBatchName), err
	})

	if res != 0 {
		return res
	}

	return res

}

// Synopsis is part of cli.Command implementation.
func (b *BatchAddCommand) Synopsis() string {
	return "Add a new Batch in Squarescale project"
}

// Help is part of cli.Command implementation.
func (b *BatchAddCommand) Help() string {
	helpText := `
usage: sqsc batch add [options] <batch_name>

  Add a new batch using the provided batch name (name is mandatory).

`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
