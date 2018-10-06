package issue

import (
	"fmt"

	lab "github.com/lighttiger2505/lab/gitlab"
)

func updateIssue(client lab.Issue, project string, iid int, opt *CreateUpdateIssueOption) (string, error) {
	// Getting exist issue
	issue, err := client.GetIssue(iid, project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if opt.Title != "" {
		updatedTitle = opt.Title
	}
	if opt.Message != "" {
		updatedMessage = opt.Message
	}

	// Do update issue
	updatedIssue, err := client.UpdateIssue(
		makeUpdateIssueOption(opt, updatedTitle, updatedMessage),
		iid,
		project,
	)
	if err != nil {
		return "", err
	}

	// Print update Issue IID
	return fmt.Sprintf("%d", updatedIssue.IID), nil
}

func updateIssueOnEditor(client lab.Issue, project string, iid int, opt *CreateUpdateIssueOption, editFunc func(program, file string) error) (string, error) {
	// Getting exist issue
	issue, err := client.GetIssue(iid, project)
	if err != nil {
		return "", err
	}

	// Create new title or description
	updatedTitle := issue.Title
	updatedMessage := issue.Description
	if opt.Title != "" {
		updatedTitle = opt.Title
	}
	if opt.Message != "" {
		updatedMessage = opt.Message
	}

	// Starting editor for edit title and description
	template := editIssueMessage(updatedTitle, updatedMessage)
	title, message, err := editIssueTitleAndDesc(template, editFunc)
	if err != nil {
		return "", err
	}

	// Do update issue
	updatedIssue, err := client.UpdateIssue(
		makeUpdateIssueOption(opt, title, message),
		iid,
		project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", updatedIssue.IID), nil
}
