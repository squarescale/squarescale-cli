package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/BenJetson/humantime"
	"github.com/squarescale/squarescale-cli/squarescale"
)

// ProjectGetCommand is a cli.Command implementation for listing project basic infos.
type ProjectGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ProjectGetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("get project", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		if *projectUUID == "" {
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			UUID = *projectUUID
		}
		project, err := client.GetProject(UUID)
		if err != nil {
			return "", err
		}

		policy, err := client.GetNetworkPolicy(UUID, "")
		if err != nil {
			return "", err
		}
		policy_status := "none"
		if policy.IsLoaded() {
			policy_status = policy.String()
		}

		var msg string

		monitoring := ""
		if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
			monitoring = project.MonitoringEngine
		}
		isHybrid := "false"
		if project.HybridClusterEnabled {
			isHybrid = "true"
		}

		location, _ := time.LoadLocation(time.Now().Location().String())

		msg = fmt.Sprintf("Name: %s\nUUID: %s\nMonitoring: %s\nProvider: %s\nCredentials: %s\n"+
			"Region: %s\nOrganization: %s\nStatus: %s\nCluster: %s\n"+
			"Extra: %s\nHybrid: %s\nSize: %s\nRootDiskSize: %d GB\nCreated: %s\nUpdated: %s\n"+
			"External ElasticSearch: %s\nSlack Webhook: %s\n\n"+
			"* Network policies * \n%s",
			project.Name,
			project.UUID,
			monitoring,
			fmt.Sprintf("%s/%s", project.Provider, project.ProviderLabel),
			project.Credentials,
			fmt.Sprintf("%s/%s", project.Region, project.RegionLabel),
			project.Organization,
			project.ProjectStatus(),
			project.ProjectStateLessCount(),
			project.ProjectStateFulCount(),
			isHybrid,
			project.NodeSize,
			project.RootDiskSizeGB,
			fmt.Sprintf("%s (%s)", project.CreatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.CreatedAt)),
			fmt.Sprintf("%s (%s)", project.UpdatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.UpdatedAt)),
			project.ExternalES,
			project.SlackWebHook,
			policy_status,
		)

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ProjectGetCommand) Synopsis() string {
	return "Get project basic informations"
}

// Help is part of cli.Command implementation.
func (cmd *ProjectGetCommand) Help() string {
	helpText := `
usage: sqsc project get [options]

  Get project basic informations.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
