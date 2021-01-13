package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectCreateCommand is a cli.Command implementation for creating project.
type ProjectCreateCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectCreateCommand) Run(args []string) int {
	// Parse flags
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	alwaysYes := yesFlag(c.flagSet)
	name := projectNameFlag(c.flagSet)
	uuid := c.flagSet.String("uuid", "", "set the uuid of the project")
	organization := c.flagSet.String("organization", "", "set the organization the project will belongs to")
	provider := c.flagSet.String("provider", "", "set the cloud provider (aws or azure)")
	region := c.flagSet.String("region", "", "set the cloud provider region (eu-west-1)")
	credential := c.flagSet.String("credential", "", "set the credential used to build the infrastructure")
	infraType := c.flagSet.String("infra-type", "high-availability", "Set the infrastructure configuration (high-availability or single-node)")
	monitoringEngine := c.flagSet.String("monitoring", "", "Set the monitoring configuration (netdata)")
	nodeSize := c.flagSet.String("node-size", "", "Set the cluster node size")
	slackURL := c.flagSet.String("slackbot", "", "Set the Slack webhook URL")

	dbEngine := c.flagSet.String("db-engine", "", "Select database engine")
	dbSize := c.flagSet.String("db-size", "", "Select database size")
	dbVersion := c.flagSet.String("db-version", "", "Select database version")
	nowait := nowaitFlag(c.flagSet)

	payload := squarescale.JSONObject{}

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *name == "" && c.flagSet.Arg(0) != "" {
		nameArg := c.flagSet.Arg(0)
		name = &nameArg
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()[1:]))
	}

	if *name == "" {
		return c.errorWithUsage(errors.New("Project name is mandatory"))
	}

	payload["name"] = *name

	if *uuid != "" {
		payload["uuid"] = *uuid
	}

	if *organization != "" {
		payload["organization_name"] = *organization
	}

	if *provider == "" {
		return c.errorWithUsage(errors.New("Cloud provider is mandatory"))
	}

	payload["provider"] = *provider

	if *region == "" {
		return c.errorWithUsage(errors.New("Cloud provider region is mandatory"))
	}

	payload["region"] = *region

	if *credential != "" {
		payload["credential_name"] = *credential
	}

	if *nodeSize == "" {
		return c.errorWithUsage(errors.New("node size is mandatory"))
	}

	payload["node_size"] = *nodeSize

	if *infraType != "high-availability" && *infraType != "single-node" {
		return c.errorWithUsage(fmt.Errorf("Unknown infrastructure type: %v. Correct values are high-availability or single-node", *infraType))
	} else if *infraType == "high-availability" {
		payload["infra_type"] = "high_availability"
	} else {
		payload["infra_type"] = "single_node"
	}

	if *monitoringEngine != "" && *monitoringEngine != "netdata" {
		return c.errorWithUsage(fmt.Errorf("Unknown monitoring engine: %v. Correct values are empty-string (aka '') or netdata", *monitoringEngine))
	} else if *monitoringEngine == "netdata" {
		payload["monitoring"] = map[string]string{"engine": "netdata", "enabled": "true"}
	} else {
		payload["monitoring"] = map[string]string{"engine": "", "enabled": "false"}
	}

	if *dbEngine != "" && *dbSize == "" {
		return c.errorWithUsage(errors.New("if db engine is present, db size must be set"))
	} else if *dbEngine == "" && *dbSize != "" {
		return c.errorWithUsage(errors.New("if db size is present, db engine must be set"))
	} else if *dbEngine != "" && *dbSize != "" {
		payload["databases"] = []map[string]string{{"engine": *dbEngine, "size": *dbSize, "version": *dbVersion}}
	}

	if *slackURL != "" {
		payload["slack_webhook"] = *slackURL
	}

	// ask confirmation
	c.Ui.Warn("Project will be created with the following configuration :")
	c.Ui.Warn(fmt.Sprintf("name : %s", *name))
	if *uuid != "" {
		c.Ui.Warn(fmt.Sprintf("uuid : %s", *uuid))
	}
	if *organization != "" {
		c.Ui.Warn(fmt.Sprintf("organization : %s", *organization))
	}
	c.Ui.Warn(fmt.Sprintf("cloud provider : %s", *provider))
	c.Ui.Warn(fmt.Sprintf("cloud provider region : %s", *region))
	c.Ui.Warn(fmt.Sprintf("credential : %s", *credential))
	c.Ui.Warn(fmt.Sprintf("node size : %s", *nodeSize))
	c.Ui.Warn(fmt.Sprintf("infra type : %s", *infraType))
	c.Ui.Warn(fmt.Sprintf("monitoring engine : %s", *monitoringEngine))
	if *slackURL != "" {
		c.Ui.Warn(fmt.Sprintf("Slack webhook : %s", *slackURL))
	}
	if *dbEngine != "" && *dbSize != "" {
		c.Ui.Warn(fmt.Sprintf("database engine : %s", *dbEngine))
		c.Ui.Warn(fmt.Sprintf("database size : %s", *dbSize))
		if *dbVersion != "" {
			c.Ui.Warn(fmt.Sprintf("database version : %s", *dbVersion))
		}
	}

	ok, err := AskYesNo(c.Ui, alwaysYes, "Proceed ?", true)
	if err != nil {
		return c.error(err)
	} else if !ok {
		return c.error(CancelledError)
	}

	var project squarescale.Project

	res := c.runWithSpinner("create project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		project, err = client.CreateProject(&payload)
		return fmt.Sprintf("Created project '%s' with UUID '%s'", project.Name, project.UUID), err
	})

	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for project creation", endpoint.String(), func(client *squarescale.Client) (string, error) {
			projectStatus, err := client.WaitProject(project.UUID, 5)
			if err != nil {
				return projectStatus, err
			} else {
				return projectStatus, nil
			}
		})
	}

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectCreateCommand) Synopsis() string {
	return "Create Project"
}

// Help is part of cli.Command implementation.
func (c *ProjectCreateCommand) Help() string {
	helpText := `
usage: sqsc project create [options]

  Creates a new project using the provided options.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func (c *ProjectCreateCommand) askConfirmation(alwaysYes *bool, name string) error {
	c.pauseSpinner()

	c.startSpinner()
	return nil
}
