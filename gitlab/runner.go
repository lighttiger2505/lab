package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Runner interface {
	ListAllRunners(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error)
	ListProjectRunners(pid string, opt *gitlab.ListProjectRunnersOptions) ([]*gitlab.Runner, error)
	GetRunnerDetails(id int) (*gitlab.RunnerDetails, error)
	RemoveRunner(iid int) error
}

type RunnerClient struct {
	Client *gitlab.Client
}

func NewRunnerClient(client *gitlab.Client) Runner {
	return &RunnerClient{Client: client}
}

func (c *RunnerClient) ListAllRunners(opt *gitlab.ListRunnersOptions) ([]*gitlab.Runner, error) {
	res, _, err := c.Client.Runners.ListAllRunners(opt)
	if err != nil {
		return nil, fmt.Errorf("failed list runners. %s", err.Error())
	}
	return res, nil
}

func (c *RunnerClient) ListProjectRunners(pid string, opt *gitlab.ListProjectRunnersOptions) ([]*gitlab.Runner, error) {
	res, _, err := c.Client.Runners.ListProjectRunners(pid, opt)
	if err != nil {
		return nil, fmt.Errorf("failed list project runners. %s", err.Error())
	}
	return res, nil
}

func (c *RunnerClient) RemoveRunner(iid int) error {
	_, err := c.Client.Runners.RemoveRunner(iid)
	if err != nil {
		return fmt.Errorf("failed delete runner. %s", err.Error())
	}
	return nil
}

func (c *RunnerClient) GetRunnerDetails(id int) (*gitlab.RunnerDetails, error) {
	res, _, err := c.Client.Runners.GetRunnerDetails(id)
	if err != nil {
		return nil, fmt.Errorf("failed delete runner. %s", err.Error())
	}
	return res, nil
}
