package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"go.githedgehog.com/gh-tools/pkg/github"
)

func main() {
	app := &cli.App{
		Name:    "gh-project-dump",
		Version: "0.0.0", // TODO load proper version using ld flags
		Action: func(ctx *cli.Context) error {
			gh := github.New()

			items, err := gh.GetProjectItems(context.Background(), "PVT_kwDOBvRah84AMA6Y")
			if err != nil {
				log.Printf("Error getting github items: %#v\n", err)
				return err
			}
			log.Println("Loaded project items:", len(items))

			for _, item := range items {
				fmt.Println(">>>", item.Milestone(), "<<<", item.Number(), item.Title(), item.Status(), item.Sprint(), item.Assignees(), item.Progress(), item.Component(), item.Estimate(), item.Jira(), item.Resource())
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "test",
				Action: func(ctx *cli.Context) error {
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err) // TODO use log
	}
}
