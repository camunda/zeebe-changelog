package github

import (
	"context"
	"github.com/google/go-github/v83/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

const (
	defaultLabelColor           = "8e8e8e"
	labelVerificationRetryDelay = 5 * time.Second
)

type Client struct {
	ctx    context.Context
	client *github.Client
	sleep  func(time.Duration)
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
		sleep:  time.Sleep,
	}
}

func (ghc *Client) EnsureLabelExists(githubOrg, githubRepo, label string) {
	exists, err := ghc.LabelExists(githubOrg, githubRepo, label)
	if err != nil {
		log.Fatalln(err)
	}

	if exists {
		return
	}

	log.Printf("Label %q was not found in %s/%s. Creating it...\n", label, githubOrg, githubRepo)
	_, _, _ = ghc.client.Issues.CreateLabel(ghc.ctx, githubOrg, githubRepo, &github.Label{
		Name:  github.Ptr(label),
		Color: github.Ptr(defaultLabelColor),
	})

	ghc.sleep(labelVerificationRetryDelay)

	exists, err = ghc.LabelExists(githubOrg, githubRepo, label)
	if err != nil {
		log.Fatalln(err)
	}
	if !exists {
		log.Fatalf("Unable to verify creation of label %q in %s/%s\n", label, githubOrg, githubRepo)
	}
}

func (ghc *Client) LabelExists(githubOrg, githubRepo, label string) (bool, error) {
	_, _, err := ghc.client.Issues.GetLabel(ghc.ctx, githubOrg, githubRepo, label)
	if err == nil {
		return true, nil
	}

	if errResp, ok := err.(*github.ErrorResponse); ok && errResp.Response != nil && errResp.Response.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, err
}

func (ghc *Client) AddLabel(githubOrg string, githubRepo string, issueId int, label string) {
	_, _, err := ghc.client.Issues.AddLabelsToIssue(ghc.ctx, githubOrg, githubRepo, issueId, []string{label})
	if err != nil {
		// Check if the error is a 404 Not Found error
		if errResp, ok := err.(*github.ErrorResponse); ok && errResp.Response.StatusCode == http.StatusNotFound {
			log.Printf("Warning: Issue #%d not found in %s/%s, skipping label addition\n", issueId, githubOrg, githubRepo)
			return
		}
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

		options.ListOptions.Page = response.NextPage
	}

	return changelog
}
