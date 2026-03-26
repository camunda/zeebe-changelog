package gitlog

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareDivergedRepo(t *testing.T) (string, string, string, string) {
	t.Helper()
	repoDir := filepath.Join(t.TempDir(), "diverged")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir temp repo: %v", err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "zcl-tests@example.com")
	runGit(t, repoDir, "config", "user.name", "zcl-tests")

	baseFile := filepath.Join(repoDir, "base.txt")
	if err := os.WriteFile(baseFile, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runGit(t, repoDir, "add", "base.txt")
	runGit(t, repoDir, "commit", "-m", "base")
	base := strings.TrimSpace(runGit(t, repoDir, "rev-parse", "HEAD"))

	runGit(t, repoDir, "checkout", "-b", "branch-a")
	aFile := filepath.Join(repoDir, "a.txt")
	if err := os.WriteFile(aFile, []byte("a\n"), 0o644); err != nil {
		t.Fatalf("write branch-a file: %v", err)
	}
	runGit(t, repoDir, "add", "a.txt")
	runGit(t, repoDir, "commit", "-m", "branch-a")
	branchA := strings.TrimSpace(runGit(t, repoDir, "rev-parse", "HEAD"))

	runGit(t, repoDir, "checkout", "-b", "branch-b", base)
	bFile := filepath.Join(repoDir, "b.txt")
	if err := os.WriteFile(bFile, []byte("b\n"), 0o644); err != nil {
		t.Fatalf("write branch-b file: %v", err)
	}
	runGit(t, repoDir, "add", "b.txt")
	runGit(t, repoDir, "commit", "-m", "branch-b")
	branchB := strings.TrimSpace(runGit(t, repoDir, "rev-parse", "HEAD"))

	return repoDir, base, branchA, branchB
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	// avoid hanging on any credential prompts
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	log.Println(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, string(out))
	}
	return string(out)
}

func prepareCamundaRepo(t *testing.T, tagA, tagB string) string {
	t.Helper()
	repoDir := filepath.Join(t.TempDir(), "camunda")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir temp repo: %v", err)
	}

	runGit(t, repoDir, "init")
	runGit(t, repoDir, "remote", "add", "origin", "https://github.com/camunda/camunda.git")

	// Fetch only the tags we care about. Use a partial clone filter if available.
	fetchArgs := []string{"fetch", "--no-tags", "--filter=blob:none", "origin", "tag", tagA, "tag", tagB}
	cmd := exec.Command("git", fetchArgs...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	log.Println(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback for older git versions without partial clone support.
		log.Printf("partial fetch failed, retrying without filter: %v\n%s", err, string(out))
		runGit(t, repoDir, "fetch", "--no-tags", "origin", "tag", tagA, "tag", tagB)
	}

	// Ensure tags are present locally (handles annotated tags)
	runGit(t, repoDir, "fetch", "--tags", "--force", "origin")
	return repoDir
}

func TestGitHistory(t *testing.T) {
	// use git command til git lib implements range feature, see https://github.com/src-d/go-git/issues/1166
	tests := map[string]struct {
		path        string
		start       string
		end         string
		size        int
		needCamunda bool
	}{
		"First commit in zcl":          {path: ".", start: "7b86247", end: "7ab8381", size: 0},
		"Between tags in camunda repo": {start: "8.5.0", end: "8.6.0-alpha1", size: 2012037, needCamunda: true},
	}

	var camundaRepo string
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			path := tc.path
			if tc.needCamunda {
				if camundaRepo == "" {
					camundaRepo = prepareCamundaRepo(t, "8.5.0", "8.6.0-alpha1")
				}
				path = camundaRepo
			}

			log := GetHistory(path, tc.start, tc.end)
			assert.Equal(t, tc.size, len(log))
		})
	}
}

func TestRegressionBackportIssue40036IncludedInReleaseRanges(t *testing.T) {
	tests := map[string]struct {
		from           string
		to             string
		expectedCommit string
	}{
		"8.6.38 range": {from: "8.6.37", to: "8.6.38", expectedCommit: "87bbd207bc0fc7e7dbfbd0ffda36af6cc35c0b50"},
		"8.7.25 range": {from: "8.7.24", to: "8.7.25", expectedCommit: "88bdb736007a93c6abeb1eb8a14093510c915af0"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			camundaRepo := prepareCamundaRepo(t, tc.from, tc.to)
			history := GetHistory(camundaRepo, tc.from, tc.to)
			assert.Contains(t, history, tc.expectedCommit)
			issueIDs := ExtractIssueIds(history)
			assert.Contains(t, issueIDs, 40036)
		})
	}
}

