package main

import (
	"fmt"
	"os"
	"os/exec"
)

type NetworkRule struct{}

const (
	networkRuleName     = "http"
	defaultInternalPort = "80"
)

func (nr *NetworkRule) create() {
	if !isNetworkRuleExists() {
		nr.openPort()
	} else {
		fmt.Println("Network rule already exists.")
	}
}

func (nr *NetworkRule) openPort() {
	internalPort := defaultInternalPort
	if _, internalPortEnvExists := os.LookupEnv(internalPortEnv); internalPortEnvExists {
		internalPort = os.Getenv(internalPortEnv)
	}

	fmt.Println(fmt.Sprintf("Opening %s port (http)...", internalPort))

	cmd := fmt.Sprintf(
		"/sqsc network-rule create -project-name %s/%s -external-protocol http -internal-port %s -internal-protocol http -name %s -service-name %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		internalPort,
		networkRuleName,
		os.Getenv(webServiceName),
	)
	executeCommand(cmd, "Fail to open http port.")
}

func isNetworkRuleExists() bool {
	_, networkRuleNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc network-rule list -project-name %s/%s -service-name %s | grep %s",
		os.Getenv(organizationName),
		os.Getenv(projectName),
		os.Getenv(webServiceName),
		networkRuleName,
	)).Output()

	return networkRuleNotExists == nil
}
