package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type ProjectVariable interface {
	GetVariables(repositoryName string) ([]*gitlab.ProjectVariable, error)
	CreateVariable(repositoryName string, opt *gitlab.CreateVariableOptions) (*gitlab.ProjectVariable, error)
	UpdateVariable(repositoryName string, key string, opt *gitlab.UpdateVariableOptions) (*gitlab.ProjectVariable, error)
	RemoveVariable(repositoryName string, key string) error
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

func (c *ProjectVariableClient) CreateVariable(repositoryName string, opt *gitlab.CreateVariableOptions) (*gitlab.ProjectVariable, error) {
	val, _, err := c.Client.ProjectVariables.CreateVariable(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("failed create project variables. %s", err.Error())
	}
	return val, nil
}

func (c *ProjectVariableClient) UpdateVariable(repositoryName string, key string, opt *gitlab.UpdateVariableOptions) (*gitlab.ProjectVariable, error) {
	val, _, err := c.Client.ProjectVariables.UpdateVariable(repositoryName, key, opt)
	if err != nil {
		return nil, fmt.Errorf("failed update project variables. %s", err.Error())
	}
	return val, nil
}

func (c *ProjectVariableClient) RemoveVariable(repositoryName string, key string) error {
	_, err := c.Client.ProjectVariables.RemoveVariable(repositoryName, key)
	if err != nil {
		return fmt.Errorf("failed update project variables. %s", err.Error())
	}
	return nil
}

type MockProjectVariableClient struct {
	ProjectVariable
	MockGetVariables   func(repositoryName string) ([]*gitlab.ProjectVariable, error)
	MockCreateVariable func(repositoryName string, opt *gitlab.CreateVariableOptions) (*gitlab.ProjectVariable, error)
	MockUpdateVariable func(repositoryName string, key string, opt *gitlab.UpdateVariableOptions) (*gitlab.ProjectVariable, error)
	MockRemoveVariable func(repositoryName string, key string) error
}

func (c *MockProjectVariableClient) GetVariables(repositoryName string) ([]*gitlab.ProjectVariable, error) {
	return c.MockGetVariables(repositoryName)
}

func (c *MockProjectVariableClient) CreateVariable(repositoryName string, opt *gitlab.CreateVariableOptions) (*gitlab.ProjectVariable, error) {
	return c.MockCreateVariable(repositoryName, opt)
}

func (c *MockProjectVariableClient) UpdateVariable(repositoryName string, key string, opt *gitlab.UpdateVariableOptions) (*gitlab.ProjectVariable, error) {
	return c.MockUpdateVariable(repositoryName, key, opt)
}

func (c *MockProjectVariableClient) RemoveVariable(repositoryName string, key string) error {
	return c.MockRemoveVariable(repositoryName, key)
}
