package command

import (
	"time"

	"github.com/briandowns/spinner"
)

func startSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
	s.Suffix = " " + suffix
	s.FinalMSG = "... done\n"
	s.Start()
	return s
}
