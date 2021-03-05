package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	checkEnvironmentVariablesExists()

	createDatabase()
	createProject()
	createWebService()
	openHTTPPort()
	scheduleWebService()
}

func checkEnvironmentVariablesExists() {
	envVars := []string{
		"SQSC_TOKEN",
		"ORGANIZATION_NAME",
		"PROJECT_NAME",
		"WEB_SERVICE_NAME",
		"DOCKER_USER",
		"DOCKER_TOKEN",
		"DOCKER_REPOSITORY",
		"DOCKER_REPOSITORY_TAG",
		"IAAS_CRED",
		"IAAS_PROVIDER",
		"IAAS_REGION",
		"NODE_TYPE",
	}

	for _, envVar := range envVars {
		_, exists := os.LookupEnv(envVar)
		if !exists {
			fmt.Println(fmt.Sprintf("%s is not set. Quitting.", envVar))
			os.Exit(1)
		}
	}
}

func createDatabase() {
	_, dbEngineExists := os.LookupEnv("DB_ENGINE")
	_, dbEngineVersionExists := os.LookupEnv("DB_ENGINE_VERSION")
	_, dbEngineSizeExists := os.LookupEnv("DB_SIZE")

	if !dbEngineExists && !dbEngineVersionExists && !dbEngineSizeExists {
		_, databaseNotExists := exec.Command("/bin/bash", "-c", fmt.Sprintf(
			"/sqsc db show -project-name %s/%s | grep \"DB enabled\" | grep true",
			os.Getenv("ORGANIZATION_NAME"),
			os.Getenv("PROJECT_NAME"),
		)).Output()

		if databaseNotExists != nil {
			fmt.Println("Creating database...")

			_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(
				"/sqsc db set -project-name %s -engine %s -engine-version %s -size %s -yes",
				os.Getenv("PROJECT_NAME"),
				os.Getenv("DB_ENGINE"),
				os.Getenv("DB_ENGINE_VERSION"),
				os.Getenv("DB_SIZE"),
			)).Output()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func createProject() {
	_, projectNotExists := exec.Command("/bin/bash", "-c", fmt.Sprintf(
		"/sqsc project get -project-name %s/%s",
		os.Getenv("ORGANIZATION_NAME"),
		os.Getenv("PROJECT_NAME"),
	)).Output()

	if projectNotExists != nil {
		fmt.Println("Creating project...")

		_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(
			"/sqsc project create -credential %s -monitoring netdata -name %s -node-size %s -infra-type high-availability -organization %s -provider %s -region %s -yes",
			os.Getenv("IAAS_CRED"),
			os.Getenv("PROJECT_NAME"),
			os.Getenv("NODE_TYPE"),
			os.Getenv("ORGANIZATION_NAME"),
			os.Getenv("IAAS_PROVIDER"),
			os.Getenv("IAAS_REGION"),
		)).Output()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func createWebService() {
	_, webServiceNotExists := exec.Command("/bin/bash", "-c", fmt.Sprintf(
		"/sqsc container list --project-name %s/%s | grep %s",
		os.Getenv("ORGANIZATION_NAME"),
		os.Getenv("PROJECT_NAME"),
		os.Getenv("WEB_SERVICE_NAME"),
	)).Output()

	if webServiceNotExists != nil {
		fmt.Println("Creating web service...")

		_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(
			"/sqsc container add -project-name %s/%s -servicename %s -name %s:%s -username %s -password %s",
			os.Getenv("ORGANIZATION_NAME"),
			os.Getenv("PROJECT_NAME"),
			os.Getenv("WEB_SERVICE_NAME"),
			os.Getenv("DOCKER_REPOSITORY"),
			os.Getenv("DOCKER_REPOSITORY_TAG"),
			os.Getenv("DOCKER_USER"),
			os.Getenv("DOCKER_TOKEN"),
		)).Output()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func openHTTPPort() {
	networkRuleName := "http"

	_, webServiceNotExists := exec.Command("/bin/bash", "-c", fmt.Sprintf(
		"/sqsc network-rule list -project-name %s/%s -service-name %s | grep %s",
		os.Getenv("ORGANIZATION_NAME"),
		os.Getenv("PROJECT_NAME"),
		os.Getenv("WEB_SERVICE_NAME"),
		networkRuleName,
	)).Output()

	if webServiceNotExists != nil {
		fmt.Println("Opening http port...")

		_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(
			"/sqsc network-rule create -project-name %s/%s -external-protocol http -internal-port 80 -internal-protocol http -name %s -service-name %s",
			os.Getenv("ORGANIZATION_NAME"),
			os.Getenv("PROJECT_NAME"),
			networkRuleName,
			os.Getenv("WEB_SERVICE_NAME"),
		)).Output()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func scheduleWebService() {
	_, err := exec.Command("/bin/bash", "-c", fmt.Sprintf(
		"/sqsc service schedule --project-name %s/%s %s",
		os.Getenv("ORGANIZATION_NAME"),
		os.Getenv("PROJECT_NAME"),
		os.Getenv("WEB_SERVICE_NAME"),
	)).Output()

	if err != nil {
		log.Fatal(err)
	}
}
