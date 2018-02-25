package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Client interface {
	Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	CreateIssue(baseurl, token string, opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	CreateMergeRequest(baseurl, token string, opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
}

type LabClient struct {
	Client
}

func NewLabClient() *LabClient {
	return &LabClient{}
}

func (g *LabClient) Issues(baseurl, token string, opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	issues, _, err := client.Issues.ListIssues(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}
	return issues, nil
}

func (g *LabClient) ProjectIssues(baseurl, token string, opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	issues, _, err := client.Issues.ListProjectIssues(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}
	return issues, nil
}

func (g *LabClient) MergeRequest(baseurl, token string, opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	mergeRequests, _, err := client.MergeRequests.ListMergeRequests(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (g *LabClient) ProjectMergeRequest(baseurl, token string, opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (l *LabClient) CreateIssue(baseurl, token string, opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	issue, _, err := client.Issues.CreateIssue(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}
	return issue, nil
}

func (l *LabClient) CreateMergeRequest(baseurl, token string, opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	client := gitlab.NewClient(nil, token)
	if err := client.SetBaseURL(baseurl); err != nil {
		return nil, fmt.Errorf("Invalid api url. %s", err.Error())
	}

	mergeRequest, _, err := client.MergeRequests.CreateMergeRequest(
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
	return m.CreateIssue(baseurl, token, opt, repositoryName)
}
