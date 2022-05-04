package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	command "github.com/stacknowledge/circleci-trigger-action/cmd"
)

func main() {
	command.Application.Flags = append(command.Application.Flags, []cli.Flag{}...)

	if err := command.Application.Run(os.Args); err != nil {
		log.Fatalf("Fatal error, exiting application: %s", err)
	}

	os.Exit(0)
}
