package github

import (
	"time"
)

type RateLimit struct {
	Limit     int
	Cost      int
	Remaining int
	ResetAt   time.Time
}

type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type FText struct {
	FValue struct {
		Text string
	} `graphql:"... on ProjectV2ItemFieldTextValue"`
}

func (f FText) V() string {
	return f.FValue.Text
}

type FNumber struct {
	FValue struct {
		Number float64
	} `graphql:"... on ProjectV2ItemFieldNumberValue"`
}

func (f FNumber) V() float64 {
	return f.FValue.Number
}

type FDate struct {
	FValue struct {
		Date time.Time
	} `graphql:"... on ProjectV2ItemFieldDateValue"`
}

func (f FDate) V() time.Time {
	return f.FValue.Date
}

type FIteration struct {
	FValue struct {
		Title     string
		StartDate string
		Duration  int
	} `graphql:"... on ProjectV2ItemFieldIterationValue"`
}

func (f FIteration) V() string {
	return f.FValue.Title
}

type FSingleSelect struct {
	FValue struct {
		Name string
	} `graphql:"... on ProjectV2ItemFieldSingleSelectValue"`
}

func (f FSingleSelect) V() string {
	return f.FValue.Name
}
