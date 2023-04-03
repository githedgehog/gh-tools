package mdbook

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.githedgehog.com/gh-tools/pkg/github"
)

const (
	TOKEN_OPEN  = "{{ $gh "
	TOKEN_CLOSE = "}}"
)

type GitHubProjectPreprocessor struct {
	org   string
	items map[string]github.ProjectItem
}

func NewGitHubProjectPreprocessor(org string, items map[string]github.ProjectItem) *GitHubProjectPreprocessor {
	return &GitHubProjectPreprocessor{
		org:   org,
		items: items,
	}
}

func (p *GitHubProjectPreprocessor) Process(in io.Reader, out io.Writer) error {
	source, err := p.readBook(in)
	if err != nil {
		return errors.Wrap(err, "error reading book")
	}

	result := strings.Builder{}
	for {
		start := strings.Index(source, TOKEN_OPEN)
		if start >= 0 {
			_, err = result.WriteString(source[:start])
			if err != nil {
				return errors.Wrap(err, "error writing pre-start string")
			}

			start += len(TOKEN_OPEN)
			end := strings.Index(source[start:], TOKEN_CLOSE) + start
			if end-start > 3 && end < len(source) {
				item := source[start : end-1]

				formatted, err := p.processExpr(item)
				if err != nil {
					return errors.Wrapf(err, "error formatting github item %s", item)
				}

				_, err = result.WriteString(formatted)
				if err != nil {
					return errors.Wrap(err, "error writing processed item")
				}

				source = source[end+len(TOKEN_CLOSE):]
			} else {
				return errors.Errorf("failed to find closing token for github processing, should be %s repo#number %s",
					TOKEN_OPEN, TOKEN_CLOSE)
			}
		} else {
			_, err = result.WriteString(source)
			if err != nil {
				return errors.Wrap(err, "error writing string w/o start")
			}
			break
		}
	}

	_, err = fmt.Fprintln(out, result.String())

	return err
}

func (p *GitHubProjectPreprocessor) readBook(in io.Reader) (string, error) {
	input := []map[string]any{}

	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		return "", errors.Wrap(err, "error parsing context and book")
	}

	bookBytes, err := json.Marshal(input[1])
	if err != nil {
		return "", errors.Wrap(err, "error marshalling book")
	}

	return string(bookBytes), nil
}

func (p *GitHubProjectPreprocessor) processExpr(ref string) (string, error) {
	parts := strings.Split(ref, "#")
	if len(parts) == 2 {
		if parts[0] == "$summary" {
			return p.formatSummary(parts[1])
		} else {
			return p.formatGitHubItem(parts[0], parts[1])
		}
	}
	return "", errors.Errorf("invalid expression: should be '$summary#sprint' or github item ref like 'repo#number'")
}

func (p *GitHubProjectPreprocessor) formatGitHubItem(repo string, number string) (string, error) {
	if item, ok := p.items[fmt.Sprintf("/%s/%s/issues/%s", p.org, repo, number)]; ok {
		return fmt.Sprintf(" <a href='%s' title='%s'>`%s#%s: %s (%.1f)`</a> ",
			item.URL(), html.EscapeString(item.Title()), repo, number, item.Status(), item.Estimate()), nil
	}

	return fmt.Sprintf(" <a href='https://github.com/githedgehog/%s/issues/%s'>`%s#%s`</a> ",
		repo, number, repo, number), nil
}

var usersMapping = map[string]string{
	"Frostman":    "Sergei L",
	"mheese":      "Marcus",
	"mkoperator":  "Mikhail",
	"sonoble":     "Steve",
	"amitlimaye":  "Amit",
	"sergeymatov": "Sergey M",
}

func (p *GitHubProjectPreprocessor) formatAssignee(name string) string {
	if name == "" {
		return "unassigned"
	}

	if mapped, ok := usersMapping[name]; ok {
		return mapped
	}

	return name
}

var statuses = []string{
	"🆕 New",
	"📋 Backlog",
	"🏗 In progress",
	"👀 In review",
	"✅ Done",
}

func (p *GitHubProjectPreprocessor) formatSummary(sprint string) (string, error) {
	stats := map[string][]github.ProjectItem{}

	for _, item := range p.items {
		if item.Sprint() != sprint {
			continue
		}

		name := p.formatAssignee(item.Assignee())
		stats[name] = append(stats[name], item)
	}

	summary := strings.Builder{}
	summary.WriteString("<table>")

	summary.WriteString("<thead>")
	summary.WriteString("<tr>")
	summary.WriteString("<th>Name</th>")
	for _, status := range statuses {
		summary.WriteString(fmt.Sprintf("<th>%s</th>", status))
	}
	summary.WriteString("<th>Total</th>")
	summary.WriteString("</tr>")
	summary.WriteString("</thead>")

	summary.WriteString("<tbody>")
	for name, items := range stats {
		summary.WriteString("<tr>")
		summary.WriteString(fmt.Sprintf("<td>%s</td>", name))

		total := 0.0
		for _, status := range statuses {
			estimate := 0.0
			for _, item := range items {
				if item.Status() != status {
					continue
				}

				estimate += item.Estimate()
			}
			total += estimate
			summary.WriteString(fmt.Sprintf("<td>%.1f</td>", estimate))
		}

		summary.WriteString(fmt.Sprintf("<td>%.1f</td>", total))
		summary.WriteString("</tr>")
	}

	summary.WriteString("</tbody>")
	summary.WriteString("</table>")

	return summary.String(), nil
}
