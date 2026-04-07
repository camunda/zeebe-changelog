package gitlog

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	lineRegex = regexp.MustCompile(`(?im)^\s*(closes?|related|relates|merges?|back\s?ports?|resolved|resolves)\s+.*$`)
	idRegex   = regexp.MustCompile(`(\s#|(https?\:\/\/(www\.)?github\.com\/)?camunda\/(camunda|zeebe)(\/|#))(\d+)`)
)

func GetHistory(path, start, end string) string {
	err := validateAncestor(path, start, end)
	if err != nil {
		log.Fatal(err)
	}

	logRange := fmt.Sprintf("%s..%s", start, end)

	// use git command til git lib implements range feature, see https://github.com/src-d/go-git/issues/1166
	// Note: We removed the --since filter because it was incorrectly filtering out backported fixes
	// that were committed before the start tag was released. The git revision range (start..end)
	// already correctly determines which commits are new between revisions.
	command := exec.Command("git", "-C", path, "log", logRange, "--merges", "--")
	log.Println(command)
	out, err := command.CombinedOutput()
	result := string(out)

	if err != nil {
		log.Fatal(result, err)
	}

	return result
}

func validateAncestor(path, start, end string) error {
	command := exec.Command("git", "-C", path, "merge-base", "--is-ancestor", start, end)
	log.Println(command)
	out, err := command.CombinedOutput()
	if err == nil {
		return nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		log.Printf("warning: git range %s..%s has start that is not an ancestor of end; continuing anyway", start, end)
		return nil
	}

	return fmt.Errorf("unable to validate git range %s..%s: %s (%w)", start, end, strings.TrimSpace(string(out)), err)
}

func ExtractIssueIds(message string) []int {
	seen := map[int]bool{}
	var issueIds []int

	for _, line := range lineRegex.FindAllString(message, -1) {
		for _, match := range idRegex.FindAllStringSubmatch(line, -1) {
			issueId, err := strconv.Atoi(match[6])
			if err != nil {
				log.Fatalln("Cannot convert issue id", match[6], err)
			}

			if !seen[issueId] {
				seen[issueId] = true
				issueIds = append(issueIds, issueId)
			}
		}
	}

	return issueIds
}
