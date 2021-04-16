package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"gopkg.in/alessio/shellescape.v1"
)

type Batches struct{}

type BatchContent struct {
	EXECUTE        bool              `json:"execute"`
	IMAGE_NAME     string            `json:"image_name"`
	IS_PRIVATE     bool              `json:"is_private"`
	IMAGE_USER     string            `json:"image_user"`
	IMAGE_PASSWORD string            `json:"image_password"`
	RUN_CMD        string            `json:"run_cmd"`
	PERIODIC       BatchPeriodic     `json:"periodic"`
	ENV            map[string]string `json:"env"`
	LIMIT_MEMORY   string            `json:"limit_memory"`
	LIMIT_CPU      string            `json:"limit_cpu"`
	LIMIT_NET      string            `json:"limit_net"`
}

type BatchPeriodic struct {
	PERIODICITY string `json:"periodicity"`
	TIMEZONE    string `json:"timezone"`
}

func (b *Batches) create() {
	if _, exists := os.LookupEnv(batchesEnv); exists {
		var batches map[string]BatchContent

		err := json.Unmarshal([]byte(os.Getenv(batchesEnv)), &batches)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error when Unmarshal: %s", err))
		}

		for batchName, batchContent := range batches {
			if !isBatchExists(batchName) {
				b.createBatch(batchName, batchContent)
				b.insertBatchEnvAndLimits(batchName, batchContent)
			} else {
				fmt.Println(fmt.Sprintf("Batch %q already exists.", batchName))
			}

			if batchContent.EXECUTE {
				b.executeBatch(batchName)
			}
		}
	}
}

func (b *Batches) createBatch(batchName string, batchContent BatchContent) {
	fmt.Println(fmt.Sprintf("Creating %q batch...", batchName))

	cmd := "/sqsc batch add"
	cmd += " -name " + batchName
	cmd += " -project-name " + getProjectName()
	cmd += " -run-command " + shellescape.Quote(batchContent.RUN_CMD)

	imageName := batchContent.IMAGE_NAME
	if imageName == "" {
		imageName = getDockerImage()
	}
	cmd += " -imageName " + imageName

	if batchContent.IS_PRIVATE {
		cmd += " -imagePrivate"
		cmd += " -imageUser " + batchContent.IMAGE_USER
		cmd += " -imagePwd " + batchContent.IMAGE_PASSWORD
	}

	if (BatchPeriodic{}) != batchContent.PERIODIC {
		periodicity := batchContent.PERIODIC.PERIODICITY
		if periodicity == "" {
			periodicity = "* * * * *"
		}

		timezone := batchContent.PERIODIC.TIMEZONE
		if timezone == "" {
			timezone = "Europe/Paris"
		}

		cmd += " -periodic"
		cmd += " -cron " + shellescape.Quote(periodicity)
		cmd += " -time \"" + timezone + "\""
	}

	executeCommand(cmd, fmt.Sprintf("Fail to add %q batch.", batchName))
}

func (b *Batches) insertBatchEnvAndLimits(batchName string, batchContent BatchContent) {

	limitMemory := batchContent.LIMIT_MEMORY
	limitNet := batchContent.LIMIT_NET
	limitCpu := batchContent.LIMIT_CPU
	environment := batchContent.ENV

	if len(environment) != 0 || limitMemory != "" || limitNet != "" || limitCpu != "" {

		cmd := "/sqsc batch set"
		cmd += " -project-name " + getProjectName()
		cmd += " -batch-name " + batchName

		if len(environment) != 0 {
			fmt.Println(fmt.Sprintf("Inserting environment variable to batch %q", batchName))

			jsonFileName := "batchEnvVar.json"
			env, _ := json.Marshal(environment)
			jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(string(env))), os.ModePerm)

			if jsonErr != nil {
				log.Fatal(fmt.Sprintf("Cannot write json file with env for batch %q.", batchName))
			}

			cmd += " -env " + jsonFileName
		}

		if limitMemory != "" {
			cmd += " -memory " + limitMemory
		}

		if limitCpu != "" {
			cmd += " -cpu " + limitCpu
		}

		if limitNet != "" {
			cmd += " -net " + limitNet
		}

		executeCommand(cmd, "Fail to insert batch env.")
	}
}

func (b *Batches) executeBatch(batchName string) {
	fmt.Println(fmt.Sprintf("Executing %q batch ...", batchName))

	cmd := fmt.Sprintf(
		"/sqsc batch exec -project-name %s %s",
		getProjectName(),
		batchName,
	)
	executeCommand(cmd, fmt.Sprintf("Fail to execute %q batch.", batchName))
}

func isBatchExists(batchName string) bool {
	_, batchNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc batch list -project-name %s | grep %s",
		getProjectName(),
		batchName,
	)).Output()

	return batchNotExists == nil
}
