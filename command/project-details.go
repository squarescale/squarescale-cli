package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"strconv"
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
	noSummary := projectDetailsNoSummaryFlag(c.flagSet)
	noComputeResources := projectDetailsNoComputeResourcesFlag(c.flagSet)
	noSchedulingGroups := projectDetailsNoSchedulingGroupsFlag(c.flagSet)
	noExternalNodes := projectDetailsNoExternalNodesFlag(c.flagSet)

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
		// TODO: add Main cluster + Extra nodes + Volumes sections
		res := ""
		if !*noSummary {
			res += ProjectSummary(projectDetails) + "\n"
		}
		if !*noComputeResources {
			res += ProjectComputeNodes(projectDetails) + "\n"
		}
		if !*noSchedulingGroups {
			res += ProjectSchedulingGroups(projectDetails) + "\n"
		}
		if !*noExternalNodes {
			res += ProjectExternalNodes(projectDetails) + "\n"
		}
		return res, nil
	})
}

// Return project summary info like in front Overview page
func ProjectSummary(project *squarescale.ProjectWithAllDetails) string {
	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf("========== Summary: %s [%s]\n", project.Project.Name, project.Project.Organization))
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

func findExternalNodeRef(name string, project *squarescale.ProjectWithAllDetails) *squarescale.ExternalNode {
	for _, n := range project.Project.Infrastructure.Cluster.ExternalNodes {
		if n.Name == name {
			return &n
		}
	}
	return nil
}

// Return project compute nodes like in front Compute Resources page
// See font/src/components/infrastructure/ComputeResourcesTab.jsx
func ProjectComputeNodes(project *squarescale.ProjectWithAllDetails) string {
	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf("========== Compute resources: %s [%s]\n", project.Project.Name, project.Project.Organization))
	table := tablewriter.NewWriter(tableString)
	// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
	// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
	table.SetAutoWrapText(false)
	// TODO: add monitoring URLs
	table.SetHeader([]string{"Hostname", "Arch", "(v)CPUs", "Freq (Ghz)", "Mem (Gb)", "Type", "Disk (Gb)", "Free Disk (Gb)", "IP Address", "OS", "Status", "Mode", "Availability Zone", "Scheduling group", "Nomad", "Consul"})

	for _, c := range project.Project.Infrastructure.Cluster.ClusterMembersDetails {
		freq, _ := strconv.ParseFloat(c.CPUFrequency, 32)
		mem, _ := strconv.ParseFloat(c.Memory, 32)
		stt, _ := strconv.ParseFloat(c.StorageBytesTotal, 32)
		stf, _ := strconv.ParseFloat(c.StorageBytesFree, 32)
		extRef := findExternalNodeRef(c.Name, project)
		mode := "Cluster"
		zone := "N/A"
		fullStorage := ""
		freeStorage := fmt.Sprintf("%.0f (%.2f%%%%)", stf / 1024.0 / 1024.0 / 1024.0, (100.0 * stf) / stt)
		schedulingGroup := "None"
		if c.SchedulingGroup.Name != "" {
			schedulingGroup = c.SchedulingGroup.Name
		}
		if extRef != nil {
			mode = fmt.Sprintf("External: %s", c.Name)
			fullStorage = fmt.Sprintf("%.0f", stt / 1024.0 / 1024.0 / 1024.0)
		} else {
			fullStorage = fmt.Sprintf("%d",  project.Project.Infrastructure.Cluster.RootDiskSize)
			zone = c.Zone
		}
		table.Append([]string{
			c.Hostname,
			c.CPUArch,
			c.CPUCores,
			fmt.Sprintf("%.1f", freq / 1000.0 ),
			fmt.Sprintf("%.2f", mem / 1024.0 / 1024.0 / 1024.0),
			c.InstanceType,
			fullStorage,
			freeStorage,
			c.PrivateIP,
			fmt.Sprintf("%s %s", c.OSName, c.OSVersion),
			c.NomadStatus,
			mode,
			zone,
			schedulingGroup,
			c.NomadVersion,
			c.ConsulVersion,
		})
	}

	ui.FormatTable(table)

	table.Render()
	return tableString.String()
}

// Return project scheduling groups like in front Compute Resources page
// See font/src/components/infrastructure/ComputeResourcesTab.jsx
func ProjectSchedulingGroups(project *squarescale.ProjectWithAllDetails) string {
	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf("========== Scheduling groups: %s [%s]\n", project.Project.Name, project.Project.Organization))
	table := tablewriter.NewWriter(tableString)
	// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
	// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
	table.SetAutoWrapText(false)
	// TODO: add monitoring URLs
	table.SetHeader([]string{"Name", "# Nodes", "Nodes", "# Services", "Services"})

	for _, s := range project.Project.Infrastructure.Cluster.SchedulingGroups {
		nodes := ""
		for _, n := range s.ClusterMembers {
			nodes += "," + n.Name
		}
		if len(nodes) > 0 {
			nodes = nodes[1:]
		}
		services := ""
		for _, n := range s.Services {
			services += "," + n.Name
		}
		if len(services) > 0 {
			services = services[1:]
		}
		table.Append([]string{
			s.Name,
			fmt.Sprintf("%d", len(s.ClusterMembers)),
			nodes,
			fmt.Sprintf("%d", len(s.Services)),
			services,
		})
	}

	ui.FormatTable(table)

	table.Render()
	return tableString.String()
}

func findComputeNodeRef(name string, project *squarescale.ProjectWithAllDetails) *squarescale.ClusterMemberDetails {
	for _, n := range project.Project.Infrastructure.Cluster.ClusterMembersDetails {
		if n.Name == name {
			return &n
		}
	}
	return nil
}

// Return project external nodes like in front Compute Resources page
// See font/src/components/infrastructure/externalNodes/ExternalNodes.jsx
func ProjectExternalNodes(project *squarescale.ProjectWithAllDetails) string {
	tableString := &strings.Builder{}
	tableString.WriteString(fmt.Sprintf("========== External nodes: %s [%s]\n", project.Project.Name, project.Project.Organization))
	table := tablewriter.NewWriter(tableString)
	// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
	// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
	table.SetAutoWrapText(false)
	// TODO: add monitoring URLs
	table.SetHeader([]string{"Hostname", "Public IP", "Status", "Private Network"})

	for _, s := range project.Project.Infrastructure.Cluster.ExternalNodes {
		status := s.Status
		computeRef := findComputeNodeRef(s.Name, project)
		if computeRef != nil {
			if computeRef.NomadStatus == "ready" {
				status = "Connected"
			} else {
				status = "Not connected"
			}
		}
		table.Append([]string{
			s.Name,
			s.PublicIP,
			status,
			s.PrivateNetwork,
		})
	}

	ui.FormatTable(table)

	table.Render()
	return tableString.String()
}

// Synopsis is part of cli.Command implementation.
func (c *ProjectDetailsCommand) Synopsis() string {
	return "Get project detailed informations"
}

// Help is part of cli.Command implementation.
func (c *ProjectDetailsCommand) Help() string {
	helpText := `
usage: sqsc project details [options]

  Get project detailed informations.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
