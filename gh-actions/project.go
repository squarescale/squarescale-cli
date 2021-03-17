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

	cmd := fmt.Sprintf(
		"/sqsc project create -credential %s -monitoring netdata -name %s -node-size %s -infra-type high-availability -organization %s -provider %s -region %s -yes",
		os.Getenv(iaasCred),
		os.Getenv(projectName),
		os.Getenv(nodeType),
		os.Getenv(organizationName),
		os.Getenv(iaasProvider),
		os.Getenv(iaasRegion),
	)
	executeCommand(cmd, fmt.Sprintf("Fail to create project %q.", os.Getenv(projectName)))
}

func isProjectExists() bool {
	_, projectNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc project get -project-name %s",
		getProjectName(),
	)).Output()

	return projectNotExists == nil
}
