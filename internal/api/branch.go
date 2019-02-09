package api

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Branch interface {
	GetBranch(project string, branch string) (*gitlab.Branch, error)
	ListBranches(project string, opt *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error)
}

type BranchClient struct {
	Client *gitlab.Client
}

func NewBranchClient(client *gitlab.Client) Branch {
	return &BranchClient{Client: client}
}

func (c *BranchClient) GetBranch(project string, branch string) (*gitlab.Branch, error) {
	result, _, err := c.Client.Branches.GetBranch(project, branch)
	if err != nil {
		return nil, fmt.Errorf("Failed list branches. Error: %s", err.Error())
	}
	return result, nil
}

func (c *BranchClient) ListBranches(project string, opt *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error) {
	results, _, err := c.Client.Branches.ListBranches(project, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list branches. Error: %s", err.Error())
	}
	return results, nil
}

type MockBranchClient struct {
	MockGetBranch    func(project string, branch string) (*gitlab.Branch, error)
	MockListBranches func(project string, opt *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error)
}

func (m *MockBranchClient) GetBranch(project string, branch string) (*gitlab.Branch, error) {
	return m.MockGetBranch(project, branch)
}

func (m *MockBranchClient) ListBranches(project string, opt *gitlab.ListBranchesOptions) ([]*gitlab.Branch, error) {
	return m.MockListBranches(project, opt)
}
