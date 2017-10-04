package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

var CancelledError error = errors.New("Cancelled")
var IsTTY bool

func init() {
	IsTTY = isatty.IsTerminal(os.Stdout.Fd())
}

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	Ui         cli.Ui
	spin       *spinner.Spinner
	spinEnable bool
	niceFormat bool
}

// DefaultMeta returns a default meta object with an initialized spinner.
func DefaultMeta(ui cli.Ui, color, niceFormat, spinnerEnable bool, spinTime time.Duration) *Meta {
	if color {
		ui = &cli.ColoredUi{
			InfoColor:  cli.UiColorCyan,
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui:         ui,
		}
	}

	if spinTime == 0 {
		spinTime = 100 * time.Millisecond
	}

	return &Meta{
		Ui:         ui,
		spin:       spinner.New(spinner.CharSets[7], spinTime),
		spinEnable: spinnerEnable,
		niceFormat: niceFormat,
	}
}

func (m *Meta) info(message string, args ...interface{}) int {
	m.Ui.Info(fmt.Sprintf(message, args...))
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
	if finalMsg != "" {
		return m.info(finalMsg)
	} else {
		return 0
	}
}

func (m *Meta) startSpinner() {
	if !m.spinEnable {
		return
	}

	m.spin.Start()
}

func (m *Meta) pauseSpinner() {
	if !m.spinEnable {
		return
	}

	m.spin.FinalMSG = "... paused\n"
	m.spin.Stop()
	time.Sleep(time.Millisecond) // leave time to the UI to refresh properly
}

func (m *Meta) stopSpinner() {
	if !m.spinEnable {
		return
	}

	m.spin.FinalMSG = "... done\n"
	m.spin.Stop()
}

func (m *Meta) errorSpinner(err error) {
	if !m.spinEnable {
		return
	}

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

func (m *Meta) FormatTable(table string, header bool) string {
	table = strings.Trim(table, "\n")
	if !m.niceFormat {
		return table
	}

	return FormatTable(table)
}

func FormatTable(table string) string {
	table = strings.Trim(table, "\n")
	var res string
	lines := strings.Split(table, "\n")
	var csize []int
	for _, line := range lines {
		cols := strings.Split(line, "\t")
		for c, col := range cols {
			if len(csize) <= c {
				csize = append(csize, len(col))
			} else if len(col) > csize[c] {
				csize[c] = len(col)
			}
		}
	}
	for _, line := range lines {
		if res != "" {
			res += "\n"
		}
		cols := strings.Split(line, "\t")
		var l string
		for c, col := range cols {
			sz := csize[c]
			l += col
			if c < len(cols)-1 {
				l += strings.Repeat(" ", sz-len(col))
				l += "\t"
			}
		}
		res += l
	}
	return res
}
