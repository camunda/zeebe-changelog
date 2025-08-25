package gitlog

import (
	"log"
	"os/exec"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestGitHistory(t *testing.T) {

	// clone zeebe repo to test with merge commits
	command := exec.Command("git", "clone", "-b", "8.5.0", "https://github.com/camunda/zeebe.git", "zeebe")
	log.Println(command)

	out, _ := command.CombinedOutput()
	log.Println(out)

	// use git command til git lib implements range feature, see https://github.com/src-d/go-git/issues/1166
	tests := map[string]struct {
		path  string
		start string
		end   string
		size  int
	}{
		"First commit in zcl":        {path: ".", start: "7b86247", end: "7ab8381", size: 0},
		"Between tags in zeebe repo": {path: "zeebe", start: "8.5.0", end: "8.6.0-alpha1", size: 1558638},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			log := GetHistory(tc.path, tc.start, tc.end)
			assert.Equal(t, tc.size, len(log))
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
		"URL reference without www"      {message: "closes https://github.com/camunda/camunda/1234", issueIds: []int{1234}},
		"Short reference"                {message: "closes camunda/camunda#1234", issueIds: []int{1234}},
		"Short old reference"            {message: "closes camunda/zeebe#1234", issueIds: []int{1234}},
		"Wrong repo reference URL"       {message: "closes https://www.github.com/camunda/operate/1234", issueIds: nil},
		"Wrong repo short reference"     {message: "closes camunda/operate#1234", issueIds: nil},
		"Wrong org URL reference"        {message: "closes https://www.github.com/zeebe-io/zeebe#1234", issueIds: nil},
		"Wrong org short reference"      {message: "closes zeebe-io/zeebe#1234", issueIds: nil},
		"Multiple issues mixed format":  {message: "closes #1234, camunda/zeebe#5678, camunda/camunda#9 and https://www.github.com/camunda/camunda/123", issueIds: []int{1234, 5678, 9, 123}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			issueIds := ExtractIssueIds(tc.message)
			assert.Equal(t, tc.issueIds, issueIds)
		})
	}
}
