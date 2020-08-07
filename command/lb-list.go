package command

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/ui"
)

// LBListCommand gets the URL of the load balancer associated to a projects and prints it on the standard output.
type LBListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LBListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	return c.runWithSpinner("list load balancer config", endpoint.String(), func(client *squarescale.Client) (string, error) {
		loadBalancers, err := client.LoadBalancerList(*projectUUID)
		if err != nil {
			return "", err
		}

		return fmtLoadBalancersListOutput(loadBalancers), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *LBListCommand) Synopsis() string {
	return "Display project's public URL if available"
}

// Help is part of cli.Command implementation.
func (c *LBListCommand) Help() string {
	helpText := `
usage: sqsc lb list [options]

  Display load balancer state (enabled, disabled). In case the load
  balancer is enabled, all the project containers are displayed together
  with their ports.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func fmtLoadBalancersListOutput(loadBalancers []squarescale.LoadBalancer) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"Public URL"})
	data := make([][]string, len(loadBalancers), len(loadBalancers))
	for i, loadBalancer := range loadBalancers {
		data[i] = []string{loadBalancer.PublicUrl}
	}

	table.AppendBulk(data)

	ui.FormatTable(table)

	table.Render()

	return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte("")))
}
