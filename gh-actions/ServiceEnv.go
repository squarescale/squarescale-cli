package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ServiceEnv struct{}

func (ev *ServiceEnv) create() {
	if _, exists := os.LookupEnv(mapEnvVar); exists {
		fmt.Println("Inserting database environment variables...")

		jsonFileName := "serviceEnvVar.json"
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(os.Getenv(mapEnvVar))), os.ModePerm)

		if jsonErr != nil {
			log.Fatal("Cannot write json file with map environment variables.")
		}

		instancesNumber := "1"

		cmd := fmt.Sprintf(
			"/sqsc container set -project-name %s/%s -env %s -service %s -instances %s -command \"%s\"",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			jsonFileName,
			os.Getenv(webServiceName),
			instancesNumber,
			getCmdEnvValue(),
		)
		fmt.Println(cmd)
		output, cmdErr := exec.Command("/bin/sh", "-c", cmd).Output()
		fmt.Println(string(output))

		if cmdErr != nil {
			fmt.Println(cmd)
			log.Fatal("Fail to import database environment variables.")
		}
	}
}
