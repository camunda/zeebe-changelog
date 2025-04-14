package github

import (
	"context"
	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
	"log"
)

type Client struct {
	ctx    context.Context
	client *github.Client
}

func NewClient(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Client{
		ctx:    ctx,
		client: client,
	}
}

func (ghc *Client) AddLabel(githubOrg string, githubRepo string, issueId int, label string) {
	_, _, err := ghc.client.Issues.AddLabelsToIssue(ghc.ctx, githubOrg, githubRepo, issueId, []string{label})
	if err != nil {
		log.Fatalln(err)
	}
}

func (ghc *Client) FetchIssues(githubOrg, githubRepo, label string) *Changelog {
	options := &github.IssueListByRepoOptions{State: "all", Labels: []string{label}}
	changelog := NewChangelog(label)

	for {
		issues, response, err := ghc.client.Issues.ListByRepo(ghc.ctx, githubOrg, githubRepo, options)
		if err != nil {
			log.Fatalln(err)
		}

		for _, issue := range issues {
			changelog.AddIssue(NewIssue(issue))
		}

		if response.NextPage == 0 {
			break
		}

		options.Page = response.NextPage
	}

	return changelog
}
