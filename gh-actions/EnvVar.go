package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hoisie/mustache"
)

type EnvVar struct{}

func (ev *EnvVar) create() {
	if _, exists := os.LookupEnv(mapEnvVar); exists {
		fmt.Println("Inserting database environment variables...")

		data := mustache.Render(os.Getenv(mapEnvVar), map[string]string{
			"DB_HOST":     getSQSCEnvValue("DB_HOST"),
			"DB_USERNAME": getSQSCEnvValue("DB_USERNAME"),
			"DB_PASSWORD": getSQSCEnvValue("DB_PASSWORD"),
		})

		jsonFileName := "mapEnvVar.json"
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(data), os.ModePerm)

		if jsonErr != nil {
			log.Fatal("Cannot write json file with map environment variables.")
		}

		fmt.Println("The json generated looks like:")
		jsonOutput, _ := exec.Command("/bin/sh", "-c", fmt.Sprintf("cat %s", jsonFileName)).Output()
		fmt.Println(string(jsonOutput))

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

func getSQSCEnvValue(key string) string {
	value, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc env get -project-name %s/%s \"%s\" | grep -v %s | tr -d '\n'",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		key,
		"...done",
	)).Output()

	if err != nil {
		fmt.Println(fmt.Sprintf("Environment variable %q does not exists in this project.", key))
		return ""
	}

	return string(value)
}
