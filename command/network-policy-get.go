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

func (b *NetworkPolicyGetCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)
	jsonFormat := jsonFormatFlag(b.flagSet)
	dumpFlag := networkPolicyDumpFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	version := networkPolicyVersionArg(b.flagSet)

	return b.runWithSpinner("show network policy status", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *NetworkPolicyGetCommand) Synopsis() string {
	return "Get network policy for a project"
}

// Help is part of cli.Command implementation.
func (b *NetworkPolicyGetCommand) Help() string {
	helpText := `
usage: sqsc network-policy get [options] [version]

  Get the current active network policy or the specified version
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
