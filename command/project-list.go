package command

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/BenJetson/humantime"
	"github.com/olekukonko/tablewriter"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/ui"
)

// ProjectListCommand is a cli.Command implementation for listing all projects.
type ProjectListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *ProjectListCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	return cmd.runWithSpinner("list projects", endpoint.String(), func(client *squarescale.Client) (string, error) {
		projects, err := client.ListProjects()
		if err != nil {
			return "", err
		}

		projectCount := len(projects)

		organizations, err := client.ListOrganizations()
		if err != nil {
			return "", err
		}

		for io, o := range organizations {
			for ip, _ := range o.Projects {
				projectCount++
				if organizations[io].Projects[ip].Organization != o.Name {
					organizations[io].Projects[ip].Organization = o.Name
				}
			}
		}

		var msg string

		if projectCount == 0 {
			msg = "No projects found"
		} else {
			msg = fmtProjectListOutput(projects, organizations)
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *ProjectListCommand) Synopsis() string {
	return "List projects"
}

// Help is part of cli.Command implementation.
func (cmd *ProjectListCommand) Help() string {
	helpText := `
usage: sqsc project list [options]

  List projects attached to the authenticated account.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}

func fmtProjectListOutput(projects []squarescale.Project, organizations []squarescale.Organization) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
	// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Name", "UUID", "Monitoring", "Provider", "Credentials", "Region", "Organization", "Status", "Cluster", "Extra", "Hybrid", "Size", "Created", "Updated", "External ElasticSearch", "Slack Webhook"})

	location, _ := time.LoadLocation(time.Now().Location().String())

	for _, project := range projects {
		monitoring := ""
		if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
			monitoring = project.MonitoringEngine
		}
		isHybrid := "false"
		if project.HybridClusterEnabled {
			isHybrid = "true"
		}
		table.Append([]string{
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
			fmt.Sprintf("%s (%s)", project.CreatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.CreatedAt)),
			fmt.Sprintf("%s (%s)", project.UpdatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.UpdatedAt)),
			project.ExternalES,
			project.SlackWebHook,
		})
	}

	for _, o := range organizations {
		for _, project := range o.Projects {
			monitoring := ""
			if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
				monitoring = project.MonitoringEngine
			}
			isHybrid := "false"
			if project.HybridClusterEnabled {
				isHybrid = "true"
			}
			table.Append([]string{
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
				fmt.Sprintf("%s (%s)", project.CreatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.CreatedAt)),
				fmt.Sprintf("%s (%s)", project.UpdatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.UpdatedAt)),
				project.ExternalES,
				project.SlackWebHook,
			})
		}
	}

	ui.FormatTable(table)

	table.Render()
	// Remove trailing \n and HT
	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}
