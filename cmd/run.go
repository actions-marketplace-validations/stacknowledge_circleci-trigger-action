package command

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/stacknowledge/circleci-trigger-action/pkg/circleci"
	"github.com/urfave/cli/v2"
)

var runCommand = cli.Command{
	Name:   "run",
	Usage:  "run",
	Action: runHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "Identifier for trigger. eg: payments-pipeline [default: uuid]", Required: false},
		&cli.StringFlag{Name: "project", Usage: "Project name to trigger the pipeline. eg: company/project"},
		&cli.StringFlag{Name: "branch", Usage: "Desired branch of the project to run the pipeline [default: master]", Required: false},
		&cli.StringFlag{Name: "token", Usage: "Circleci personal token (machine recommended)"},
		&cli.IntFlag{Name: "timeout", Usage: "Timeout in minutes for the triggered pipeline to run [default: 10]", Required: false},
	},
}

func runHandler(cliContext *cli.Context) error {
	project, branch, token, id, timeout, err := parseInput(cliContext)
	if err != nil {
		log.Fatal(err)
	}

	api := circleci.NewCircleCIAPI(token)
	log.Printf("triggering pipeline on %s project @ %s branch...", project, branch)

	pipelineID, err := api.TriggerPipeline(project, branch, id)
	if err != nil {
		log.Fatal(err)
	}

	timeoutTime := time.After(time.Duration(timeout) * time.Minute)
	ticker := time.NewTicker(2 * time.Second)
	log.Printf("pipeline triggered with id: %s", pipelineID)

	for {
		select {
		case <-timeoutTime:
			log.Fatal("operation timed out")
		case <-ticker.C:
			log.Print("fetching pipeline status...")
			pipeline, err := api.GetPipelineStatus(pipelineID)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("pipeline status: %s", pipeline.Status)
			switch pipeline.Status {
			case circleci.StatusCanceled:
				log.Fatal("pipeline was canceled on circleci")
			case circleci.StatusFailed:
				log.Fatal("pipeline failed, check error on pipeline logs")
			case circleci.StatusSuccess:
				log.Print("pipeline ran successfully")
				return nil
			}
		}
	}
}

func parseInput(context *cli.Context) (string, string, string, string, int, error) {
	var id, project, branch, token string
	var timeout int

	id = uuid.New().String()
	branch = "master"
	timeout = 10

	if !context.IsSet("project") {
		return "", "", "", "", 0, fmt.Errorf("project not specified")
	}

	project = context.String("project")

	if !context.IsSet("token") {
		return "", "", "", "", 0, fmt.Errorf("circleci token not specified")
	}
	token = context.String("token")

	if context.IsSet("branch") {
		branch = context.String("branch")
	}

	if context.IsSet("timeout") {
		timeout = context.Int("timeout")
	}

	if context.IsSet("id") {
		id = context.String("id")
	}

	return project, branch, token, id, timeout, nil
}
