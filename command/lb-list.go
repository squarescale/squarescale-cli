package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
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
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return c.runWithSpinner("load balancer config", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		loadBalancers, err := client.LoadBalancerGet(UUID)
		if err != nil {
			return "", err
		}

		var msg string
		var activeIcon string
		var certBodyIcon string
		var httpsIcon string
		msg += fmt.Sprintf("ID\tActive\tCertificateBody\tHTTPS\tPublicURL\n")
		msg += fmt.Sprintf("--\t------\t---------------\t-----\t---------\n")
		for _, lb := range loadBalancers {
			if lb.Active {
				activeIcon = "✅"
			} else {
				activeIcon = "❌"
			}
			if lb.CertificateBody != "" {
				certBodyIcon = "✅"
			} else {
				certBodyIcon = "❌"
			}
			if lb.HTTPS {
				httpsIcon = "✅"
			} else {
				httpsIcon = "❌"
			}
			msg += fmt.Sprintf("%d\t%s\t%s\t\t%s\t%s\n", lb.ID, activeIcon, certBodyIcon, httpsIcon, lb.PublicURL)
		}
		return msg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *LBListCommand) Synopsis() string {
	return "Display project's public URL if available"
}

// Help is part of cli.Command implementation.
func (c *LBListCommand) Help() string {
	helpText := `
usage: sqsc lb get [options]

  Display load balancer state (enabled, disabled). In case the load
  balancer is enabled, all the project containers are displayed together
  with their ports.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
