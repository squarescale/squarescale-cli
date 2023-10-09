package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/BenJetson/humantime"
	"github.com/olekukonko/tablewriter"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/ui"
)

// ProjectDetailsCommand is a cli.Command implementation for listing project details.
type ProjectDetailsCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ProjectDetailsCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("detail project", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		projectDetails, err := client.GetProjectDetails(UUID)
		if err != nil {
			return "", err
		}
		if projectDetails == nil {
			return "No details to show", nil
		}
		return ProjectSummary(projectDetails), nil
	})
}

// Return project summary info like in front Overview page
func ProjectSummary(project *squarescale.ProjectWithAllDetails) string {
	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf("========== Project Summary: %s [%s]\n", project.Project.Name, project.Project.Organization))
	table := tablewriter.NewWriter(tableString)
	// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
	// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
	table.SetAutoWrapText(false)

	location, _ := time.LoadLocation(time.Now().Location().String())

	table.Append([]string{
		fmt.Sprintf("Provider: %s/%s", project.Project.Infrastructure.CloudProvider, project.Project.Infrastructure.CloudProviderLabel),
		fmt.Sprintf("Region: %s/%s", project.Project.Infrastructure.Region, project.Project.Infrastructure.RegionLabel),
		fmt.Sprintf("Credentials: %s", project.Project.Infrastructure.CredentialName),
	})
	table.Append([]string{
		fmt.Sprintf("Infrastructure type: %s", project.Project.Infrastructure.Type),
		fmt.Sprintf("Size: %s", project.Project.Infrastructure.NodeSize),
		fmt.Sprintf("Hybrid: %v", project.Project.HybridClusterEnabled),
	})
	table.Append([]string{
		fmt.Sprintf("Created: %s (%s)", project.Project.CreatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.Project.CreatedAt)),
		fmt.Sprintf("External ElasticSearch: %s", project.Project.ExternalElasticSearch),
		"",
	})
	table.Append([]string{
		fmt.Sprintf("Updated: %s (%s)", project.Project.UpdatedAt.In(location).Format("2006-01-02 15:04"), humantime.Since(project.Project.UpdatedAt)),
		fmt.Sprintf("Slack WebHook: %s", project.Project.SlackWebHook),
		"",
	})

	ui.FormatTable(table)

	table.Render()
	return tableString.String()
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectDetailsCommand) Synopsis() string {
	return "Get project informations"
}

// Help is part of cli.Command implementation.
func (c *ProjectDetailsCommand) Help() string {
	helpText := `
usage: sqsc project details [options]

  Get project detailed informations.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
