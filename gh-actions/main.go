package main

import (
	"fmt"
	"os"
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
	batchesEnv          = "BATCHES"
)

func main() {
	fmt.Println("gh-action version 0.0.1")
	checkEnvironmentVariablesExists()

	project := Project{}
	project.create()

	webservice := Webservice{}
	webservice.create()

	networkRule := NetworkRule{}
	networkRule.create()

	database := Database{}
	database.create()

	serviceEnv := ServiceEnv{}
	serviceEnv.create()

	batches := Batches{}
	batches.create()

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
		if _, exists := os.LookupEnv(envVar); !exists {
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
	executeCommand(cmd, "Fail to schedule web service")
}
