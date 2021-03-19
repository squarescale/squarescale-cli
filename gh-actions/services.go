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
	IMAGE_NAME     string              `json:"image_name"`
	IS_PRIVATE     bool                `json:"is_private"`
	IMAGE_USER     string              `json:"image_user"`
	IMAGE_PASSWORD string              `json:"image_password"`
	RUN_CMD        string              `json:"run_cmd"`
	NETWORK_RULES  NetworkRulesContent `json:"network_rules"`
	ENV            map[string]string   `json:"env"`
}

type NetworkRulesContent struct {
	NAME          string `json:"name"`
	INTERNAL_PORT string `json:"internal_port"`
	DOMAIN        string `json:"domain"`
	PATH_PREFIX   string `json:"path_prefix"`
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
				s.insertServiceEnv(serviceName, serviceContent)
				s.insertNetworkRules(serviceName, serviceContent)
			} else {
				fmt.Println(fmt.Sprintf("Service %q already exists.", serviceName))
			}

			s.schedule(serviceName)
		}
	}
}

func (s *Services) createService(serviceName string, serviceContent ServiceContent) {
	fmt.Println(fmt.Sprintf("Creating %q service...", serviceName))

	cmd := "/sqsc container add"
	cmd += " -project-name " + getProjectName()
	cmd += " -servicename " + serviceName
	cmd += " -run-command " + serviceContent.RUN_CMD

	imageName := serviceContent.IMAGE_NAME
	if imageName == "" {
		imageName = getDockerImage()
	}
	cmd += " -name " + imageName

	if serviceContent.IS_PRIVATE {
		cmd += " -username " + serviceContent.IMAGE_USER
		cmd += " -password " + serviceContent.IMAGE_PASSWORD
	}

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

			pathPrefix := serviceContent.NETWORK_RULES.PATH_PREFIX
			if pathPrefix == "" {
				pathPrefix = "/"
			}

			domain := serviceContent.NETWORK_RULES.DOMAIN

			cmd := "/sqsc network-rule create"
			cmd += " -project-name " + getProjectName()
			cmd += " -external-protocol http"
			cmd += " -internal-port " + internalPort
			cmd += " -internal-protocol http"
			cmd += " -name " + networkRuleName
			cmd += " -service-name " + serviceName
			cmd += " -path-prefix " + pathPrefix

			if domain != "" {
				cmd += " -domain-expression " + domain
			}

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
