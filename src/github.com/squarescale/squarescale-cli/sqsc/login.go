package main

import (
	"fmt"
	"github.com/squarescale/squarescale-cli/github"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
	"os"
)

var cmdLogin = &Command{
	Run:       runLogin,
	UsageLine: "login ",
	Short:     "Login to Squarescale",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdLogin.Flag.BoolVar(&flagA, "a", false, "")
}

type AuthTokenResponse struct {
	AuthToken string `json:"auth_token"`
}

// runLogin executes login command and return exit code.
func runLogin(args []string) int {

	login, password, one_time_password, gh_token, gh_token_url, err := github.GeneratePersonalToken("Squarescale CLI")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 1
	}

	fmt.Println("Forward GitHub authorization to Squarescale...")

	sqsc_token, err := squarescale.ObtainTokenFromGitHub(SquarescaleEndpoint, gh_token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 1
	}

	fmt.Printf("Store Squarescale token: %s\n", sqsc_token)

	err = tokenstore.SaveToken(SquarescaleEndpoint, sqsc_token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 1
	}

	fmt.Println("Revoke temporary GitHub token...")

	err = github.RevokePersonalToken(gh_token_url, login, password, one_time_password)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		return 1
	}

	return 0
}
