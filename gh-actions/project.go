package main

import (
	"fmt"
	"log"
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
	fmt.Println(cmd)
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	fmt.Println(string(output))

	if err != nil {
		fmt.Println(cmd)
		log.Fatal(fmt.Sprintf("Creating project fails with error:\n %s", err))
	}
}

func isProjectExists() bool {
	_, projectNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc project get -project-name %s/%s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
	)).Output()

	return projectNotExists == nil
}
