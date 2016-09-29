package command

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
)

func runWithSpinner(text, endpoint string, action func(token string) error) error {
	s := spinner.New(spinner.CharSets[24], 100*time.Millisecond)
	defer s.Stop()
	s.Suffix = " " + text
	s.FinalMSG = "... done\n"
	s.Start()

	token, err := tokenstore.GetToken(endpoint)
	if err != nil {
		return err
	}

	err = squarescale.ValidateToken(endpoint, token)
	if err != nil {
		return fmt.Errorf("You're not authenticated, please login first (%v)", err)
	}

	return action(token)
}
