package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func executeCommand(cmd string, errorMsg string) {
	fmt.Println(cmd)
	output, err := exec.Command("/bin/sh", "-c", cmd).Output()

	if err != nil {
		fmt.Println(string(output))
		log.Fatal(errorMsg)
	}

	fmt.Println(string(output))
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

func isNetworkRuleExists(networkRuleName string, serviceName string) bool {
	_, networkRuleNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc network-rule list -project-name %s/%s -service-name %s | grep %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		serviceName,
		networkRuleName,
	)).Output()

	return networkRuleNotExists == nil
}
