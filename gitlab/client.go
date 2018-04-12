package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Client interface {
	// Issue
	GetIssue(pid int, repositoryName string) (*gitlab.Issue, error)
	Issues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	ProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error)
	// Merge Request
	GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error)
	MergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	ProjectMergeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error)
	// Project
	Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
	// Pipeline
	ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]gitlab.Job, error)
	// Lint
	Lint(content string) (*gitlab.LintResult, error)
	// User
	ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
	Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
}

type LabClient struct {
	Client *gitlab.Client
}

func NewLabClient(client *gitlab.Client) *LabClient {
	return &LabClient{Client: client}
}

func (l *LabClient) GetIssue(pid int, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := l.Client.Issues.GetIssue(repositoryName, pid)
	if err != nil {
		return nil, fmt.Errorf("Failed get issue. %s", err.Error())
	}
	return issue, nil
}

func (l *LabClient) Issues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	issues, _, err := l.Client.Issues.ListIssues(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}
	return issues, nil
}

func (l *LabClient) ProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	issues, _, err := l.Client.Issues.ListProjectIssues(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}
	return issues, nil
}

func (l *LabClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := l.Client.Issues.CreateIssue(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed create issue. %s", err.Error())
	}
	return issue, nil
}

func (l *LabClient) UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := l.Client.Issues.UpdateIssue(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed update issue. %s", err.Error())
	}
	return issue, nil
}

func (l *LabClient) GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.GetMergeRequest(repositoryName, pid)
	if err != nil {
		return nil, fmt.Errorf("Failed get merge request. %s", err.Error())
	}
	return mergeRequest, nil
}

func (l *LabClient) MergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	mergeRequests, _, err := l.Client.MergeRequests.ListMergeRequests(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (l *LabClient) ProjectMergeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	mergeRequests, _, err := l.Client.MergeRequests.ListProjectMergeRequests(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (l *LabClient) CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.CreateMergeRequest(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}
	return mergeRequest, nil
}

func (l *LabClient) UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.UpdateMergeRequest(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed get merge request. %s", err.Error())
	}
	return mergeRequest, nil
}

func (l *LabClient) Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := l.Client.Projects.ListProjects(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list projects. Error: %s", err.Error())
	}
	return projects, nil
}

func (l *LabClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	pipelines, _, err := l.Client.Pipelines.ListProjectPipelines(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipelines. Error: %s", err.Error())
	}
	return pipelines, nil
}

func (l *LabClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]gitlab.Job, error) {
	jobs, _, err := l.Client.Jobs.ListPipelineJobs(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipeline jobs. Error: %s", err.Error())
	}
	return jobs, nil
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
	// Issue
	MockGetIssue      func(pid int, repositoryName string) (*gitlab.Issue, error)
	MockIssues        func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	MockProjectIssues func(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MockCreateIssue   func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	MockUpdateIssue   func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error)
	// Merge Request
	MockGetMergeRequest     func(pid int, repositoryName string) (*gitlab.MergeRequest, error)
	MockMergeRequest        func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	MockProjectMergeRequest func(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	MockCreateMergeRequest  func(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	MockUpdateMergeRequest  func(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error)
	// Project
	MockProjects func(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
	// Pipeline
	MockProjectPipelines    func(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	MockProjectPipelineJobs func(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]gitlab.Job, error)
	// Lint
	MockLint func(content string) (*gitlab.LintResult, error)
	// User
	MockProjectUsers func(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
	MockUsers        func(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
}

func (m *MockLabClient) Issues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	return m.MockIssues(opt)
}

func (m *MockLabClient) GetIssue(pid int, repositoryName string) (*gitlab.Issue, error) {
	return m.MockGetIssue(pid, repositoryName)
}

func (m *MockLabClient) ProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	return m.MockProjectIssues(opt, repositoryName)
}

func (m *MockLabClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	return m.MockCreateIssue(opt, repositoryName)
}

func (m *MockLabClient) UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
	return m.MockUpdateIssue(opt, pid, repositoryName)
}

func (m *MockLabClient) GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockGetMergeRequest(pid, repositoryName)
}

func (m *MockLabClient) MergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	return m.MockMergeRequest(opt)
}

func (m *MockLabClient) ProjectMergeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	return m.MockProjectMergeRequest(opt, repositoryName)
}

func (m *MockLabClient) CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockCreateMergeRequest(opt, repositoryName)
}

func (m *MockLabClient) UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockUpdateMergeRequest(opt, pid, repositoryName)
}

func (m *MockLabClient) Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	return m.MockProjects(opt)
}

func (m *MockLabClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	return m.MockProjectPipelines(repositoryName, opt)
}

func (m *MockLabClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]gitlab.Job, error) {
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
