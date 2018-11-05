package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Client interface {
	// Pipeline
	ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error)
	// Lint
	Lint(content string) (*gitlab.LintResult, error)
	// User
	ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
	Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
}

type LabClient struct {
	Client *gitlab.Client
}

func (l *LabClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	pipelines, _, err := l.Client.Pipelines.ListProjectPipelines(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipelines. Error: %s", err.Error())
	}
	return pipelines, nil
}

func (l *LabClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error) {
	jobs, _, err := l.Client.Jobs.ListPipelineJobs(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipeline jobs. Error: %s", err.Error())
	}
	return jobs, nil
}

func NewLabClient(client *gitlab.Client) *LabClient {
	return &LabClient{Client: client}
}

func (l *LabClient) Lint(content string) (*gitlab.LintResult, error) {
	lintResult, _, err := l.Client.Validate.Lint(content)
	if err != nil {
		return nil, fmt.Errorf("Failed lint. Error: %s", err.Error())
	}
	return lintResult, nil
}

func (l *LabClient) Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
	results, _, err := l.Client.Users.ListUsers(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list users. Error: %s", err.Error())
	}
	return results, nil
}

func (l *LabClient) ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
	results, _, err := l.Client.Projects.ListProjectsUsers(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project users. Error: %s", err.Error())
	}
	return results, nil
}

type MockLabClient struct {
	Client
	// Pipeline
	MockProjectPipelines    func(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	MockProjectPipelineJobs func(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error)
	// Lint
	MockLint func(content string) (*gitlab.LintResult, error)
	// User
	MockProjectUsers func(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
	MockUsers        func(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
}

func (m *MockLabClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	return m.MockProjectPipelines(repositoryName, opt)
}

func (m *MockLabClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error) {
	return m.MockProjectPipelineJobs(repositoryName, opt, pid)
}

func (m *MockLabClient) Lint(content string) (*gitlab.LintResult, error) {
	return m.MockLint(content)
}

func (m *MockLabClient) Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
	return m.MockUsers(opt)
}

func (m *MockLabClient) ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
	return m.MockProjectUsers(repositoryName, opt)
}
