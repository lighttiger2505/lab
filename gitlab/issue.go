package gitlab

import (
	"fmt"
	"testing"

	gitlab "github.com/xanzy/go-gitlab"
)

type Issue interface {
	GetIssue(pid int, repositoryName string) (*gitlab.Issue, error)
	GetAllProjectIssues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	GetProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error)
}

type IssueClient struct {
	Issue
	Client *gitlab.Client
}

func NewIssueClient(client *gitlab.Client) *IssueClient {
	return &IssueClient{Client: client}
}

func (c *IssueClient) GetIssue(pid int, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := c.Client.Issues.GetIssue(repositoryName, pid)
	if err != nil {
		return nil, fmt.Errorf("Failed get issue. %s", err.Error())
	}
	return issue, nil
}

func (c *IssueClient) GetAllProjectIssues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	issues, _, err := c.Client.Issues.ListIssues(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}
	return issues, nil
}

func (c *IssueClient) GetProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	issues, _, err := c.Client.Issues.ListProjectIssues(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}
	return issues, nil
}

func (c *IssueClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := c.Client.Issues.CreateIssue(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed create issue. %s", err.Error())
	}
	return issue, nil
}

func (c *IssueClient) UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
	issue, _, err := c.Client.Issues.UpdateIssue(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed update issue. %s", err.Error())
	}
	return issue, nil
}

type MockLabIssueClient struct {
	Issue
	t                       *testing.T
	MockGetIssue            func(pid int, repositoryName string) (*gitlab.Issue, error)
	MockGetAllProjectIssues func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error)
	MockGetProjectIssues    func(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error)
	MockCreateIssue         func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error)
	MockUpdateIssue         func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error)
}

func (m *MockLabIssueClient) GetIssue(pid int, repositoryName string) (*gitlab.Issue, error) {
	return m.MockGetIssue(pid, repositoryName)
}

func (m *MockLabIssueClient) GetAllProjectIssues(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
	return m.MockGetAllProjectIssues(opt)
}

func (m *MockLabIssueClient) GetProjectIssues(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
	return m.MockGetProjectIssues(opt, repositoryName)
}

func (m *MockLabIssueClient) CreateIssue(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
	return m.MockCreateIssue(opt, repositoryName)
}

func (m *MockLabIssueClient) UpdateIssue(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
	return m.MockUpdateIssue(opt, pid, repositoryName)
}
