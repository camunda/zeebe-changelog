package github

import (
	"fmt"
	"github.com/google/go-github/v83/github"
)

const (
	brokerLabel     = "scope/broker"
	gatewayLabel    = "scope/gateway"
	javaClientLabel = "scope/clients-java"
	goClientLabel   = "scope/clients-go"
	zbctlLabel      = "scope/zbctl"
	featureLabel    = "kind/feature"
	bugLabel        = "kind/bug"
	docsLabel       = "kind/documentation"
	toilLabel       = "kind/toil"
	taskLabel       = "kind/task"
)

var knownLabels = []string{
	brokerLabel,
	gatewayLabel,
	javaClientLabel,
	goClientLabel,
	featureLabel,
	bugLabel,
	docsLabel,
	zbctlLabel,
	toilLabel,
	taskLabel,
}

type Issue struct {
	title       *string
	number      *int
	url         *string
	labels      map[string]bool
	pullRequest bool
}

func NewIssue(issue *github.Issue) *Issue {
	return &Issue{
		title:       issue.Title,
		number:      issue.Number,
		url:         issue.HTMLURL,
		labels:      mapLabels(issue.Labels),
		pullRequest: issue.IsPullRequest(),
	}
}

func mapLabels(labelList []*github.Label) map[string]bool {
	labels := make(map[string]bool)
	for _, label := range labelList {
		labelName := label.GetName()
		for _, knownLabel := range knownLabels {
			if labelName == knownLabel {
				labels[knownLabel] = true
			}
		}
	}
	return labels
}

func (i *Issue) HasBrokerLabel() bool {
	return i.hasLabel(brokerLabel)
}

func (i *Issue) HasGatewayLabel() bool {
	return i.hasLabel(gatewayLabel)
}

func (i *Issue) HasJavaClientLabel() bool {
	return i.hasLabel(javaClientLabel)
}

func (i *Issue) HasGoClientLabel() bool {
	return i.hasLabel(goClientLabel)
}

func (i *Issue) HasEnhancementLabel() bool {
	return i.hasLabel(featureLabel)
}

func (i *Issue) HasBugLabel() bool {
	return i.hasLabel(bugLabel)
}

func (i *Issue) HasDocsLabel() bool {
	return i.hasLabel(docsLabel)
}

func (i *Issue) HasZbctlLabel() bool {
	return i.hasLabel(zbctlLabel)
}

func (i *Issue) HasToilLabel() bool {
	return i.hasLabel(toilLabel)
}

func (i *Issue) HasTaskLabel() bool {
	return i.hasLabel(taskLabel)
}

func (i *Issue) hasLabel(label string) bool {
	return i.labels[label]
}

func (i *Issue) IsPullRequest() bool {
	return i.pullRequest
}

func (i *Issue) String() string {
	return fmt.Sprintf("%s ([#%d](%s))", *i.title, *i.number, *i.url)
}
