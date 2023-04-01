package main

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.githedgehog.com/gh-tools/pkg/github"
	"go.githedgehog.com/gh-tools/pkg/mdbook"
)

func main() {
	app := &cli.App{
		Name:    "mdbook-gh-project",
		Version: "0.0.0", // TODO load proper version using ld flags
		Action: func(ctx *cli.Context) error {
			gh := github.New()

			items, err := gh.GetProjectItems(context.Background(), "PVT_kwDOBvRah84AMA6Y")
			if err != nil {
				return errors.Wrap(err, "error getting github project items")
			}

			return errors.Wrap(mdbook.NewGitHubProjectPreprocessor("githedgehog", items).Process(os.Stdin, os.Stdout),
				"error processing book")
		},
		Commands: []*cli.Command{
			{
				Name: "supports",
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
