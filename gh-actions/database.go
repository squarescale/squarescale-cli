package main

import (
	"fmt"
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
		if !isDatabaseExists() {
			d.createDatabase()
		} else {
			fmt.Println("Database already exists.")
		}
	} else {
		fmt.Println(fmt.Sprintf("%s, %s, %s are not set correctly. No database will be created.", dbEngine, dbEngineVersion, dbSize))
	}
}

func (d *Database) createDatabase() {
	fmt.Println("Creating database...")

	cmd := fmt.Sprintf(
		"%s db set -project-name %s -engine \"%s\" -engine-version \"%s\" -size \"%s\" -yes",
		getSqscCLICmd(),
		getProjectName(),
		os.Getenv(dbEngine),
		os.Getenv(dbEngineVersion),
		os.Getenv(dbSize),
	)
	executeCommand(cmd, fmt.Sprintf("Failed to create database %q %q %q for %q project.", os.Getenv(dbEngine), os.Getenv(dbEngineVersion), os.Getenv(dbSize), getProjectName()))
}

func isDatabaseExists() bool {
	_, databaseNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"%s db show -project-name %s | grep \"DB enabled\" | grep true",
		getSqscCLICmd(),
		getProjectName(),
	)).Output()

	return databaseNotExists == nil
}

func mapDatabaseEnv(env string) string {
	return mustache.Render(env, map[string]string{
		"DB_HOST":     getSQSCEnvValue("DB_HOST"),
		"DB_USERNAME": getSQSCEnvValue("DB_USERNAME"),
		"DB_PASSWORD": getSQSCEnvValue("DB_PASSWORD"),
		"DB_NAME":     getSQSCEnvValue("DB_NAME"),
		"DB_PORT":     getSQSCEnvValue("DB_PORT"),
	})
}
