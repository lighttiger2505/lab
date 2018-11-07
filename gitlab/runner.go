package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Runner interface {
	ListRunners(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error)
}

type RunnerClient struct {
	Client *gitlab.Client
}

func NewRunnerClient(client *gitlab.Client) Runner {
	return &RunnerClient{Client: client}
}

func (c *RunnerClient) ListRunners(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error) {
	res, _, err := c.Client.Runners.ListRunners(opt)
	if err != nil {
		return nil, fmt.Errorf("failed list tree. %s", err.Error())
	}
	return res, nil
}

type MockRunnerClient struct {
	MockListRunners func(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error)
}

func (m *MockRunnerClient) ListRunners(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error) {
	return m.MockListRunners(opt)
}
