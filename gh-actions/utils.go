package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const sqscCLICmd string = "/sqsc"

func getSqscCLICmd() string {
	envVars := []string{
		"SQSC_CLI_COLOR",
		"SQSC_CLI_FORMAT",
		"SQSC_CLI_PROGRESS",
	}
	ret := sqscCLICmd
	for _, envVar := range envVars {
		value, exists := os.LookupEnv(envVar)
		if exists && len(strings.TrimSpace(value)) > 0 {
			ret += fmt.Sprintf(" %s %s", strings.ToLower(strings.ReplaceAll(envVar, "SQSC_CLI_", "")), value)
		}
	}
	return ret
}

func checkEnvironmentVariablesExists() {
	fmt.Println("Checking environment variables...")

	envVars := []string{
		sqscToken,
		dockerRepository,
		projectName,
		iaasProvider,
		iaasRegion,
		iaasCred,
		nodeType,
		infraType,
	}

	for _, envVar := range envVars {
		value, exists := os.LookupEnv(envVar)
		if !exists || len(strings.TrimSpace(value)) == 0 {
			fmt.Println(fmt.Sprintf("%s is not set or is empty/filled with blanks. Quitting.", envVar))
			os.Exit(1)
		}
	}
}

func executeCommand(cmd string, errorMsg string) {
	fmt.Println(cmd)
	execCmd := exec.Command("/bin/sh", "-c", cmd)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stdout
	err := execCmd.Run()

	if err != nil {
		log.Fatal(errorMsg)
	}
}

func getProjectName() string {
	if _, exists := os.LookupEnv(organizationName); exists {
		return fmt.Sprintf("%s/%s", os.Getenv(organizationName), os.Getenv(projectName))
	}

	return fmt.Sprintf("%s", os.Getenv(projectName))
}

func getDockerImage() string {
	if _, exists := os.LookupEnv(dockerRepositoryTag); exists {
		return fmt.Sprintf("%s:%s", os.Getenv(dockerRepository), os.Getenv(dockerRepositoryTag))
	}

	return fmt.Sprintf("%s", os.Getenv(dockerRepository))
}

func getSQSCEnvValue(key string) string {
	value, err := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"%s env get -project-name %s \"%s\" | grep -v %s | tr -d '\n'",
		getSqscCLICmd(),
		getProjectName(),
		key,
		"...done",
	)).Output()

	if err != nil {
		fmt.Println(fmt.Sprintf("Environment variable %q does not exists in this project.", key))
		return ""
	}

	return string(value)
}

func isNetworkRuleExists(networkRuleName string, serviceName string) bool {
	_, networkRuleNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"%s network-rule list -project-name %s -service-name %s | grep %s",
		getSqscCLICmd(),
		getProjectName(),
		serviceName,
		networkRuleName,
	)).Output()

	return networkRuleNotExists == nil
}
