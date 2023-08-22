package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"go.githedgehog.com/gh-tools/pkg/gdrive"
	"go.githedgehog.com/gh-tools/pkg/github"
)

func main() {
	var name, sprint string

	app := &cli.App{
		Name:    "gh-project-dump",
		Version: "0.0.0", // TODO load proper version using ld flags
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Value:       "0w0-test",
				Usage:       "update name",
				Destination: &name,
			},
			&cli.StringFlag{
				Name:        "sprint",
				Aliases:     []string{"s"},
				Value:       "0w0",
				Usage:       "sprint name",
				Destination: &sprint,
			},
		},
		Action: func(ctx *cli.Context) error {
			log.Println("Sprint:", sprint)
			log.Println("Name:", name)

			githubProjectID := "PVT_kwDOBvRah84AMA6Y"

			sheetID := 0
			spreadsheetID := "1hP2dy9JXPgTqMx3-x66cykojhp9JHbweBumfdwypQ3Y"

			gh := github.New()

			items, err := gh.GetProjectItems(context.Background(), githubProjectID)
			if err != nil {
				log.Printf("Error getting github items: %#v\n", err)
				return err
			}
			log.Println("Loaded project items:", len(items))

			updated := time.Now().Format(time.DateTime)

			records := [][]string{}
			for _, item := range items {
				if item.Sprint() != sprint {
					continue
				}

				resource := strings.ReplaceAll(item.Resource()[13:], "/issues/", "#")
				records = append(records, []string{
					name,
					updated,
					fmt.Sprintf("=HYPERLINK(\"%s\",\"%s\")", item.URL(), resource),
					item.Title(),
					item.Status(),
					item.Assignee(),
					fmt.Sprintf("%.2f", item.Estimate()),
					item.Progress(),
					fmt.Sprintf("%.2f", item.ProgressPercentage()),
					fmt.Sprintf("%.2f", item.ProgressEstimate()),
					item.Sprint(),
					item.Component(),
					item.Jira(),
					fmt.Sprint(item.Assignees()),
					item.SprintStartDate(),
				})
			}
			log.Println("Records in current sprint:", len(records))

			err = gdrive.New().AppendToSpreadsheet(context.Background(), spreadsheetID, sheetID, records)
			if err != nil {
				return err
			}
			log.Println("Spreadsheet updated")

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
