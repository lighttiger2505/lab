package issue

import (
	"fmt"

	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

func detail(client lab.Issue, project string, iid int) (string, error) {
	issue, err := client.GetIssue(iid, project)
	if err != nil {
		return "", err
	}
	return issueDetailOutput(issue), nil
}

func issueDetailOutput(issue *gitlab.Issue) string {
	base := `#%d
Title: %s
Assignee: %s
Author: %s
State: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		issue.IID,
		issue.Title,
		issue.Assignee.Name,
		issue.Author.Name,
		issue.State,
		issue.CreatedAt.String(),
		issue.UpdatedAt.String(),
		issue.Description,
	)
	return detial
}
