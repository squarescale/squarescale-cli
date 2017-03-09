package command

import (
	"errors"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

var CancelledError error = errors.New("Cancelled")

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
	if err == CancelledError {
		return m.cancelled()
	} else {
		m.Ui.Error(err.Error())
		return 1
	}
}

func (m *Meta) cancelled() int {
	return 2
}

func (m *Meta) errorWithUsage(err error) int {
	m.error(errors.New(err.Error() + "\n"))
	return cli.RunResultHelp
}

func (m *Meta) runWithSpinner(text, endpoint string, action func(*squarescale.Client) (string, error)) int {
	m.startSpinner()
	m.spin.Suffix = " " + text

	client, err := m.ensureLogin(endpoint)
	if err != nil {
		m.errorSpinner(err)
		return m.error(err)
	}

	finalMsg, err := action(client)
	if err != nil {
		m.errorSpinner(err)
		return m.error(err)
	}

	m.stopSpinner()
	return m.info(finalMsg)
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

func (m *Meta) errorSpinner(err error) {
	if err == CancelledError {
		m.spin.FinalMSG = "... cancelled\n"
	} else {
		m.spin.FinalMSG = "... error\n"
	}
	m.spin.Stop()
}

func (m *Meta) ensureLogin(endpoint string) (*squarescale.Client, error) {
	token, err := tokenstore.GetToken(endpoint)
	if err != nil {
		return nil, err
	}

	client := squarescale.NewClient(endpoint, token)
	if err := client.ValidateToken(); err != nil {
		return nil, err
	}

	return client, nil
}
