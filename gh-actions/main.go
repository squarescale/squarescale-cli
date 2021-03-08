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
)

func main() {
	checkEnvironmentVariablesExists()

	createProject()
	createDatabase()
	createWebService()
	openHTTPPort()
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

func createProject() {
	_, projectNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc project get -project-name %s/%s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
	)).Output()

	if projectNotExists != nil {
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
		_, err := exec.Command("/bin/sh", "-c", cmd).Output()

		if err != nil {
			fmt.Println(cmd)
			log.Fatal(fmt.Sprintf("Creating project fails with error:\n %s", err))
		}
	}
}

func createDatabase() {
	_, dbEngineExists := os.LookupEnv(dbEngine)
	_, dbEngineVersionExists := os.LookupEnv(dbEngineVersion)
	_, dbEngineSizeExists := os.LookupEnv(dbSize)

	if dbEngineExists && dbEngineVersionExists && dbEngineSizeExists {
		_, databaseNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
			"/sqsc db show -project-name %s/%s | grep \"DB enabled\" | grep true",
			os.Getenv(organizationName),
			os.Getenv(projectName),
		)).Output()

		if databaseNotExists != nil {
			fmt.Println("Creating database...")

			cmd := fmt.Sprintf(
				"/sqsc db set -project-name %s/%s -engine \"%s\" -engine-version \"%s\" -size \"%s\" -yes",
				os.Getenv(organizationName),
				os.Getenv(projectName),
				os.Getenv(dbEngine),
				os.Getenv(dbEngineVersion),
				os.Getenv(dbSize),
			)
			_, err := exec.Command("/bin/sh", "-c", cmd).Output()

			if err != nil {
				fmt.Println(cmd)
				log.Fatal(fmt.Sprintf("Creating database fails with error:\n %s", err))
			}
		}
	} else {
		fmt.Println(fmt.Sprintf("%s, %s, %s are not set. No database will be created.", dbEngine, dbEngineVersion, dbSize))
	}
}

func createWebService() {
	_, webServiceNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc container list --project-name %s/%s | grep %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
	)).Output()

	if webServiceNotExists != nil {
		fmt.Println("Creating web service...")

		cmd := fmt.Sprintf(
			"/sqsc container add -project-name %s/%s -servicename %s -name %s:%s -username %s -password %s",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			os.Getenv(webServiceName),
			os.Getenv(dockerRepository),
			os.Getenv(dockerRepositoryTag),
			os.Getenv(dockerUser),
			os.Getenv(dockerToken),
		)
		_, err := exec.Command("/bin/sh", "-c", cmd).Output()

		if err != nil {
			fmt.Println(cmd)
			log.Fatal(fmt.Sprintf("Creating web service fails with error:\n %s", err))
		}
	}
}

func openHTTPPort() {
	networkRuleName := "http"

	_, webServiceNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc network-rule list -project-name %s/%s -service-name %s | grep %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
		networkRuleName,
	)).Output()

	if webServiceNotExists != nil {
		fmt.Println("Opening http port...")

		cmd := fmt.Sprintf(
			"/sqsc network-rule create -project-name %s/%s -external-protocol http -internal-port 80 -internal-protocol http -name %s -service-name %s",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			networkRuleName,
			os.Getenv(webServiceName),
		)
		_, err := exec.Command("/bin/sh", "-c", cmd).Output()

		if err != nil {
			fmt.Println(cmd)
			log.Fatal(fmt.Sprintf("Opening http port fails with error:\n%s", err))
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
