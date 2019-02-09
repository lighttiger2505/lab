package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/internal/api"
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
	if opt.MilestoneID != 0 {
		createIssueOption.MilestoneID = gitlab.Int(opt.MilestoneID)
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
	client  api.Issue
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
	issueClient      api.Issue
	repositoryClient api.Repository
	opt              *CreateUpdateOption
	editFunc         func(program, file string) error
	project          string
}

const templateDir = ".gitlab/issue_templates"

func (m *createOnEditorMethod) Process() (string, error) {
	templateFilename := m.opt.Template
	var template string
	if templateFilename != "" {
		filename := templateDir + "/" + templateFilename
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
