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

type MockLabClient struct {
	MockIssues              func(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	MockProjectIssues       func(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MockMergeRequest        func(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	MockProjectMergeRequest func(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	MockCreateIssue         func(baseurl, token string, opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	MockCreateMergeRequest  func(baseurl, token string, opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
}

func (m *MockLabClient) Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	return m.MockIssues(baseurl, token, opt)
}

func (m *MockLabClient) ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	return m.MockProjectIssues(baseurl, token, opt, repositoryName)
}

func (m *MockLabClient) MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	return m.MockMergeRequest(baseurl, token, opt)
}

func (m *MockLabClient) ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	return m.MockProjectMergeRequest(baseurl, token, opt, repositoryName)
}

func (m *MockLabClient) CreateIssue(baseurl, token string, opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	return m.MockCreateIssue(baseurl, token, opt, repositoryName)
}

func (m *MockLabClient) CreateMergeRequest(baseurl, token string, opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockCreateMergeRequest(baseurl, token, opt, repositoryName)
}
