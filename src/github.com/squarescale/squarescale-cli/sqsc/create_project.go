package main

import (
	"fmt"
	"github.com/squarescale/squarescale-cli/squarescale"
	"github.com/squarescale/squarescale-cli/tokenstore"
	"os"
)

var cmdCreateProject = &Command{
	Run:       runCreateProject,
	UsageLine: "create-project ",
	Short:     "Create Project in Squarescale",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdLogin.Flag.BoolVar(&flagA, "a", false, "")
}

// runLogin executes login command and return exit code.
func runCreateProject(args []string) int {
	var project_name string
	if len(args) > 0 {
		project_name = args[0]
	}

	token, err := tokenstore.GetToken(SquarescaleEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	project_name, err = squarescale.CreateProject(SquarescaleEndpoint, token, project_name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Printf("Created project %s\n", project_name)

	return 0
}
