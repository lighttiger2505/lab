package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type Job interface {
	GetProjectJobs(opt *gitlab.ListJobsOptions, repositoryName string) ([]gitlab.Job, error)
}

type JobClient struct {
	Job
	Client *gitlab.Client
}

func NewJobClient(client *gitlab.Client) *JobClient {
	return &JobClient{Client: client}
}

func (c *JobClient) GetProjectJobs(opt *gitlab.ListJobsOptions, repositoryName string) ([]gitlab.Job, error) {
	jobs, _, err := c.Client.Jobs.ListProjectJobs(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}
	return jobs, nil
}

type MockLabJobClient struct {
	Job
	MockGetProjectJobs func(opt *gitlab.ListJobsOptions, repositoryName string) ([]*gitlab.Job, error)
}

func (m *MockLabJobClient) GetProjectJobs(opt *gitlab.ListJobsOptions, repositoryName string) ([]*gitlab.Job, error) {
	return m.MockGetProjectJobs(opt, repositoryName)
}
