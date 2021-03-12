package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Webservice struct{}

func (ws *Webservice) create() {
	if !isWebserviceExists() {
		ws.createWebservice()
	} else {
		fmt.Println("Service already exists.")
	}
}

func (ws *Webservice) createWebservice() {
	fmt.Println("Creating web service...")

	cmdEnvValue := getCmdEnvValue()

	cmd := fmt.Sprintf(
		"/sqsc container add -project-name %s/%s -servicename %s -name %s:%s -username %s -password %s -run-command \"%s\"",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
		os.Getenv(dockerRepository),
		os.Getenv(dockerRepositoryTag),
		os.Getenv(dockerUser),
		os.Getenv(dockerToken),
		cmdEnvValue,
	)
	fmt.Println(cmd)
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	fmt.Println(string(output))

	if err != nil {
		fmt.Println(cmd)
		log.Fatal(fmt.Sprintf("Creating web service fails with error:\n %s", err))
	}
}

func isWebserviceExists() bool {
	_, webServiceNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc container list --project-name %s/%s | grep %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
	)).Output()

	return webServiceNotExists == nil
}
