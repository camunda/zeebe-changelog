package main

import (
	"fmt"
	"log"
	"os"

	"github.com/camunda/zeebe-changelog/pkg/github"
	"github.com/camunda/zeebe-changelog/pkg/gitlog"
	"github.com/camunda/zeebe-changelog/pkg/progress"
	"github.com/urfave/cli"
)

const (
	appName           = "zcl"
	gitApiTokenFlag   = "token"
	gitApiTokenEnv    = "GITHUB_TOKEN"
	gitDirFlag        = "gitDir"
	gitDirEnv         = "ZCL_GIT_DIR"
	labelFlag         = "label"
	labelEnv          = "ZCL_LABEL"
	fromFlag          = "from"
	fromEnv           = "ZCL_FROM_REV"
	targetFlag        = "target"
	targetEnv         = "ZCL_TARGET_REV"
	githubOrgFlag     = "org"
	githubOrgEnv      = "ZCL_ORG"
	githubOrgDefault  = "zeebe-io"
	githubRepoFlag    = "repo"
	githubRepoEnv     = "ZCL_REPO"
	githubRepoDefault = "zeebe"
)

var (
	version = "development"
	commit  = "HEAD"
)

func main() {
	app := createApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = appName
	app.HelpName = appName
	app.Usage = "Zeebe Changelog Helper"
	app.Version = fmt.Sprintf("%s (commit: %s)", version, commit)

	app.Commands = []cli.Command{
		{
			Name:      "add-labels",
			ShortName: "a",
			Usage:     "Add GitHub labels to issues and PRs",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     gitDirFlag,
					Usage:    "Git working directory",
					EnvVar:   gitDirEnv,
					Required: false,
					Value:    ".",
				},
				cli.StringFlag{
					Name:     labelFlag,
					EnvVar:   labelEnv,
					Usage:    "GitHub label to attach to issues and PRs",
					Required: true,
				},
				cli.StringFlag{
					Name:     fromFlag,
					EnvVar:   fromEnv,
					Usage:    "Git revision to start start processing",
					Required: true,
				},
				cli.StringFlag{
					Name:     targetFlag,
					EnvVar:   targetEnv,
					Usage:    "Git revision to stop commit processing",
					Required: true,
				},
				cli.StringFlag{
					Name:     gitApiTokenFlag,
					Usage:    "GitHub API Token",
					EnvVar:   gitApiTokenEnv,
					Required: true,
				},
				cli.StringFlag{
					Name:   githubOrgFlag,
					Usage:  "GitHub organization",
					EnvVar: githubOrgEnv,
					Value:  githubOrgDefault,
				},
				cli.StringFlag{
					Name:   githubRepoFlag,
					Usage:  "GitHub repository",
					EnvVar: githubRepoEnv,
					Value:  githubRepoDefault,
				},
			},
			Action: addLabels,
		},
		{
			Name:      "generate",
			ShortName: "g",
			Usage:     "Generate change log",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     labelFlag,
					EnvVar:   labelEnv,
					Usage:    "GitHub label name to generate changelog from",
					Required: true,
				},
				cli.StringFlag{
					Name:     gitApiTokenFlag,
					Usage:    "GitHub API Token",
					EnvVar:   gitApiTokenEnv,
					Required: true,
				},
				cli.StringFlag{
					Name:   githubOrgFlag,
					Usage:  "GitHub organization",
					EnvVar: githubOrgEnv,
					Value:  githubOrgDefault,
				},
				cli.StringFlag{
					Name:   githubRepoFlag,
					Usage:  "GitHub repository",
					EnvVar: githubRepoEnv,
					Value:  githubRepoDefault,
				},
			},
			Action: generateChangelog,
		},
	}

	return app
}

func addLabels(c *cli.Context) error {
	token := c.String(gitApiTokenFlag)
	gitDir := c.String(gitDirFlag)
	from := c.String(fromFlag)
	target := c.String(targetFlag)
	githubOrg := c.String(githubOrgFlag)
	githubRepo := c.String(githubRepoFlag)
	label := c.String(labelFlag)

	log.Println("Fetching git history in dir", gitDir, "for", from, "..", target)

	commits := gitlog.GetHistory(gitDir, from, target)

	log.Println("Collection issue ids")
	issueIds := gitlog.ExtractIssueIds(commits)

	issueCount := len(issueIds)
	log.Println("Updating", issueCount, "issues")

	client := github.NewClient(token)
	bar := progress.NewProgressBar(issueCount)

	for _, id := range issueIds {
		client.AddLabel(githubOrg, githubRepo, id, label)
		bar.Increase()
	}

	return nil
}

func generateChangelog(c *cli.Context) error {
	token := c.String(gitApiTokenFlag)
	githubOrg := c.String(githubOrgFlag)
	githubRepo := c.String(githubRepoFlag)
	label := c.String(labelFlag)

	client := github.NewClient(token)

	log.Println("Fetching issues for GitHub label", label)
	changelog := client.FetchIssues(githubOrg, githubRepo, label)

	log.Println("Generating changelog for GitHub label", label)
	fmt.Println(changelog.String())

	return nil
}
