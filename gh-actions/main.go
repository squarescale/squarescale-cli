package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	sqscToken           = "SQSC_TOKEN"
	organizationName    = "ORGANIZATION_NAME"
	projectName         = "PROJECT_NAME"
	webServiceName      = "WEB_SERVICE_NAME"
	dockerUser          = "DOCKER_USER"
	dockerToken         = "DOCKER_TOKEN"
	dockerRepository    = "DOCKER_REPOSITORY"
	dockerRepositoryTag = "DOCKER_REPOSITORY_TAG"
	iaasCred            = "IAAS_CRED"
	iaasProvider        = "IAAS_PROVIDER"
	iaasRegion          = "IAAS_REGION"
	nodeType            = "NODE_TYPE"
	dbEngine            = "DB_ENGINE"
	dbEngineVersion     = "DB_ENGINE_VERSION"
	dbSize              = "DB_SIZE"
	cmdEnv              = "CMD"
	internalPortEnv     = "INTERNAL_PORT"
	mapEnvVar           = "MAP_ENV_VAR"
)

func main() {
	checkEnvironmentVariablesExists()

	project := Project{}
	project.create()

	webservice := Webservice{}
	webservice.create()

	networkRule := NetworkRule{}
	networkRule.create()

	database := Database{}
	database.create()

	scheduleWebService()
}

func checkEnvironmentVariablesExists() {
	fmt.Println("Checking environment variables...")

	envVars := []string{
		sqscToken,
		organizationName,
		projectName,
		webServiceName,
		dockerUser,
		dockerToken,
		dockerRepository,
		dockerRepositoryTag,
		iaasCred,
		iaasProvider,
		iaasRegion,
		nodeType,
	}

	for _, envVar := range envVars {
		_, exists := os.LookupEnv(envVar)
		if !exists {
			fmt.Println(fmt.Sprintf("%s is not set. Quitting.", envVar))
			os.Exit(1)
		}
	}
}

func scheduleWebService() {
	fmt.Println("Scheduling web service...")

	cmd := fmt.Sprintf(
		"/sqsc service schedule --project-name %s/%s %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
	)
	_, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		fmt.Println(cmd)
		log.Fatal(fmt.Sprintf("Scheduling web service fails with error:\n%s", err))
	}
}

func getCmdEnvValue() string {
	cmdEnvValue := ""
	if _, cmdExists := os.LookupEnv(cmdEnv); cmdExists {
		cmdEnvValue = os.Getenv(cmdEnv)
	}
	return cmdEnvValue
}
