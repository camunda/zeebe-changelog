package github

import (
	"github.com/google/go-github/v83/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIssue_IsPullRequest(t *testing.T) {
	tests := map[string]struct {
		issue         *Issue
		isPullRequest bool
	}{
		"Issue":        {issue: createIssue("", 0, "", false), isPullRequest: false},
		"Pull Request": {issue: createIssue("", 0, "", true), isPullRequest: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.isPullRequest, tc.issue.IsPullRequest())
		})
	}
}

func TestIssue_HasBrokerLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, brokerLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, featureLabel, brokerLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasBrokerLabel())
		})
	}
}

func TestIssue_HasGatewayLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, gatewayLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, featureLabel, gatewayLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasGatewayLabel())
		})
	}
}

func TestIssue_HasJavaClientLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, brokerLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, featureLabel, brokerLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasJavaClientLabel())
		})
	}
}

func TestIssue_HasGoClientLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, goClientLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, goClientLabel, featureLabel, brokerLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasGoClientLabel())
		})
	}
}

func TestIssue_HasEnhancementLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, featureLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, featureLabel, brokerLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasEnhancementLabel())
		})
	}
}

func TestIssue_HasBugLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, bugLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, bugLabel, brokerLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasBugLabel())
		})
	}
}

func TestIssue_HasDocsLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, docsLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, brokerLabel, docsLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasDocsLabel())
		})
	}
}

func TestIssue_HasToilLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, toilLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, brokerLabel, toilLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasToilLabel())
		})
	}
}

func TestIssue_HasTaskLabel(t *testing.T) {
	tests := map[string]struct {
		issue    *Issue
		hasLabel bool
	}{
		"No Label":        {issue: createIssue("", 0, "", false), hasLabel: false},
		"Different Label": {issue: createIssue("", 0, "", false, javaClientLabel), hasLabel: false},
		"Has Label":       {issue: createIssue("", 0, "", false, taskLabel), hasLabel: true},
		"Multiple Labels": {issue: createIssue("", 0, "", false, javaClientLabel, brokerLabel, taskLabel), hasLabel: true},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.hasLabel, tc.issue.HasTaskLabel())
		})
	}
}

func TestIssue_String(t *testing.T) {
	issue := createIssue("Test Issue", 1234, "https://github.com", false)
	assert.Equal(t, "Test Issue ([#1234](https://github.com))", issue.String())
}

func createIssue(title string, number int, url string, pullRequest bool, labels ...string) *Issue {
	var pullRequestLinks *github.PullRequestLinks
	var labelList []*github.Label

	if labels != nil {
		labelList = []*github.Label{}
		for _, label := range labels {
			labelCopy := label
			labelList = append(labelList, &github.Label{Name: &labelCopy})
		}
	}

	if pullRequest {
		pullRequestLinks = &github.PullRequestLinks{}
	}

	return NewIssue(&github.Issue{
		Title:            &title,
		Number:           &number,
		HTMLURL:          &url,
		Labels:           labelList,
		PullRequestLinks: pullRequestLinks,
	})
}
