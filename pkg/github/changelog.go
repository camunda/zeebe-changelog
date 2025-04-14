package github

import (
	"bytes"
	"fmt"
	"log"
)

type Changelog struct {
	title        string
	enhancements *Section
	fixes        *Section
	docs         []*Issue
	toil         []*Issue
	pullRequests []*Issue
}

func NewChangelog(title string) *Changelog {
	return &Changelog{
		title:        title,
		enhancements: NewSection(),
		fixes:        NewSection(),
		docs:         []*Issue{},
		toil:         []*Issue{},
		pullRequests: []*Issue{},
	}
}

func (c *Changelog) AddIssue(issue *Issue) *Changelog {
	if issue.IsPullRequest() {
		c.pullRequests = append(c.pullRequests, issue)
	} else {
		if issue.HasEnhancementLabel() {
			c.enhancements.AddIssue(issue)
		}
		if issue.HasBugLabel() {
			c.fixes.AddIssue(issue)
		}
		if issue.HasDocsLabel() {
			c.docs = append(c.docs, issue)
		}
		if issue.HasToilLabel() {
			c.toil = append(c.toil, issue)
		}
	}
	return c
}

func (c *Changelog) String() string {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("# %s\n", c.title))

	chapterToString(&b, "Enhancements", c.enhancements)
	chapterToString(&b, "Bug Fixes", c.fixes)

	issueListToString(&b, "Maintenance", c.toil)
	issueListToString(&b, "Documentation", c.docs)
	issueListToString(&b, "Merged Pull Requests", c.pullRequests)

	return b.String()
}

func issueListToString(b *bytes.Buffer, title string, issues []*Issue) {
	if len(issues) > 0 {
		b.WriteString(fmt.Sprintf("## %s\n", title))
		for _, issue := range issues {
		  log.Println("issueListToString: ", issue.String())
			b.WriteString(fmt.Sprintf("* %s\n", issue.String()))
		}
	}
}

func chapterToString(b *bytes.Buffer, title string, section *Section) {
	if !section.IsEmpty() {
		b.WriteString(fmt.Sprintf("## %s\n%s", title, section.String()))
	}
}
