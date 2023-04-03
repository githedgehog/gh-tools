package github

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func (gh *GitHub) GetProjectItems(ctx context.Context, project string) (map[string]ProjectItem, error) {
	var q struct {
		RateLimit RateLimit
		Node      struct {
			Project struct {
				Items struct {
					TotalCount int
					PageInfo   PageInfo
					Nodes      []ProjectItem
				} `graphql:"items(first: $limit, after: $cursor)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $project)"`
	}

	vars := map[string]interface{}{
		"project": githubv4.ID(project),
		"limit":   githubv4.Int(100),
		"cursor":  githubv4.String(""),
	}

	items := map[string]ProjectItem{}

	for {
		err := gh.client.Query(ctx, &q, vars)
		if err != nil {
			return nil, err
		}

		// log.Println("RateLimit", q.RateLimit)

		for _, item := range q.Node.Project.Items.Nodes {
			items[item.Resource()] = item
		}

		if !q.Node.Project.Items.PageInfo.HasNextPage {
			break
		}

		vars["cursor"] = githubv4.String(q.Node.Project.Items.PageInfo.EndCursor)
	}

	return items, nil
}

type ProjectItem struct {
	Type       string
	FStatus    FSingleSelect `graphql:"status: fieldValueByName(name: \"Status\")"`
	FEstimate  FNumber       `graphql:"estimate: fieldValueByName(name: \"Estimate\")"`
	FIteration FIteration    `graphql:"iteration: fieldValueByName(name: \"Iteration\")"`
	FJira      FText         `graphql:"jira: fieldValueByName(name: \"Jira\")"`
	Content    struct {
		Issue struct {
			Title        string
			Number       int
			State        string
			URL          string
			ResourcePath string
			Assignees    struct {
				Nodes []struct {
					Login string
				}
			} `graphql:"assignees(first: 10)"`
		} `graphql:"... on Issue"`
	}
}

func (i ProjectItem) Status() string {
	return i.FStatus.V()
}

func (i ProjectItem) Estimate() float64 {
	return i.FEstimate.V()
}

func (i ProjectItem) Sprint() string {
	return i.FIteration.V()
}

func (i ProjectItem) Jira() string {
	return i.FJira.V()
}

func (i ProjectItem) Title() string {
	return i.Content.Issue.Title
}

func (i ProjectItem) Number() int {
	return i.Content.Issue.Number
}

func (i ProjectItem) State() string {
	return i.Content.Issue.State
}

func (i ProjectItem) URL() string {
	return i.Content.Issue.URL
}

func (i ProjectItem) Resource() string {
	return i.Content.Issue.ResourcePath
}

func (i ProjectItem) Assignee() string {
	if len(i.Content.Issue.Assignees.Nodes) == 0 {
		return ""
	}

	return i.Content.Issue.Assignees.Nodes[0].Login
}

func (i ProjectItem) String() string {
	return fmt.Sprintf("Issue{%s %s est=%.1f iter=%s jira=%s}",
		i.Resource(), i.Status(), i.Estimate(), i.Sprint(), i.Jira())
}
