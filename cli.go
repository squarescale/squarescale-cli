package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command"
)

func defValueFromEnv(envname string, def bool) bool {
	env := strings.ToLower(os.Getenv(envname))
	if env == "0" || env == "false" || env == "off" || env == "no" {
		def = false
	} else if env != "" {
		def = true
	}
	return def
}

func defDurationFromEnv(envname string, def time.Duration) time.Duration {
	env, err := time.ParseDuration(os.Getenv(envname))
	if err == nil {
		return env
	} else {
		return def
	}
}

func Run(args []string) int {
	var f flag.FlagSet

	color := f.Bool("color", defValueFromEnv("SQSC_COLOR", command.IsTTY), "Colored output")
	format := f.Bool("format", defValueFromEnv("SQSC_FORMAT", true), "Enable nice output")
	spin := f.Bool("progress", defValueFromEnv("SQSC_PROGRESS", command.IsTTY), "Enable progress spinner")

	err := f.Parse(args)
	if err == flag.ErrHelp {
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}

	// Meta-option for executables.
	// It defines output color and its stdout/stderr stream.
	meta := command.DefaultMeta(&cli.BasicUi{
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
		Reader:      os.Stdin,
	}, *color, *format, *spin, defDurationFromEnv("SQSC_SPIN_TIME", 0))

	return RunCustom(f.Args(), Commands(meta))
}

func RunCustom(args []string, commands map[string]cli.CommandFactory) int {
	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	cli := &cli.CLI{
		Args:       args,
		Commands:   commands,
		Version:    Version,
		HelpFunc:   cli.BasicHelpFunc(Name),
		HelpWriter: os.Stdout,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute: %s\n", err.Error())
	}

	return exitCode
}
