package command

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// LBSetCommand configures the load balancer associated to a given project.
type LBSetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

// Run is part of cli.Command implementation.
func (c *LBSetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := c.flagSet.String("project-uuid", "", "uuid of the targeted project")
	disabledArg := disabledFlag(c.flagSet, "Disable load balancer")
	certArg := certFlag(c.flagSet)
	certChainArg := certChainFlag(c.flagSet)
	secretKeyArg := secretKeyFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	var cert, secretKey string
	var certChain string

	// Read certificate files and re-encode them for backend

	if *certArg != "" {
		data, err := ioutil.ReadFile(*certArg)
		if err != nil {
			return c.error(err)
		}
		cert = string(data)
	}

	if *certChainArg != "" {
		data, err := ioutil.ReadFile(*certChainArg)
		if err != nil {
			return c.error(err)
		}
		certChain = string(data)
	}

	if *secretKeyArg != "" {
		data, err := ioutil.ReadFile(*secretKeyArg)
		if err != nil {
			return c.error(err)
		}
		secretKey = string(data)
	}

	res := c.runWithSpinner("configure load balancer", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		if *disabledArg {
			err = client.DisableLB(*projectUUID)
			return fmt.Sprintf("Successfully disabled load balancer for project '%s'", *projectUUID), err
		}

		err = client.ConfigLB(*projectUUID, cert, certChain, secretKey)

		return fmt.Sprintf("Successfully update load balancer for project '%s'", *projectUUID), err
	})

	return res
}

// Synopsis is part of cli.Command implementation.
func (c *LBSetCommand) Synopsis() string {
	return "Configure the load balancer associated to a project"
}

// Help is part of cli.Command implementation.
func (c *LBSetCommand) Help() string {
	helpText := `
usage: sqsc lb set [options]

  Configure the load balancer associated to a Squarescale project.

	Used to define a TLS certificate
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
