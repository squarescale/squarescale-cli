package main

import (
	"fmt"
	"os"
	"os/exec"
)

type Project struct{}

func (p *Project) create() {
	if !isProjectExists() {
		p.createProject()
	} else {
		fmt.Println("Project already exists.")
	}
}

func (p *Project) createProject() {
	fmt.Println("Creating project...")

	cmd := fmt.Sprintf("%s project create", sqscCLICmd)
	cmd += " -credential " + os.Getenv(iaasCred)
	cmd += " -name " + os.Getenv(projectName)
	cmd += " -node-size " + os.Getenv(nodeType)
	cmd += " -infra-type " + os.Getenv(infraType)
	value, exists := os.LookupEnv(organizationName)
	if exists && len(value) > 0 {
		cmd += " -organization " + value
	}
	cmd += " -provider " + os.Getenv(iaasProvider)
	cmd += " -region " + os.Getenv(iaasRegion)

	if os.Getenv(monitoring) != "" {
		cmd += " -monitoring " + os.Getenv(monitoring)
	}

	cmd += " -yes "

	executeCommand(cmd, fmt.Sprintf("Failed to create project %q.", os.Getenv(projectName)))
}

func isProjectExists() bool {
	_, projectNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc project get -project-name %s",
		getProjectName(),
	)).Output()

	return projectNotExists == nil
}
