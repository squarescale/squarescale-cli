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
	table.SetHeader([]string{"Name", "UUID", "Monitoring", "Provider", "Credentials", "Region", "Organization", "Status", "Nodes", "Extra", "Size", "Created", "Age", "External ElasticSearch", "Slack Webhook"})
	data := make([][]string, len(projects), len(projects))

	location, _ := time.LoadLocation(time.Now().Location().String())

	for i, project := range projects {
		monitoring := ""
		if project.MonitoringEnabled && len(project.MonitoringEngine) > 0 {
			monitoring = project.MonitoringEngine
		}
		data[i] = []string{
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
			project.NodeSize,
			project.CreatedAt.In(location).Format("2006-01-02 15:04"),
			humantime.Since(project.CreatedAt),
			project.ExternalES,
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
				project.NodeSize,
				project.CreatedAt.In(location).Format("2006-01-02 15:04"),
				humantime.Since(project.CreatedAt),
			})
		}
	}

	ui.FormatTable(table)

	table.Render()
	// Remove trailing \n and HT
	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}
