package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/camunda/zeebe-changelog/pkg/github"
	"github.com/camunda/zeebe-changelog/pkg/gitlog"
	"github.com/camunda/zeebe-changelog/pkg/progress"
	"github.com/urfave/cli/v3"
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
	workersFlag       = "workers"
	workersEnv        = "ZCL_WORKERS"
	workersDefault    = 10
	dryRunFlag        = "dry-run"
	dryRunEnv         = "ZCL_DRY_RUN"
)

var (
	version = "development"
	commit  = "HEAD"
)

func main() {
	app := createApp()
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createApp() *cli.Command {
	return &cli.Command{
		Name:    appName,
		Usage:   "Zeebe Changelog Helper",
		Version: fmt.Sprintf("%s (commit: %s)", version, commit),
		Commands: []*cli.Command{
			{
				Name:    "add-labels",
				Aliases: []string{"a"},
				Usage:   "Add GitHub labels to issues and PRs",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    gitDirFlag,
						Usage:   "Git working directory",
						Sources: cli.EnvVars(gitDirEnv),
						Value:   ".",
					},
					&cli.StringFlag{
						Name:     labelFlag,
						Sources:  cli.EnvVars(labelEnv),
						Usage:    "GitHub label to attach to issues and PRs",
						Required: true,
					},
					&cli.StringFlag{
						Name:     fromFlag,
						Sources:  cli.EnvVars(fromEnv),
						Usage:    "Git revision to start start processing",
						Required: true,
					},
					&cli.StringFlag{
						Name:     targetFlag,
						Sources:  cli.EnvVars(targetEnv),
						Usage:    "Git revision to stop commit processing",
						Required: true,
					},
					&cli.StringFlag{
						Name:     gitApiTokenFlag,
						Usage:    "GitHub API Token",
						Sources:  cli.EnvVars(gitApiTokenEnv),
						Required: true,
					},
					&cli.StringFlag{
						Name:    githubOrgFlag,
						Usage:   "GitHub organization",
						Sources: cli.EnvVars(githubOrgEnv),
						Value:   githubOrgDefault,
					},
					&cli.StringFlag{
						Name:    githubRepoFlag,
						Usage:   "GitHub repository",
						Sources: cli.EnvVars(githubRepoEnv),
						Value:   githubRepoDefault,
					},
					&cli.IntFlag{
						Name:    workersFlag,
						Usage:   "Number of concurrent workers for labeling",
						Sources: cli.EnvVars(workersEnv),
						Value:   workersDefault,
					},
					&cli.BoolFlag{
						Name:    dryRunFlag,
						Usage:   "Print issues that would be labeled without making any changes",
						Sources: cli.EnvVars(dryRunEnv),
					},
				},
				Action: addLabels,
			},
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "Generate change log",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     labelFlag,
						Sources:  cli.EnvVars(labelEnv),
						Usage:    "GitHub label name to generate changelog from",
						Required: true,
					},
					&cli.StringFlag{
						Name:     gitApiTokenFlag,
						Usage:    "GitHub API Token",
						Sources:  cli.EnvVars(gitApiTokenEnv),
						Required: true,
					},
					&cli.StringFlag{
						Name:    githubOrgFlag,
						Usage:   "GitHub organization",
						Sources: cli.EnvVars(githubOrgEnv),
						Value:   githubOrgDefault,
					},
					&cli.StringFlag{
						Name:    githubRepoFlag,
						Usage:   "GitHub repository",
						Sources: cli.EnvVars(githubRepoEnv),
						Value:   githubRepoDefault,
					},
				},
				Action: generateChangelog,
			},
		},
	}
}

func addLabelsParallel(client *github.Client, githubOrg, githubRepo string, issueIds []int, label string, bar *progress.Bar, numWorkers int) {
	// Use a worker pool pattern with reasonable concurrency
	jobs := make(chan int, len(issueIds))
	var wg sync.WaitGroup

	// Start worker goroutines
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for issueId := range jobs {
				client.AddLabel(githubOrg, githubRepo, issueId, label)
				bar.Increase()
			}
		}()
	}

	// Send all issue IDs to workers
	for _, id := range issueIds {
		jobs <- id
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
}

func addLabels(_ context.Context, cmd *cli.Command) error {
	token := cmd.String(gitApiTokenFlag)
	gitDir := cmd.String(gitDirFlag)
	from := cmd.String(fromFlag)
	target := cmd.String(targetFlag)
	githubOrg := cmd.String(githubOrgFlag)
	githubRepo := cmd.String(githubRepoFlag)
	label := cmd.String(labelFlag)
	numWorkers := cmd.Int(workersFlag)
	dryRun := cmd.Bool(dryRunFlag)

	// Validate number of workers
	if numWorkers <= 0 {
		log.Fatalf("Number of workers must be positive, got: %d", numWorkers)
	}

	log.Println("Fetching git history in dir", gitDir, "for", from, "..", target)

	commits := gitlog.GetHistory(gitDir, from, target)

	log.Println("Collection issue ids")
	issueIds := gitlog.ExtractIssueIds(commits)

	issueCount := len(issueIds)

	if dryRun {
		log.Println("[dry-run] Would add label", label, "to", issueCount, "issues in", githubOrg+"/"+githubRepo)
	} else {
		log.Println("Adding label", label, "to", issueCount, "issues in", githubOrg+"/"+githubRepo)
	}
	for _, id := range issueIds {
		fmt.Printf("  https://github.com/%s/%s/issues/%d\n", githubOrg, githubRepo, id)
	}

	client := github.NewClient(token)
	client.EnsureLabelExists(githubOrg, githubRepo, label, dryRun)

	if dryRun {
		return nil
	}

	log.Println("Updating", issueCount, "issues with", numWorkers, "workers")
	bar := progress.NewProgressBar(issueCount)

	addLabelsParallel(client, githubOrg, githubRepo, issueIds, label, bar, numWorkers)

	return nil
}

func generateChangelog(_ context.Context, cmd *cli.Command) error {
	token := cmd.String(gitApiTokenFlag)
	githubOrg := cmd.String(githubOrgFlag)
	githubRepo := cmd.String(githubRepoFlag)
	label := cmd.String(labelFlag)

	client := github.NewClient(token)

	log.Println("Fetching issues for GitHub label", label)
	changelog := client.FetchIssues(githubOrg, githubRepo, label)

	log.Println("Generating changelog for GitHub label", label)
	fmt.Println(changelog.String())

	return nil
}
