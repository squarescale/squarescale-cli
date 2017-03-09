package command

import (
	"github.com/mitchellh/cli"
	"strings"
)

func AskYesNo(ui cli.Ui, alwaysYes *bool, question string, defVal bool) (bool, error) {
	if defVal {
		question = question + " [Y/n]"
	} else {
		question = question + " [y/N]"
	}

	if *alwaysYes {
		ui.Info(question + " y")
		return true, nil
	}

	ans, err := ui.Ask(question)
	if err != nil {
		return defVal, err
	}

	ans = strings.ToLower(ans)
	if defVal {
		return ans != "n", nil
	} else {
		return ans != "y", nil
	}
}
