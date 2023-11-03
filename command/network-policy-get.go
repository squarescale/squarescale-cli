package command

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/squarescale/squarescale-cli/squarescale"
)

// NetworkPolicyGetCommand is a cli.Command implementation for get current network policies
type NetworkPolicyGetCommand struct {
	Meta
	flagSet *flag.FlagSet
}

func (c *NetworkPolicyGetCommand) Run(args []string) int {
	c.flagSet = newFlagSet(c, c.Ui)
	endpoint := endpointFlag(c.flagSet)
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	jsonFormat := jsonFormatFlag(c.flagSet)
	dumpFlag := networkPolicyDumpFlag(c.flagSet)

	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if c.flagSet.NArg() > 1 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	version := networkPolicyVersionArg(c.flagSet)

	return c.runWithSpinner("show network policy status", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		policy, err := client.GetNetworkPolicy(UUID, version)
		if err != nil {
			return "", err
		}

		if *dumpFlag && policy.IsLoaded() {
			// dump YAML into file
			data, err := policy.DumpYAMLRules()

			if err != nil {
				return "", err
			}
			err = ioutil.WriteFile(policy.Version+".yaml", data, 0644)
			if err != nil {
				return "", err
			}
		}

		if *jsonFormat {
			if policy.IsLoaded() {
				j, _ := json.Marshal(policy)
				return string(j), nil
			}
			return "{}", nil
		}

		policy_status := "none"
		if policy.IsLoaded() {
			policy_status = policy.String()
		}

		if *dumpFlag && policy.IsLoaded() {
			policy_status = fmt.Sprintf(
				"%s\n(network policy rules dumped to %s file)\n",
				policy_status,
				policy.Version+".yaml",
			)
		}

		return policy_status, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *NetworkPolicyGetCommand) Synopsis() string {
	return "Get network policy for a project"
}

// Help is part of cli.Command implementation.
func (c *NetworkPolicyGetCommand) Help() string {
	helpText := `
usage: sqsc network-policy get [options] [version]

  Get the current active network policy or the specified version
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
