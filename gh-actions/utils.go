package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func executeCommand(cmd string, errorMsg string) {
	fmt.Println(cmd)
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		log.Fatal(errorMsg)
	}

	fmt.Println(string(output))
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

func getCmdEnvValue() string {
	cmdEnvValue := ""
	if _, cmdExists := os.LookupEnv(cmdEnv); cmdExists {
		cmdEnvValue = os.Getenv(cmdEnv)
	}
	return cmdEnvValue
}
