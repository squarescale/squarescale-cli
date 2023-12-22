package command

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LBGetCommand gets the URL of the load balancer associated to a projects and prints it on the standard output.
type LBGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *LBGetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)

	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return cmd.runWithSpinner("load balancer config", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
		if len(loadBalancers) != 1 {
			log.Printf("Warning project %s has %d Load Balancers (only 1 expected)", UUID, len(loadBalancers))
		}

		var msg string
		for _, lb := range loadBalancers {
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
		return "", errors.New(fmt.Sprintf("No load balancer found in project %s", projectToShow))
	})
}

// Synopsis is part of cli.Command implementation.
func (cmd *LBGetCommand) Synopsis() string {
	return "Display project's public URL if available"
}

// Help is part of cli.Command implementation.
func (cmd *LBGetCommand) Help() string {
	helpText := `
usage: sqsc lb get [options]

  Display load balancer state (enabled, disabled). In case the load
  balancer is enabled, all the project containers are displayed together
  with their ports.
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
