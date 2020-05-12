package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command"
)

// func main() {
// 	log.SetOutput(os.Stdout)
//
// 	logLevel := os.Getenv("LOG_LEVEL")
// 	switch logLevel {
// 	case "debug":
// 		log.SetLevel(log.DebugLevel)
// 	case "info":
// 		log.SetLevel(log.InfoLevel)
// 	default:
// 		log.SetLevel(log.WarnLevel)
// 	}
//
// 	os.Exit(Run(os.Args[1:]))
// }

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	ui := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}
	cmds := command.Map(ui)
	var names []string
	for c := range cmds {
		names = append(names, c)
	}

	cli := &cli.CLI{
		Args:         os.Args[1:],
		Commands:     cmds,
		Autocomplete: true,
		Name:         "sqsc",
		HelpFunc:     cli.FilteredHelpFunc(names, cli.BasicHelpFunc("sqsc")),
		HelpWriter:   os.Stdout,
		// ErrorWriter:  os.Stderr,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
