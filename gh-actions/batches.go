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

func (b *Batches) insertBatchEnv(batchName string, env map[string]interface{}) {
	if len(env) != 0 {
		jsonFileName := "batchEnvVar.json"
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(os.Getenv(batchesEnv))), os.ModePerm)

		if jsonErr != nil {
			log.Fatal(fmt.Sprintf("Cannot write json file with env for batch %q.", batchName))
		}

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
