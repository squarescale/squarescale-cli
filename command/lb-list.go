package command

import (
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
	project := projectFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := validateProjectName(*project); err != nil {
		return c.errorWithUsage(err)
	}

	var msg string
	err := c.runWithSpinner("list load balancer config", *endpoint, func(token string) error {
		enabled, err := squarescale.LoadBalancerEnabled(*endpoint, token, *project)
		if err != nil {
			return err
		}

		msg = "state: enabled"
		if !enabled {
			msg = "state: disabled"
			return nil
		}

		containers, err := squarescale.GetContainers(*endpoint, token, *project)
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			msg += "\n"
			return nil
		}

		checkedChars := map[bool]string{false: " ", true: "✓"}
		var msgBody []string
		for _, container := range containers {
			fmtMsg := fmt.Sprintf("[%s] %s:%d", checkedChars[container.Web], container.ShortName, container.WebPort)
			msgBody = append(msgBody, fmtMsg)
		}

		msg += " ([✓] -> maps to <container>:<port>, no mapping otherwise)\n\n"
		msg += strings.Join(msgBody, "\n")
		return nil
	})

	if err != nil {
		return c.error(err)
	}

	return c.info(msg)
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