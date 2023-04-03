package mdbook

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
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
				log.Println("Processing github item:", item)

				formatted, err := p.formatGitHubItem(item)
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

func (p *GitHubProjectPreprocessor) formatGitHubItem(ref string) (string, error) {
	parts := strings.Split(ref, "#")
	if len(parts) != 2 {
		return "", errors.Errorf("invalid github item ref, should be repo#number")
	}
	repo := parts[0]
	number := parts[1]

	if item, ok := p.items[fmt.Sprintf("/%s/%s/issues/%s", p.org, repo, number)]; ok {
		return fmt.Sprintf(" <a href='%s' title='%s'>`%s#%s: %s (%.1f)`</a> ",
			item.URL(), html.EscapeString(item.Title()), repo, number, item.Status(), item.Estimate()), nil
	}

	return fmt.Sprintf(" <a href=\"https://github.com/%s/issues/%s\">`%s#%s`</a> ", repo, number, repo, number), nil
}
