package command

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LBSetCommand configures the load balancer associated to a given project.
type LBSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (cmd *LBSetCommand) Run(args []string) int {
	cmd.flagSet = newFlagSet(cmd, cmd.Ui)
	endpoint := endpointFlag(cmd.flagSet)
	projectUUID := projectUUIDFlag(cmd.flagSet)
	projectName := projectNameFlag(cmd.flagSet)
	disableArg := loadBalancerDisableFlag(cmd.flagSet)
	certArg := certFlag(cmd.flagSet)
	certChainArg := certChainFlag(cmd.flagSet)
	secretKeyArg := secretKeyFlag(cmd.flagSet)
	if err := cmd.flagSet.Parse(args); err != nil {
		return 1
	}

	if cmd.flagSet.NArg() > 0 {
		return cmd.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", cmd.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return cmd.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	var cert, secretKey string
	var certChain string

	// Read certificate files and re-encode them for backend

	if *certArg != "" {
		data, err := ioutil.ReadFile(*certArg)
		if err != nil {
			return cmd.error(err)
		}
		cert = string(data)
	}

	if *certChainArg != "" {
		data, err := ioutil.ReadFile(*certChainArg)
		if err != nil {
			return cmd.error(err)
		}
		certChain = string(data)
	}

	if *secretKeyArg != "" {
		data, err := ioutil.ReadFile(*secretKeyArg)
		if err != nil {
			return cmd.error(err)
		}
		secretKey = string(data)
	}

	res := cmd.runWithSpinner("configure load balancer", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		if *disableArg {
			err = client.LoadBalancerDisable(UUID, loadBalancers[0].ID)
			return fmt.Sprintf("Successfully disable load balancer for project '%s'", projectToShow), err
		}

		err = client.LoadBalancerEnable(UUID, loadBalancers[0].ID, cert, certChain, secretKey)
		return fmt.Sprintf("Successfully update load balancer for project '%s'", projectToShow), err
	})

	return res
}

// Synopsis is part of cli.Command implementation.
func (cmd *LBSetCommand) Synopsis() string {
	return "Configure load balancer associated to project"
}

// Help is part of cli.Command implementation.
func (cmd *LBSetCommand) Help() string {
	helpText := `
usage: sqsc lb set [options]

  Configure load balancer associated to project.

	Used to define a TLS certificate
`
	return strings.TrimSpace(helpText + optionsFromFlags(cmd.flagSet))
}
