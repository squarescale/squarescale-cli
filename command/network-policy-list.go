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

func (b *NetworkPolicyListCommand) Run(args []string) int {
	b.flagSet = newFlagSet(b, b.Ui)
	endpoint := endpointFlag(b.flagSet)
	projectUUID := projectUUIDFlag(b.flagSet)
	projectName := projectNameFlag(b.flagSet)
	jsonFormat := jsonFormatFlag(b.flagSet)

	if err := b.flagSet.Parse(args); err != nil {
		return 1
	}

	if *projectUUID == "" && *projectName == "" {
		return b.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	if b.flagSet.NArg() > 1 {
		return b.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", b.flagSet.Args()))
	}

	return b.runWithSpinner("list network policy versions", endpoint.String(), func(client *squarescale.Client) (string, error) {
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
func (b *NetworkPolicyListCommand) Synopsis() string {
	return "List network policies"
}

// Help is part of cli.Command implementation.
func (b *NetworkPolicyListCommand) Help() string {
	helpText := `
usage: sqsc network-policy list [options]

  list network policies related to a project
`
	return strings.TrimSpace(helpText + optionsFromFlags(b.flagSet))
}
