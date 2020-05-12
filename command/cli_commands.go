package command

import (
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command/volume"
	volumelist "github.com/squarescale/squarescale-cli/command/volume/list"
)

func init() {
	Register("volume", func(cli.Ui) (cli.Command, error) { return volume.New(), nil })
	Register("volume list", func(ui cli.Ui) (cli.Command, error) { return volumelist.New(ui), nil })
}
