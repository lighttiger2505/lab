package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Pipeline interface {
	ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error)
}

type PipelineClient struct {
	Client *gitlab.Client
}

func NewPipelineClient(client *gitlab.Client) Pipeline {
	return &PipelineClient{Client: client}
}

func (c *PipelineClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	pipelines, _, err := c.Client.Pipelines.ListProjectPipelines(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipelines. Error: %s", err.Error())
	}
	return pipelines, nil
}

func (c *PipelineClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error) {
	jobs, _, err := c.Client.Jobs.ListPipelineJobs(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list pipeline jobs. Error: %s", err.Error())
	}
	return jobs, nil
}

type MockPipelineClient struct {
	MockProjectPipelines    func(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error)
	MockProjectPipelineJobs func(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error)
}

func (m *MockPipelineClient) ProjectPipelines(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
	return m.MockProjectPipelines(repositoryName, opt)
}

func (m *MockPipelineClient) ProjectPipelineJobs(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error) {
	return m.MockProjectPipelineJobs(repositoryName, opt, pid)
}
