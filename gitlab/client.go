package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Client interface {
	Issues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	ProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	ProjectMergeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
	ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
}

type LabClient struct {
	Client *gitlab.Client
}

func NewLabClient(client *gitlab.Client) *LabClient {
	return &LabClient{Client: client}
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

func (l *LabClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := l.Client.Issues.CreateIssue(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}
	return issue, nil
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

type MockLabClient struct {
	Client
	MockIssues              func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	MockProjectIssues       func(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MockMergeRequest        func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	MockProjectMergeRequest func(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	MockCreateIssue         func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	MockCreateMergeRequest  func(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	MockProjects            func(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
}

func (m *MockLabClient) Issues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	return m.MockIssues(opt)
}

func (m *MockLabClient) ProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	return m.MockProjectIssues(opt, repositoryName)
}

func (m *MockLabClient) MergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	return m.MockMergeRequest(opt)
}

func (m *MockLabClient) ProjectMergeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	return m.MockProjectMergeRequest(opt, repositoryName)
}

func (m *MockLabClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	return m.MockCreateIssue(opt, repositoryName)
}

func (m *MockLabClient) CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockCreateMergeRequest(opt, repositoryName)
}

func (m *MockLabClient) Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	return m.MockProjects(opt)
}
