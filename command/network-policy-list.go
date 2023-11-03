package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkPolicyListCommand is a cli.Command implementation for listing network policies
type NetworkPolicyListCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *NetworkPolicyListCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	jsonFormat := jsonFormatFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	return c.runWithSpinner("list network policy versions", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		current, versions, err := client.ListNetworkPolicies(UUID)
		if err != nil {
			return "", err
		}

		var buffer bytes.Buffer
		for idx, version := range versions {
			buffer.WriteString(version.Version)
			versions[idx].Status = squarescale.NETWORK_POLICY_DEFAULT_STATUS
			buffer.WriteString(" " + version.Name)
			if current.Version == version.Version {
				versions[idx].Active = true
				versions[idx].Status = current.StatusStr
				buffer.WriteString(" ** status: " + current.StatusStr)

			}
			buffer.WriteString("\n")
		}

		if len(versions) == 0 {
			buffer.WriteString("No network policies found")
		} else {
			if *jsonFormat {
				j, _ := json.Marshal(versions)
				return string(j), nil
			}
		}

		return buffer.String(), nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkPolicyListCommand) Synopsis() string {
	return "List network policies"
}

// Help is part of cli.Command implementation.
func (c *NetworkPolicyListCommand) Help() string {
	helpText := `
usage: sqsc network-policy list [options]

  list network policies related to a project
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
