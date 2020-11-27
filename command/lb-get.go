package command

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LBGetCommand gets the URL of the load balancer associated to a projects and prints it on the standard output.
type LBGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LBGetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	projectName := c.flagSet.String("project-name", "", "set the name of the project")
	lbIDFlag := c.flagSet.String("lb-id", "", "id of the targeted load balancer inside project")
	var lbID int64
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if *lbIDFlag == "" {
		return c.errorWithUsage(errors.New("Load balancer ID is mandatory"))
	} else {
		lbIDint, _ := strconv.Atoi(*lbIDFlag)
		lbID = int64(lbIDint)
	}

	return c.runWithSpinner("load balancer config", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var UUID string
		var err error
		var projectToShow string
		if *projectUUID == "" {
			projectToShow = *projectName
			UUID, err = client.ProjectByName(*projectName)
			if err != nil {
				return "", err
			}
		} else {
			projectToShow = *projectUUID
			UUID = *projectUUID
		}

		loadBalancers, err := client.LoadBalancerGet(UUID)
		if err != nil {
			return "", err
		}

		var msg string
		for _, lb := range loadBalancers {
			if lb.ID == lbID {
				msg += fmt.Sprintf("Public URL: %s\n", lb.PublicURL)
				if lb.Active {
					msg += fmt.Sprintf("Active: ✅\n")
				} else {
					msg += fmt.Sprintf("Active: ❌\n")
				}
				if lb.HTTPS {
					msg += fmt.Sprintf("HTTPS: ✅\n")
				} else {
					msg += fmt.Sprintf("HTTPS: ❌\n")
				}
				if lb.CertificateBody != "" {
					msg += fmt.Sprintf("Certificate body:\n%s\n", lb.CertificateBody)
				} else {
					msg += fmt.Sprintf("Certificate body: ❌\n")
				}
				return msg, nil
			}
		}
		return "", errors.New(fmt.Sprintf("Load balancer with ID %d not found in project %s", lbID, projectToShow))
	})
}

// Synopsis is part of cli.Command implementation.
func (c *LBGetCommand) Synopsis() string {
	return "Display project's public URL if available"
}

// Help is part of cli.Command implementation.
func (c *LBGetCommand) Help() string {
	helpText := `
usage: sqsc lb get [options]

  Display load balancer state (enabled, disabled). In case the load
  balancer is enabled, all the project containers are displayed together
  with their ports.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
