package main

import (
	"fmt"
	"os"
	"os/exec"
)

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
