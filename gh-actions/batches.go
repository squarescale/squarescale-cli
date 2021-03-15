package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Batches struct{}

type BatchContent struct {
	Cmd string
	Env map[string]interface{}
}

func (b *Batches) create() {
	if _, exists := os.LookupEnv(batchesEnv); exists {
		var batches map[string]BatchContent
		json.Unmarshal([]byte(os.Getenv(batchesEnv)), &batches)

		for batchName, batchContent := range batches {
			b.createBatch(batchName, batchContent)
			b.insertBatchEnv(batchName, batchContent.Env)
			b.executeBatch(batchName)
		}
	}
}

func (b *Batches) createBatch(batchName string, batchContent BatchContent) {
	cmd := fmt.Sprintf(
		"/sqsc batch add -project-name %s/%s -imageName %s:%s -imagePrivate -imageUser %s -imagePwd %s -name %s -run-command %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(dockerRepository),
		os.Getenv(dockerRepositoryTag),
		os.Getenv(dockerUser),
		os.Getenv(dockerToken),
		batchName,
		batchContent.Cmd,
	)
	executeCommand(cmd, fmt.Sprintf("Fail to add batch %q.", batchName))
}

func (b *Batches) insertBatchEnv(batchName string, batchContentEnv map[string]interface{}) {
	if len(batchContentEnv) != 0 {
		jsonFileName := "batchEnvVar.json"
		d, _ := json.Marshal(batchContentEnv)
		data := mapDatabaseEnv(string(d))
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(data), os.ModePerm)

		if jsonErr != nil {
			log.Fatal(fmt.Sprintf("Cannot write json file with env for batch %q.", batchName))
		}

		//TODO: remove next line juste below
		executeCommand(fmt.Sprintf("cat %s", jsonFileName), fmt.Sprintf("Fail to cat %s", jsonFileName))

		cmd := fmt.Sprintf(
			"./sqsc batch set -project-name %s/%s -batch-name %s -env %s",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			batchName,
			jsonFileName,
		)
		executeCommand(cmd, "Fail to insert batch env.")
	}
}

func (b *Batches) executeBatch(batchName string) {
	fmt.Println(fmt.Sprintf("Executing batch %q ...", batchName))

	cmd := fmt.Sprintf(
		"/sqsc batch exec -project-name %s/%s %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		batchName,
	)
	executeCommand(cmd, fmt.Sprintf("Fail to execute batch %q", batchName))
}