func TestExtractIssueIds(t *testing.T) {
	tests := map[string]struct {
		message  string
		issueIds []int
	}{
		"No issue id":                   {message: "hello world", issueIds: nil},
		"No keyword":                    {message: "#1234", issueIds: nil},
		"Close keyword":                 {message: "close #1234", issueIds: []int{1234}},
		"Closes keyword":                {message: "closes #1234", issueIds: []int{1234}},
		"Related keyword":               {message: "related #1234", issueIds: []int{1234}},
		"Merge keyword":                 {message: "merge #1234", issueIds: []int{1234}},
		"Merges keyword":                {message: "merges #1234", issueIds: []int{1234}},
		"Relates keyword":               {message: "relates to #1234", issueIds: []int{1234}},
		"Backport keyword":              {message: "backport #1234", issueIds: []int{1234}},
		"Backports keyword":             {message: "backports #1234", issueIds: []int{1234}},
		"Back ports keyword":            {message: "back ports #1234", issueIds: []int{1234}},
		"Keyword uppercase":             {message: "Closes #1234", issueIds: []int{1234}},
		"Spacing in front of keyword":   {message: "  \t closes #1234", issueIds: []int{1234}},
		"Multiple issues":               {message: "closes #1234, #5678, #9 and #123", issueIds: []int{1234, 5678, 9, 123}},
		"Duplicate issue ids":           {message: "closes #123, #234, #123 and #23", issueIds: []int{123, 234, 23}},
		"Multiple lines":                {message: "foo bar\n\ncloses #1234\ntest", issueIds: []int{1234}},
		"Multiple IDs without keywords": {message: "foo\n\nbar #234\n\nmerges #1", issueIds: []int{1}},
		"ID with text after":            {message: "closes #4002 drop multi column families usage", issueIds: []int{4002}},
		"Multiple ID with text after":   {message: "closes #5137 low load causes defragmentation\ncloses #4560 unstable cluster on bigger state", issueIds: []int{5137, 4560}},
		"Full URL reference":            {message: "closes https://www.github.com/camunda/camunda/1234", issueIds: []int{1234}},
		"Full old URL reference":        {message: "closes https://www.github.com/camunda/zeebe/1234", issueIds: []int{1234}},
		"URL reference without www":     {message: "closes https://github.com/camunda/camunda/1234", issueIds: []int{1234}},
		"Short reference":               {message: "closes camunda/camunda#1234", issueIds: []int{1234}},
		"Short old reference":           {message: "closes camunda/zeebe#1234", issueIds: []int{1234}},
		"Wrong repo reference URL":      {message: "closes https://www.github.com/camunda/operate/1234", issueIds: nil},
		"Wrong repo short reference":    {message: "closes camunda/operate#1234", issueIds: nil},
		"Wrong org URL reference":       {message: "closes https://www.github.com/zeebe-io/zeebe#1234", issueIds: nil},
		"Wrong org short reference":     {message: "closes zeebe-io/zeebe#1234", issueIds: nil},
		"Multiple issues mixed format":  {message: "closes #1234, camunda/zeebe#5678, camunda/camunda#9 and https://www.github.com/camunda/camunda/123", issueIds: []int{1234, 5678, 9, 123}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			issueIds := ExtractIssueIds(tc.message)
			assert.Equal(t, tc.issueIds, issueIds)
		})
	}
}

func TestValidateAncestor(t *testing.T) {
	repoDir, base, branchA, branchB := prepareDivergedRepo(t)

	t.Run("valid ancestor", func(t *testing.T) {
		err := validateAncestor(repoDir, base, branchA)
		assert.NoError(t, err)
	})

	t.Run("non-ancestor range", func(t *testing.T) {
		err := validateAncestor(repoDir, branchA, branchB)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start is not an ancestor")
	})
}
