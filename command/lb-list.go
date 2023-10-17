package command

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/ui"
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
	projectUUID := projectUUIDFlag(c.flagSet)
	projectName := projectNameFlag(c.flagSet)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if c.flagSet.NArg() > 0 {
		return c.errorWithUsage(fmt.Errorf("Unparsed arguments on the command line: %v", c.flagSet.Args()))
	}

	if *projectUUID == "" && *projectName == "" {
		return c.errorWithUsage(errors.New("Project name or uuid is mandatory"))
	}

	return c.runWithSpinner("load balancer config", endpoint.String(), func(client *squarescale.Client) (string, error) {
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

		loadBalancers, err := client.LoadBalancerGet(UUID)
		if err != nil {
			return "", err
		}

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		// reset by ui/table.go FormatTable function: table.SetAutoFormatHeaders(false)
		// seems like this should be taken into account earlier than in the ui/table.go FormatTable function to have effect on fields
		table.SetAutoWrapText(false)
		table.SetHeader([]string{"Active", "Certificate Body", "HTTPS", "Public URL"})
		var activeIcon string
		var certBodyIcon string
		var httpsIcon string
		extraMsg := ""
		termType := os.Getenv("TERM_PROGRAM")
		matchTerm, _ := regexp.MatchString(".*[Tt][Mm][Uu][Xx].*", termType)
		if matchTerm {
			extraMsg = "\n\nPlease note that as you are using Tmux the UTF-8 icons might not be displayed properly unless you used the `-u` option"
		}
		for _, lb := range loadBalancers {
			if lb.Active {
				activeIcon = "✅"
			} else {
				activeIcon = "❌"
			}
			if lb.CertificateBody != "" {
				certBodyIcon = "✅"
			} else {
				certBodyIcon = "❌"
			}
			if lb.HTTPS {
				httpsIcon = "✅"
			} else {
				httpsIcon = "❌"
			}
			table.Append([]string{
				activeIcon,
				certBodyIcon,
				httpsIcon,
				lb.PublicURL,
			})
		}
		ui.FormatTable(table)

		table.Render()
		// Remove trailing \n and HT
		return string(regexp.MustCompile(`[\n\x09][\n\x09]*$`).ReplaceAll([]byte(tableString.String()), []byte(""))) + extraMsg, nil
	})
}

// Synopsis is part of cli.Command implementation.
func (c *LBListCommand) Synopsis() string {
	return "Display project's list of load balancers"
}

// Help is part of cli.Command implementation.
func (c *LBListCommand) Help() string {
	helpText := `
usage: sqsc lb list [options]

  Display load balancer list for given project.
`
	return strings.TrimSpace(helpText + optionsFromFlags(c.flagSet))
}
