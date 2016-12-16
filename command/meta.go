package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui   cli.Ui
	spin *spinner.Spinner
}

// DefaultMeta returns a default meta object with an initialized spinner.
func DefaultMeta(ui cli.Ui) *Meta {
	return &Meta{
		Ui:   ui,
		spin: spinner.New(spinner.CharSets[24], 100*time.Millisecond),
	}
}

func (m *Meta) info(message string) int {
	m.Ui.Info(message)
	return 0
}

func (m *Meta) error(err error) int {
	m.Ui.Error(fmt.Sprintf("Error: %v", err))
	return 1
}

func (m *Meta) errorWithUsage(err error) int {
	m.error(errors.New(err.Error() + "\n"))
	return cli.RunResultHelp
}

func (m *Meta) runWithSpinner(text, endpoint string, action func(*squarescale.Client) error) error {
	m.startSpinner()
	m.spin.Suffix = " " + text
	defer m.stopSpinner()

	client, err := m.ensureLogin(endpoint)
	if err != nil {
		return err
	}

	return action(client)
}

func (m *Meta) startSpinner() {
	m.spin.Start()
}

func (m *Meta) pauseSpinner() {
	m.spin.FinalMSG = "... paused\n"
	m.spin.Stop()
	time.Sleep(time.Millisecond) // leave time to the UI to refresh properly
}

func (m *Meta) stopSpinner() {
	m.spin.FinalMSG = "... done\n"
	m.spin.Stop()
}

func (m *Meta) ensureLogin(endpoint string) (*squarescale.Client, error) {
	token, err := tokenstore.GetToken(endpoint)
	if err != nil {
		return nil, err
	}

	client := squarescale.NewClient(endpoint, token)
	if err := client.ValidateToken(); err != nil {
		return nil, fmt.Errorf("You're not authenticated, please login first (%v)", err)
	}

	return client, nil
}
