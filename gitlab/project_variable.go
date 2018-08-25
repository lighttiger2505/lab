package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type ProjectVariable interface {
	GetVariables(repositoryName string) ([]*gitlab.ProjectVariable, error)
}

type ProjectVariableClient struct {
	ProjectVariable
	Client *gitlab.Client
}

func NewProjectVariableClient(client *gitlab.Client) ProjectVariable {
	return &ProjectVariableClient{Client: client}
}

func (c *ProjectVariableClient) GetVariables(repositoryName string) ([]*gitlab.ProjectVariable, error) {
	vals, _, err := c.Client.ProjectVariables.ListVariables(repositoryName)
	if err != nil {
		return nil, fmt.Errorf("failed list project variables. %s", err.Error())
	}
	return vals, nil
}

type MockProjectVariableClient struct {
	ProjectVariable
	MockGetVariables func(repositoryName string) ([]*gitlab.ProjectVariable, error)
}

func (c *MockProjectVariableClient) GetVariables(repositoryName string) ([]*gitlab.ProjectVariable, error) {
	return c.MockGetVariables(repositoryName)
}
