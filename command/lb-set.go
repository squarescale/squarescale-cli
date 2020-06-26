package command

import (
	"encoding/pem"
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
	serviceArg := c.flagSet.String("service", "", "select the service")
	portArg := portFlag(c.flagSet)
	exprArg := exprFlag(c.flagSet)
	disabledArg := disabledFlag(c.flagSet, "Disable load balancer")
	httpsArg := httpsFlag(c.flagSet)
	certArg := certFlag(c.flagSet)
	certChainArg := certChainFlag(c.flagSet)
	secretKeyArg := secretKeyFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *disabledArg && (*serviceArg != "" || *portArg > 0) {
		return c.errorWithUsage(errors.New("Cannot specify service or port when disabling load balancer."))
	}

	if *projectUUID == "" {
		return c.errorWithUsage(errors.New("Project uuid is mandatory"))
	}

	var taskId int
	var cert, secretKey string
	var certChain []string

	// Read certificate files and re-encode them for backend

	if *certArg != "" {
		data, err := ioutil.ReadFile(*certArg)
		if err != nil {
			return c.error(err)
		}
		block, data := pem.Decode(data)
		if block == nil {
			return c.error(fmt.Errorf("%s: Could not decode PEM", *secretKeyArg))
		} else if block.Type != "CERTIFICATE" {
			return c.error(fmt.Errorf("%s: Does not contain a certificate, it contains %s", *secretKeyArg, block.Type))
		} else {
			cert = string(pem.EncodeToMemory(block))
		}
	}

	if *certChainArg != "" {
		data, err := ioutil.ReadFile(*certChainArg)
		if err != nil {
			return c.error(err)
		}
		block, data := pem.Decode(data)
		for block != nil {
			if block.Type != "CERTIFICATE" {
				return c.error(fmt.Errorf("%s: Does not contain a certificate, it contains %s", *secretKeyArg, block.Type))
			} else {
				certChain = append(certChain, string(pem.EncodeToMemory(block)))
			}
			block, data = pem.Decode(data)
		}
	}

	if *secretKeyArg != "" {
		data, err := ioutil.ReadFile(*secretKeyArg)
		if err != nil {
			return c.error(err)
		}
		block, data := pem.Decode(data)
		if block == nil {
			return c.error(fmt.Errorf("%s: Could not decode PEM", *secretKeyArg))
		} else if !strings.Contains(block.Type, "PRIVATE KEY") {
			return c.error(fmt.Errorf("%s: Does not contain a private key, it contains %s", *secretKeyArg, block.Type))
		} else {
			secretKey = string(pem.EncodeToMemory(block))
		}
	}

	res := c.runWithSpinner("configure load balancer", endpoint.String(), func(client *squarescale.Client) (string, error) {
		var err error
		if *disabledArg {
			taskId, err = client.DisableLB(*projectUUID)
			return fmt.Sprintf("[#%d] Successfully disabled load balancer for project '%s'", taskId, *projectUUID), err
		}

		container, err := client.GetServicesInfo(*projectUUID, *serviceArg)
		if err != nil {
			return "", err
		}

		if *portArg > 0 {
			container.WebPort = *portArg
		}

		taskId, err = client.ConfigLB(*projectUUID, container.ID, container.WebPort, *exprArg, *httpsArg, cert, certChain, secretKey)
		msg := fmt.Sprintf(
			"[#%d] Successfully configured load balancer (enabled = '%v', container = '%s', port = '%d') for project '%s'",
			taskId, true, *serviceArg, container.WebPort, *projectUUID)

		return msg, err
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
  "--container" and "--port" flags must be specified together.

  If you want access the service through a CNAME dns entry, you have
  to specify a regex expression that match this entry with "--expr".
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

// validateLBSetCommandArgs ensures that the following predicate is satisfied:
// - 'disabled' is true and (container and port are not both specified)
// - 'disabled' is false and (container or port are specified)
func validateLBSetCommandArgs(container string, port int, disabled bool, cert, certChain, secretKey string) error {

	return nil
}
