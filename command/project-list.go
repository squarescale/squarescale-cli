package command

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

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
func (c *ProjectListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("list projects", endpoint.String(), func(client *squarescale.Client) (string, error) {
		projects, err := client.ListProjects()
		if err != nil {
			return "", err
		}

		organizations, err := client.ListOrganizations()
		if err != nil {
			return "", err
		}

		var msg string

		if len(projects) == 0 {
			msg = "No projects found"
		} else {
			msg = fmtProjectListOutput(projects, organizations)
		}

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectListCommand) Synopsis() string {
	return "List projects"
}

// Help is part of cli.Command implementation.
func (c *ProjectListCommand) Help() string {
	helpText := `
usage: sqsc project list [options]

  List projects attached to the authenticated account.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func fmtProjectListOutput(projects []squarescale.Project, organizations []squarescale.Organization) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "UUID", "Monitoring", "Provider", "Region", "Organization", "Status", "Size", "Slack Webhook"})
	data := make([][]string, len(projects), len(projects))

	for i, project := range projects {
		monitoring := ""
		if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
			monitoring = project.MonitoringEngine
		}
		data[i] = []string{
			project.Name,
			project.UUID,
			monitoring,
			project.Provider,
			project.Region,
			project.Organization,
			project.InfraStatus,
			fmt.Sprintf("%d/%d", project.NomadNodesReady, project.ClusterSize),
			project.SlackWebHook,
		}
	}

	table.AppendBulk(data)

	for _, o := range organizations {
		for _, project := range o.Projects {
			monitoring := ""
			if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
				monitoring = project.MonitoringEngine
			}
			table.Append([]string{
				fmt.Sprintf("%s/%s", o.Name, project.Name),
				project.UUID,
				monitoring,
				project.Provider,
				project.Region,
				project.Organization,
				project.InfraStatus,
				fmt.Sprintf("%d/%d", project.NomadNodesReady, project.ClusterSize),
				project.SlackWebHook,
			})
		}
	}

	ui.FormatTable(table)

	table.Render()
	// Remove trailing \n and HT
	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}
