package main

import (
	"context"
	"fmt"

	"go.githedgehog.com/gh-tools/pkg/github"
)

func main() {
	gh := github.New()

	items, err := gh.GetProjectItems(context.Background(), "PVT_kwDOBvRah84AMA6Y")
	if err != nil {
		panic(err)
	}

	fmt.Println(items)
	fmt.Println(len(items))
}
