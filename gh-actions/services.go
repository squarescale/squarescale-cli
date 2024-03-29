package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/alessio/shellescape"
)

type Services struct{}

type ServiceContent struct {
	IMAGE_NAME     string                `json:"image_name"`
	IS_PRIVATE     bool                  `json:"is_private"`
	IMAGE_USER     string                `json:"image_user"`
	IMAGE_PASSWORD string                `json:"image_password"`
	RUN_CMD        string                `json:"run_cmd"`
	INSTANCES      string                `json:"instances"`
	NETWORK_RULES  []NetworkRulesContent `json:"network_rules"`
	ENV            map[string]string     `json:"env"`
	LIMIT_MEMORY   string                `json:"limit_memory"`
	LIMIT_CPU      string                `json:"limit_cpu"`
}

type NetworkRulesContent struct {
	NAME              string `json:"name"`
	INTERNAL_PORT     string `json:"internal_port"`
	DOMAIN            string `json:"domain"`
	PATH_PREFIX       string `json:"path_prefix"`
	INTERNAL_PROTOCOL string `json:"internal_protocol"`
	EXTERNAL_PROTOCOL string `json:"external_protocol"`
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
			s.insertServiceEnvAndLimits(serviceName, serviceContent)
			s.insertNetworkRules(serviceName, serviceContent)
			s.schedule(serviceName)
		}
	}
}

func (s *Services) createService(serviceName string, serviceContent ServiceContent) {
	fmt.Println(fmt.Sprintf("Creating %q service...", serviceName))

	cmd := fmt.Sprintf("%s container add", getSqscCLICmd())
	cmd += " -project-name " + getProjectName()
	cmd += " -servicename " + serviceName

	imageName := serviceContent.IMAGE_NAME
	if imageName == "" {
		imageName = getDockerImage()
	}
	cmd += " -name " + imageName

	runCommand := serviceContent.RUN_CMD
	if runCommand != "" {
		cmd += " -run-command " + shellescape.Quote(serviceContent.RUN_CMD)
	}

	instances := serviceContent.INSTANCES
	if instances != "" {
		cmd += " -instances " + instances
	}

	if serviceContent.IS_PRIVATE {
		cmd += " -username " + serviceContent.IMAGE_USER
		cmd += " -password " + serviceContent.IMAGE_PASSWORD
	}

	executeCommand(cmd, fmt.Sprintf("Failed to create %q service.", serviceName))
}

func (s *Services) insertServiceEnvAndLimits(serviceName string, serviceContent ServiceContent) {

	limitMemory := serviceContent.LIMIT_MEMORY
	limitCpu := serviceContent.LIMIT_CPU
	instances := serviceContent.INSTANCES
	command := serviceContent.RUN_CMD
	environment := serviceContent.ENV

	if len(environment) != 0 || limitMemory != "" || limitCpu != "" {

		cmd := fmt.Sprintf("%s container set", getSqscCLICmd())
		cmd += " -project-name " + getProjectName()
		cmd += " -service " + serviceName

		if len(environment) != 0 {
			fmt.Println(fmt.Sprintf("Inserting environment variable to service %q", serviceName))

			jsonFileName := "serviceEnvVar.json"
			env, _ := json.Marshal(environment)
			jsonErr := ioutil.WriteFile(jsonFileName, []byte(mapDatabaseEnv(string(env))), os.ModePerm)

			if jsonErr != nil {
				log.Fatal("Cannot write json file with map environment variables.")
			}

			cmd += " -env " + jsonFileName
		}

		if limitMemory != "" {
			cmd += " -memory " + limitMemory
		}

		if limitCpu != "" {
			cmd += " -cpu " + limitCpu
		}

		if instances != "" {
			cmd += " -instances " + instances
		}

		if command != "" {
			cmd += " -command " + command
		}

		executeCommand(cmd, fmt.Sprintf("Failed to import %q service environment variables.", serviceName))
	}
}

func (s *Services) insertNetworkRules(serviceName string, serviceContent ServiceContent) {
	for _, networkRules := range serviceContent.NETWORK_RULES {
		if (NetworkRulesContent{}) != networkRules {
			networkRuleName := networkRules.NAME

			if networkRuleName == "" {
				networkRuleName = "http"
			}

			if !isNetworkRuleExists(networkRuleName, serviceName) {
				internalPort := networkRules.INTERNAL_PORT
				if internalPort == "" {
					internalPort = "80"
				}

				fmt.Println(fmt.Sprintf("Opening %s port (http)...", internalPort))

				pathPrefix := networkRules.PATH_PREFIX
				if pathPrefix == "" {
					pathPrefix = "/"
				}

				internalProtocol := networkRules.INTERNAL_PROTOCOL
				if internalProtocol == "" {
					internalProtocol = "http"
				}

				externalProtocol := networkRules.EXTERNAL_PROTOCOL
				if externalProtocol == "" {
					externalProtocol = "http"
				}

				domain := networkRules.DOMAIN

				cmd := fmt.Sprintf("%s network-rule create", getSqscCLICmd())
				cmd += " -project-name " + getProjectName()
				cmd += " -external-protocol " + externalProtocol
				cmd += " -internal-port " + internalPort
				cmd += " -internal-protocol " + internalProtocol
				cmd += " -name " + networkRuleName
				cmd += " -service-name " + serviceName
				cmd += " -path-prefix " + pathPrefix

				if domain != "" {
					cmd += " -domain-expression " + domain
				}

				executeCommand(cmd, fmt.Sprintf("Failed to insert %q network rule on %s port in %q service", networkRuleName, internalPort, serviceName))
			} else {
				fmt.Println(fmt.Sprintf("Network rule %s already exists in %q service.", networkRuleName, serviceName))
			}
		}
	}
}

func (s *Services) schedule(serviceName string) {
	fmt.Println(fmt.Sprintf("Scheduling %q service...", serviceName))

	cmd := fmt.Sprintf(
		"%s service schedule --project-name %s %s",
		getSqscCLICmd(),
		getProjectName(),
		serviceName,
	)
	executeCommand(cmd, fmt.Sprintf("Failed to schedule %q service.", serviceName))
}

func isServiceExists(serviceName string) bool {
	_, webServiceNotExists := exec.Command("/bin/sh", "-c", fmt.Sprintf(
		"%s container list -project-name %s | grep %s",
		getSqscCLICmd(),
		getProjectName(),
		serviceName,
	)).Output()

	return webServiceNotExists == nil
}
