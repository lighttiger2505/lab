package issue

import (
	"fmt"

	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

func makeCreateIssueOptions(opt *CreateUpdateOption, title, description string) *gitlab.CreateIssueOptions {
	createIssueOption := &gitlab.CreateIssueOptions{
		Title:       gitlab.String(title),
		Description: gitlab.String(description),
	}
	if opt.AssigneeID != 0 {
		createIssueOption.AssigneeIDs = []int{opt.AssigneeID}
	}
	return createIssueOption
}

func create(client lab.Issue, project string, opt *CreateUpdateOption) (string, error) {
	// Do create issue
	issue, err := client.CreateIssue(
		makeCreateIssueOptions(opt, opt.Title, opt.Message),
		project,
	)
	if err != nil {
		return "", err
	}

	// Print created Issue IID
	return fmt.Sprintf("%d", issue.IID), nil
}

func createOnEditor(
	client lab.Issue,
	project string,
	templateContent string,
	opt *CreateUpdateOption,
	editFunc func(program, file string) error,
) (string, error) {

	var title, message string
	title = opt.Title
	message = templateContent
	if opt.Message != "" {
		message = opt.Message
	}

	template := editIssueMessage(title, message)
	title, message, err := editIssueTitleAndDesc(template, editFunc)
	if err != nil {
		return "", err
	}

	// Do create issue
	issue, err := client.CreateIssue(
		makeCreateIssueOptions(opt, title, message),
		project,
	)
	if err != nil {
		return "", err
	}

	// Print created Issue IID
	return fmt.Sprintf("%d", issue.IID), nil
}
