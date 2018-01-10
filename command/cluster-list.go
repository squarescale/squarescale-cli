package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// ClusterListCommand is a cli.Command implementation to list the cluster node sizes available for use.
type ClusterListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *ClusterListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("list available cluster node sizes", endpoint.String(), func(client *squarescale.Client) (string, error) {
		nodeSizes, err := client.GetClusterNodeSizes()
		if err != nil {
			return "", err
		}

		msg := fmtClusterListOutput("Available node sizes", nodeSizes.ListHuman())

		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *ClusterListCommand) Synopsis() string {
	return "List all available cluster node sizes"
}

// Help is part of cli.Command implementation.
func (c *ClusterListCommand) Help() string {
	helpText := `
usage: sqsc cluster list

  List all cluster node sizes available for use
  in the Squarescale services.

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

func fmtClusterListOutput(title string, lines []string) string {
	var res string
	for _, line := range lines {
		res += fmt.Sprintf("  %s\n", line)
	}

	res = strings.Trim(res, "\n")
	return fmt.Sprintf("%s\n\n%s", title, res)
}
