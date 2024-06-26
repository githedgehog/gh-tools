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
	FSprint    FIteration    `graphql:"iteration: fieldValueByName(name: \"Sprint\")"`
	FJira      FText         `graphql:"jira: fieldValueByName(name: \"Jira\")"`
	FMilestone FSingleSelect `graphql:"milestone: fieldValueByName(name: \"DevMilestone\")"`
	FProgress  FSingleSelect `graphql:"progress: fieldValueByName(name: \"Progress\")"`
	FComponent FSingleSelect `graphql:"component: fieldValueByName(name: \"Component\")"`
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
	return i.FSprint.V()
}

func (i ProjectItem) SprintStartDate() string {
	return i.FSprint.FValue.StartDate
}

func (i ProjectItem) Jira() string {
	return i.FJira.V()
}

func (i ProjectItem) Milestone() string {
	return i.FMilestone.V()
}

func (i ProjectItem) Progress() string {
	return i.FProgress.V()
}

func (i ProjectItem) ProgressPercentage() float32 {
	status := i.Status()

	if status == "✅ Done" {
		return 1
	}
	if status != "🏗 In progress" {
		return 0
	}

	pr := i.Progress()
	if len(pr) == 0 {
		return 0
	}

	switch pr[0] {
	case '1':
		return 0.1
	case '2':
		return 0.3
	case '4':
		return 0.5
	case '6':
		return 0.7
	case '8':
		return 0.9
	}

	return 0
}

func (i ProjectItem) ProgressEstimate() float64 {
	return i.Estimate() * float64(i.ProgressPercentage())
}

func (i ProjectItem) Component() string {
	return i.FComponent.V()
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

func (i ProjectItem) Assignees() []string {
	assignees := make([]string, len(i.Content.Issue.Assignees.Nodes))

	for i, node := range i.Content.Issue.Assignees.Nodes {
		assignees[i] = node.Login
	}

	return assignees
}

func (i ProjectItem) String() string {
	return fmt.Sprintf("Issue{%s %s est=%.1f iter=%s jira=%s}",
		i.Resource(), i.Status(), i.Estimate(), i.Sprint(), i.Jira())
}
