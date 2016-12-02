package main

import (
	"github.com/mitchellh/cli"
	"github.com/squarescale/squarescale-cli/command"
)

// Commands returns a map of command factories.
func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &command.LoginCommand{
				Meta: *meta,
			}, nil
		},
		"project": func() (cli.Command, error) {
			return &command.ProjectCommand{}, nil
		},
		"project create": func() (cli.Command, error) {
			return &command.ProjectCreateCommand{
				Meta: *meta,
			}, nil
		},
		"project list": func() (cli.Command, error) {
			return &command.ProjectListCommand{
				Meta: *meta,
			}, nil
		},
		"url": func() (cli.Command, error) {
			return &command.UrlCommand{
				Meta: *meta,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &command.DBCommand{}, nil
		},
		"db scale": func() (cli.Command, error) {
			return &command.DBScaleCommand{
				Meta: *meta,
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &command.StatusCommand{
				Meta: *meta,
			}, nil
		},
		"repository": func() (cli.Command, error) {
			return &command.RepositoryCommand{}, nil
		},
		"repository add": func() (cli.Command, error) {
			return &command.RepositoryAddCommand{
				Meta: *meta,
			}, nil
		},
		"repository list": func() (cli.Command, error) {
			return &command.RepositoryListCommand{
				Meta: *meta,
			}, nil
		},
		"logs": func() (cli.Command, error) {
			return &command.LogsCommand{
				Meta: *meta,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
