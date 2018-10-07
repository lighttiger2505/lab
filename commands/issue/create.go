package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
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

func makeIssueTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}

type createMethod struct {
	internal.Method
	client  lab.Issue
	opt     *CreateUpdateOption
	project string
}

func (m *createMethod) Process() (string, error) {
	issue, err := m.client.CreateIssue(
		makeCreateIssueOptions(m.opt, m.opt.Title, m.opt.Message),
		m.project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", issue.IID), nil
}

type createOnEditorMethod struct {
	internal.Method
	issueClient      lab.Issue
	repositoryClient lab.Repository
	opt              *CreateUpdateOption
	editFunc         func(program, file string) error
	project          string
}

func (m *createOnEditorMethod) Process() (string, error) {
	templateFilename := m.opt.Template
	var template string
	if templateFilename != "" {
		filename := TemplateDir + "/" + templateFilename
		res, err := m.repositoryClient.GetFile(
			m.project,
			filename,
			makeIssueTemplateOption(),
		)
		if err != nil {
			return "", err
		}
		template = res
	}

	var title, message string
	title = m.opt.Title
	message = template
	if m.opt.Message != "" {
		message = m.opt.Message
	}

	content := editIssueMessage(title, message)
	title, message, err := editIssueTitleAndDesc(content, m.editFunc)
	if err != nil {
		return "", err
	}
	issue, err := m.issueClient.CreateIssue(
		makeCreateIssueOptions(m.opt, title, message),
		m.project,
	)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", issue.IID), nil
}
