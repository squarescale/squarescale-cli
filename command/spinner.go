package command

import (
	"time"

	"github.com/briandowns/spinner"
)

func startSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
	s.Suffix = suffix
	s.Start()
	return s
}

func stopSpinner(s *spinner.Spinner) {
	s.Stop()
	time.Sleep(10 * time.Millisecond) // leave enough time for the Ui to refresh
}
