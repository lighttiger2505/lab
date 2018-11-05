package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Project interface {
	Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
}

type ProjectClient struct {
	Client *gitlab.Client
}

func NewProjectClient(client *gitlab.Client) *ProjectClient {
	return &ProjectClient{Client: client}
}

func (c *ProjectClient) Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	projects, _, err := c.Client.Projects.ListProjects(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list projects. Error: %s", err.Error())
	}
	return projects, nil
}

type MockProjectClient struct {
	MockProjects func(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error)
}

func (m *MockProjectClient) Projects(opt *gitlab.ListProjectsOptions) ([]*gitlab.Project, error) {
	return m.MockProjects(opt)
}
