package main

import (
	"fmt"
	"os"
)

const (
	sqscToken           = "SQSC_TOKEN"
	dockerUser          = "DOCKER_USER"
	dockerToken         = "DOCKER_TOKEN"
	dockerRepository    = "DOCKER_REPOSITORY"
	dockerRepositoryTag = "DOCKER_REPOSITORY_TAG"
	organizationName    = "ORGANIZATION_NAME"
	projectName         = "PROJECT_NAME"
	iaasProvider        = "IAAS_PROVIDER"
	iaasRegion          = "IAAS_REGION"
	iaasCred            = "IAAS_CRED"
	nodeType            = "NODE_TYPE"
	dbEngine            = "DB_ENGINE"
	dbEngineVersion     = "DB_ENGINE_VERSION"
	dbSize              = "DB_SIZE"
	servicesEnv         = "SERVICES"
	batchesEnv          = "BATCHES"
)

func main() {
	checkEnvironmentVariablesExists()

	project := Project{}
	project.create()

	database := Database{}
	database.create()

	services := Services{}
	services.create()

	batches := Batches{}
	batches.create()
}

func checkEnvironmentVariablesExists() {
	fmt.Println("Checking environment variables...")

	envVars := []string{
		sqscToken,
		dockerUser,
		dockerToken,
		dockerRepository,
		dockerRepositoryTag,
		organizationName,
		projectName,
		iaasProvider,
		iaasRegion,
		iaasCred,
		nodeType,
	}

	for _, envVar := range envVars {
		if _, exists := os.LookupEnv(envVar); !exists {
			fmt.Println(fmt.Sprintf("%s is not set. Quitting.", envVar))
			os.Exit(1)
		}
	}
}
