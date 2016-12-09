package command

import (
	"errors"
	"flag"
	"fmt"
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
	project := projectFlag(c.flagSet)
	containerArg := containerFlag(c.flagSet)
	portArg := portFlag(c.flagSet)
	disabledArg := disabledFlag(c.flagSet, "Disable load balancer")
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := validateLBSetCommandArgs(*project, *containerArg, *portArg, *disabledArg); err != nil {
		return c.errorWithUsage(err)
	}

	var msg string
	err := c.runWithSpinner("configure load balancer", *endpoint, func(token string) error {
		if *disabledArg {
			msg = fmt.Sprintf("Successfully disabled load balancer for project '%s'", *project)
			return squarescale.DisableLB(*endpoint, token, *project)
		}

		container, err := squarescale.GetContainerInfo(*endpoint, token, *project, *containerArg)
		if err != nil {
			return err
		}

		if *portArg > 0 {
			container.WebPort = *portArg
		}

		msg = fmt.Sprintf(
			"Successfully configured load balancer (enabled = '%v', container = '%s', port = '%d') for project '%s'",
			true, *containerArg, container.WebPort, *project)

		return squarescale.ConfigLB(*endpoint, token, *project, container.ID, container.WebPort)
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
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

`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}

// validateLBSetCommandArgs ensures that the following predicate is satisfied:
// - 'disabled' is true and (container and port are not both specified)
// - 'disabled' is false and (container or port are specified)
func validateLBSetCommandArgs(project, container string, port int, disabled bool) error {
	if err := validateProjectName(project); err != nil {
		return err
	}

	if disabled && (container != "" || port > 0) {
		return errors.New("Cannot specify container or port when disabling load balancer.")
	}

	if !disabled && (container == "" && port <= 0) {
		return errors.New("Instance and engine cannot be both empty.")
	}

	return nil
}