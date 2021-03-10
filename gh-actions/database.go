package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/hoisie/mustache"
)

type Database struct{}

func (d *Database) create() {
	_, dbEngineExists := os.LookupEnv(dbEngine)
	_, dbEngineVersionExists := os.LookupEnv(dbEngineVersion)
	_, dbEngineSizeExists := os.LookupEnv(dbSize)

	if dbEngineExists && dbEngineVersionExists && dbEngineSizeExists {
		if !isDabataseExists() {
			fmt.Println("Creating database...")

			createDatabase()
			insertDatabaseEnvironement()
		} else {
			fmt.Println("Database already exists.")
		}
	} else {
		fmt.Println(fmt.Sprintf("%s, %s, %s are not set. No database will be created.", dbEngine, dbEngineVersion, dbSize))
	}
}

func createDatabase() {
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

func insertDatabaseEnvironement() {
	if _, exists := os.LookupEnv(mapEnvVar); exists {
		fmt.Println("Inserting database environment variables...")

		data := mustache.Render(os.Getenv(mapEnvVar), map[string]string{
			"DB_HOST":     getSQSCEnvValue("DB_HOST"),
			"DB_USERNAME": getSQSCEnvValue("DB_USERNAME"),
			"DB_PASSWORD": getSQSCEnvValue("DB_PASSWORD"),
		})

		jsonFileName := "mapEnvVar.json"
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(data), os.ModePerm)

		if jsonErr != nil {
			log.Fatal("Cannot write json file with map environment variables.")
		}

		cmd := fmt.Sprintf(
			"/sqsc container set -project-name %s/%s -env %s -service %s -instances 1",
			os.Getenv(organizationName),
			os.Getenv(projectName),
			jsonFileName,
			os.Getenv(webServiceName),
		)
		_, cmdErr := exec.Command("/bin/sh", "-c", cmd).Output()

		if cmdErr != nil {
			fmt.Println(cmd)
			log.Fatal("Fail to import database environment variables.")
		}
	}
}

func isDabataseExists() bool {
	_, databaseNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc db show -project-name %s/%s | grep \"DB enabled\" | grep true",
		os.Getenv(organizationName),
		os.Getenv(projectName),
	)).Output()

	return databaseNotExists == nil
}

func getSQSCEnvValue(key string) string {
	value, _ := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc env get -project-name %s/%s \"%s\" | grep -v %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		key,
		"...done",
	)).Output()

	return string(value)
}
