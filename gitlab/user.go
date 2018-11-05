package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type User interface {
	Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
	ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
}

type UserClient struct {
	Client *gitlab.Client
}

func NewUserClient(client *gitlab.Client) *UserClient {
	return &UserClient{Client: client}
}

func (c *UserClient) Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
	results, _, err := c.Client.Users.ListUsers(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list users. Error: %s", err.Error())
	}
	return results, nil
}

func (c *UserClient) ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
	results, _, err := c.Client.Projects.ListProjectsUsers(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project users. Error: %s", err.Error())
	}
	return results, nil
}

type MockUserClient struct {
	MockUsers        func(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error)
	MockProjectUsers func(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error)
}

func (m *MockUserClient) Users(opt *gitlab.ListUsersOptions) ([]*gitlab.User, error) {
	return m.MockUsers(opt)
}

func (m *MockUserClient) ProjectUsers(repositoryName string, opt *gitlab.ListProjectUserOptions) ([]*gitlab.ProjectUser, error) {
	return m.MockProjectUsers(repositoryName, opt)
}
