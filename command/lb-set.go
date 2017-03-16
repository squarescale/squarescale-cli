package command

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

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
	nowait := nowaitFlag(c.flagSet)
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

	var taskId int

	res := c.runWithSpinner("configure load balancer", *endpoint, func(client *squarescale.Client) (string, error) {
		var err error
		if *disabledArg {
			taskId, err = client.DisableLB(*project)
			return fmt.Sprintf("Successfully disabled load balancer for project '%s'", *project), err
		}

		container, err := client.GetContainerInfo(*project, *containerArg)
		if err != nil {
			return "", err
		}

		if *portArg > 0 {
			container.WebPort = *portArg
		}

		msg := fmt.Sprintf(
			"Successfully configured load balancer (enabled = '%v', container = '%s', port = '%d') for project '%s'",
			true, *containerArg, container.WebPort, *project)

		taskId, err = client.ConfigLB(*project, container.ID, container.WebPort)
		return msg, err
	})
	if res != 0 {
		return res
	}

	if !*nowait {
		res = c.runWithSpinner("wait for load balancer change", *endpoint, func(client *squarescale.Client) (string, error) {
			return client.WaitTask(taskId, time.Second)
		})
	}

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
