package list

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command/flags"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

func New(ui cli.Ui) *cmd {
	c := &cmd{UI: ui}
	c.init()
	return c
}

type cmd struct {
	UI       cli.Ui
	flags    *flag.FlagSet
	sqsc     *flags.SQSCFlags
	help     string
	endpoint string
	project  string
}

func (c *cmd) init() {
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	// c.flags.StringVar(&c.project, "project", "", "Project name")

	c.sqsc = &flags.SQSCFlags{}
	flags.Merge(c.flags, c.sqsc.CommonFlags())
	flags.Merge(c.flags, c.sqsc.ProjectFlags())

	c.help = flags.Usage(help, c.flags)
}

func (c *cmd) Run(args []string) int {
	if err := c.flags.Parse(args); err != nil {
		return 1
	}

	endPoint := c.sqsc.EndPoint()
	project, err := c.sqsc.Project()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error on parsing parameters: %s", err))
		return 1
	}

	token, _ := tokenstore.GetToken(endPoint)

	client := squarescale.NewClient(endPoint, token)

	volumes, err := client.GetVolumes(project)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error on get volume list: %s", err))
		return 1
	}

	var msg string = "Name\tSize\tType\tZone\t\tNode\tStatus\n"
	for _, v := range volumes {
		msg += fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\n", v.Name, v.Size, v.Type, v.Zone, v.StatefullNodeName, v.Status)
	}

	if len(volumes) == 0 {
		msg = "No volumes found"
	}

	fmt.Printf(msg)
	fmt.Printf("\n")
	return 0
}

func (c *cmd) Synopsis() string {
	return synopsis
}

func (c *cmd) Help() string {
	return c.help
}

const synopsis = "Lists volume from a project"
const help = `
Usage: sqsc volume list [options]

  To retrieve the volume list for the project named "foo":

      $ sqsc volume list -project foo
`
