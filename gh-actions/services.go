package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Services struct{}

type ServiceContent struct {
	RUN_CMD       string              `json:"run_cmd"`
	NETWORK_RULES NetworkRulesContent `json:"network_rules"`
	ENV           map[string]string   `json:"env"`
}

type NetworkRulesContent struct {
	NAME          string `json:"name"`
	INTERNAL_PORT string `json:"internal_port"`
}

func (s *Services) create() {
	if _, exists := os.LookupEnv(servicesEnv); exists {
		var services map[string]ServiceContent

		err := json.Unmarshal([]byte(os.Getenv(servicesEnv)), &services)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error when Unmarshal %s", err))
		}

		for serviceName, serviceContent := range services {
			if !isServiceExists(serviceName) {
				s.createService(serviceName, serviceContent)
			} else {
				fmt.Println(fmt.Sprintf("Service %q already exists.", serviceName))
			}

			s.insertServiceEnv(serviceName, serviceContent)
			s.insertNetworkRules(serviceName, serviceContent)
			s.schedule(serviceName)
		}
	}
}

func (s *Services) createService(serviceName string, serviceContent ServiceContent) {
	fmt.Println(fmt.Sprintf("Creating %q service...", serviceName))

	cmd := fmt.Sprintf(
		"/sqsc container add -project-name %s -servicename %s -name %s:%s -username %s -password %s -run-command \"%s\"",
		getProjectName(),
		serviceName,
		os.Getenv(dockerRepository),
		os.Getenv(dockerRepositoryTag),
		os.Getenv(dockerUser),
		os.Getenv(dockerToken),
		serviceContent.RUN_CMD,
	)
	executeCommand(cmd, fmt.Sprintf("Fail to create service %q.", serviceName))
}

func (s *Services) insertServiceEnv(serviceName string, serviceContent ServiceContent) {
	if len(serviceContent.ENV) != 0 {
		fmt.Println(fmt.Sprintf("Inserting environment variable to service %q", serviceName))

		jsonFileName := "serviceEnvVar.json"
		env, _ := json.Marshal(serviceContent.ENV)
		jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(string(env))), os.ModePerm)

		if jsonErr != nil {
			log.Fatal("Cannot write json file with map environment variables.")
		}

		instancesNumber := "1"

		cmd := fmt.Sprintf(
			"/sqsc container set -project-name %s -env %s -service %s -instances %s -command \"%s\"",
			getProjectName(),
			jsonFileName,
			serviceName,
			instancesNumber,
			serviceContent.RUN_CMD,
		)
		executeCommand(cmd, fmt.Sprintf("Fail to import service %q environment variables.", serviceName))
	}
}

func (s *Services) insertNetworkRules(serviceName string, serviceContent ServiceContent) {
	if (NetworkRulesContent{}) != serviceContent.NETWORK_RULES {
		networkRuleName := serviceContent.NETWORK_RULES.NAME

		if networkRuleName == "" {
			networkRuleName = "http"
		}

		if !isNetworkRuleExists(networkRuleName, serviceName) {
			internalPort := serviceContent.NETWORK_RULES.INTERNAL_PORT

			if internalPort == "" {
				internalPort = "80"
			}

			fmt.Println(fmt.Sprintf("Opening %s port (http)...", internalPort))

			cmd := fmt.Sprintf(
				"/sqsc network-rule create -project-name %s -external-protocol http -internal-port %s -internal-protocol http -name %s -service-name %s",
				getProjectName(),
				internalPort,
				networkRuleName,
				serviceName,
			)
			executeCommand(cmd, fmt.Sprintf("Fail to insert %q network rule on %s port", networkRuleName, internalPort))
		} else {
			fmt.Println(fmt.Sprintf("Network rule %s already exists.", networkRuleName))
		}
	}
}

func (s *Services) schedule(serviceName string) {
	fmt.Println(fmt.Sprintf("Scheduling %q service...", serviceName))

	cmd := fmt.Sprintf(
		"/sqsc service schedule --project-name %s %s",
		getProjectName(),
		serviceName,
	)
	executeCommand(cmd, fmt.Sprintf("Fail to schedule %q service.", serviceName))
}

func isServiceExists(serviceName string) bool {
	_, webServiceNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"/sqsc container list -project-name %s | grep %s",
		getProjectName(),
		serviceName,
	)).Output()

	return webServiceNotExists == nil
}
