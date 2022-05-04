package command

import "github.com/urfave/cli/v2"

// Application initializes the client application and its description
var Application = cli.App{
	Name:     "circleci-trigger",
	Usage:    "Triggers pipelines in a remote circleci project",
	HelpName: "circleci-trigger",
	Commands: []*cli.Command{
		&runCommand,
	},
	EnableBashCompletion: true,
	HideHelp:             false,
}
