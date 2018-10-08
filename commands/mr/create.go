package mr

import (
	"fmt"

	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

type createMethod struct {
	internal.Method
	client  lab.MergeRequest
	opt     *CreateUpdateMergeRequestOption
	project string
}

func (m *createMethod) Process() (string, error) {
	// Get source branch. current branch from local repository when non specific flags
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}
	if m.opt.SourceBranch != "" {
		// TODO Checking source branch exitst
		currentBranch = m.opt.SourceBranch
	}

	// Do create merge request
	mergeRequest, err := m.client.CreateMergeRequest(
		makeCreateMergeRequestOption(m.opt, m.opt.Title, m.opt.Message, currentBranch),
		m.project,
	)
	if err != nil {
		return "", err
	}

	// Print created merge request IID
	return fmt.Sprintf("%d", mergeRequest.IID), nil
}

type createOnEditorMethod struct {
	internal.Method
	client           lab.MergeRequest
	repositoryClient lab.Repository
	opt              *CreateUpdateMergeRequestOption
	project          string
	editFunc         func(program, file string) error
}

const templateDir = ".gitlab/merge_request_templates"

func (m *createOnEditorMethod) Process() (string, error) {
	templateFilename := m.opt.Template
	var template string
	if templateFilename != "" {
		filename := templateDir + "/" + templateFilename
		res, err := m.repositoryClient.GetFile(
			m.project,
			filename,
			makeMergeRequestTemplateOption(),
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

	content := editMergeRequestTemplate(title, message)
	title, message, err := editIssueTitleAndDesc(content, m.editFunc)
	if err != nil {
		return "", err
	}

	// Get source branch. current branch from local repository when non specific flags
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return "", err
	}
	if m.opt.SourceBranch != "" {
		// TODO Checking source branch exitst
		currentBranch = m.opt.SourceBranch
	}

	// Do create merge request
	mergeRequest, err := m.client.CreateMergeRequest(
		makeCreateMergeRequestOption(m.opt, title, message, currentBranch),
		m.project,
	)
	if err != nil {
		return "", err
	}

	// Print created merge request IID
	return fmt.Sprintf("%d", mergeRequest.IID), nil
}

func makeCreateMergeRequestOption(opt *CreateUpdateMergeRequestOption, title, description, branch string) *gitlab.CreateMergeRequestOptions {
	createMergeRequestOption := &gitlab.CreateMergeRequestOptions{
		Title:           gitlab.String(title),
		Description:     gitlab.String(description),
		SourceBranch:    gitlab.String(branch),
		TargetBranch:    gitlab.String(opt.TargetBranch),
		TargetProjectID: nil,
	}
	if opt.AssigneeID != 0 {
		createMergeRequestOption.AssigneeID = gitlab.Int(opt.AssigneeID)
	}
	return createMergeRequestOption
}

func makeMergeRequestTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}
