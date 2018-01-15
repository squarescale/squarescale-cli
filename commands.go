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
		"cluster": func() (cli.Command, error) {
			return &command.ClusterCommand{}, nil
		},
		"cluster set": func() (cli.Command, error) {
			return &command.ClusterSetCommand{
				Meta: *meta,
			}, nil
		},
		"cluster list": func() (cli.Command, error) {
			return &command.ClusterListCommand{
				Meta: *meta,
			}, nil
		},
		"container": func() (cli.Command, error) {
			return &command.ContainerCommand{}, nil
		},
		"container list": func() (cli.Command, error) {
			return &command.ContainerListCommand{
				Meta: *meta,
			}, nil
		},
		"container show": func() (cli.Command, error) {
			return &command.ContainerShowCommand{
				Meta: *meta,
			}, nil
		},
		"container set": func() (cli.Command, error) {
			return &command.ContainerSetCommand{
				Meta: *meta,
			}, nil
		},
		"container delete": func() (cli.Command, error) {
			return &command.ContainerDeleteCommand{
				Meta: *meta,
			}, nil
		},
		"env": func() (cli.Command, error) {
			return &command.EnvCommand{}, nil
		},
		"env get": func() (cli.Command, error) {
			return &command.EnvGetCommand{
				Meta: *meta,
			}, nil
		},
		"env set": func() (cli.Command, error) {
			return &command.EnvSetCommand{
				Meta: *meta,
			}, nil
		},
		"image": func() (cli.Command, error) {
			return &command.ImageCommand{}, nil
		},
		"image add": func() (cli.Command, error) {
			return &command.ImageAddCommand{
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
		"project remove": func() (cli.Command, error) {
			return &command.ProjectRemoveCommand{
				Meta: *meta,
			}, nil
		},
		"project slackbot": func() (cli.Command, error) {
			return &command.ProjectSlackbotCommand{
				Meta: *meta,
			}, nil
		},
		"task show": func() (cli.Command, error) {
			return &command.TaskShowCommand{
				Meta: *meta,
			}, nil
		},
		"task list": func() (cli.Command, error) {
			return &command.TaskListCommand{
				Meta: *meta,
			}, nil
		},
		"task wait": func() (cli.Command, error) {
			return &command.TaskWaitCommand{
				Meta: *meta,
			}, nil
		},
		"db": func() (cli.Command, error) {
			return &command.DBCommand{}, nil
		},
		"db list": func() (cli.Command, error) {
			return &command.DBListCommand{
				Meta: *meta,
			}, nil
		},
		"db set": func() (cli.Command, error) {
			return &command.DBSetCommand{
				Meta: *meta,
			}, nil
		},
		"db show": func() (cli.Command, error) {
			return &command.DBShowCommand{
				Meta: *meta,
			}, nil
		},
		"lb": func() (cli.Command, error) {
			return &command.LBCommand{}, nil
		},
		"lb list": func() (cli.Command, error) {
			return &command.LBListCommand{
				Meta: *meta,
			}, nil
		},
		"lb set": func() (cli.Command, error) {
			return &command.LBSetCommand{
				Meta: *meta,
			}, nil
		},
		"lb url": func() (cli.Command, error) {
			return &command.LBURLCommand{
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
