package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Batches struct{}

type BatchContent struct {
	Cmd string
	Env map[string]interface{}
}

func (b *Batches) create() {
	os.Setenv(batchesEnv, `{"database-setup": {"cmd": "bundle exec rails db:setup","env": {"a": "a", "b": "b"}}, "database-seed": {"cmd": "bundle exec rails db:sogilis:seed"}}`)

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
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to add batch %q", batchName))
	}

	fmt.Println(output)
}

func (b *Batches) insertBatchEnv(batchName string, env map[string]interface{}) {
	if len(env) != 0 {
		jsonFileName := "batchEnvVar.json"
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(os.Getenv(batchesEnv))), os.ModePerm)

		if jsonErr != nil {
			log.Fatal(fmt.Sprintf("Cannot write json file with env for batch %q", batchName))
		}

		cmd := fmt.Sprintf(
			"./sqsc batch set -project-name %s/%s -batch-name %s -env %s",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			batchName,
			jsonFileName,
		)
		fmt.Println(cmd)
		output, err := exec.Command("/bin/sh", "-c", cmd).Output()

		if err != nil {
			log.Fatal("Fail to insert batch env.")
		}

		fmt.Println(output)
	}
}
