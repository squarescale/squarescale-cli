package main

import (
	"fmt"
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
			createDatabase()
		} else {
			fmt.Println("Database already exists.")
		}
	} else {
		fmt.Println(fmt.Sprintf("%s, %s, %s are not set. No database will be created.", dbEngine, dbEngineVersion, dbSize))
	}
}

func createDatabase() {
	fmt.Println("Creating database...")

	cmd := fmt.Sprintf(
		"/sqsc db set -project-name %s/%s -engine \"%s\" -engine-version \"%s\" -size \"%s\" -yes",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(dbEngine),
		os.Getenv(dbEngineVersion),
		os.Getenv(dbSize),
	)
	fmt.Println(cmd)
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()
	fmt.Println(string(output))

	if err != nil {
		fmt.Println(cmd)
		log.Fatal(fmt.Sprintf("Creating database fails with error:\n %s", err))
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
	value, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc env get -project-name %s/%s \"%s\" | grep -v %s | tr -d '\n'",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		key,
		"...done",
	)).Output()

	if err != nil {
		fmt.Println(fmt.Sprintf("Environment variable %q does not exists in this project.", key))
		return ""
	}

	return string(value)
}

func mapDatabaseEnv(env string) string {
	return mustache.Render(os.Getenv(mapEnvVar), map[string]string{
		"DB_HOST":     getSQSCEnvValue("DB_HOST"),
		"DB_USERNAME": getSQSCEnvValue("DB_USERNAME"),
		"DB_PASSWORD": getSQSCEnvValue("DB_PASSWORD"),
	})
}
