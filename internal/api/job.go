package api

import (
	"fmt"
	"io"

	gitlab "github.com/xanzy/go-gitlab"
)

type Job interface {
	GetProjectJobs(opt *gitlab.ListJobsOptions, repositoryName string) ([]gitlab.Job, error)
	GetJob(repositoryName string, jobID int) (*gitlab.Job, error)
	GetTraceFile(repositoryName string, jobID int) (io.Reader, error)
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
		return nil, fmt.Errorf("Failed list project jobs. %s", err.Error())
	}
	return jobs, nil
}

func (c *JobClient) GetJob(repositoryName string, jobID int) (*gitlab.Job, error) {
	job, _, err := c.Client.Jobs.GetJob(repositoryName, jobID)
	if err != nil {
		return nil, fmt.Errorf("Failed get job. %s", err.Error())
	}
	return job, nil
}

func (c *JobClient) GetTraceFile(repositoryName string, jobID int) (io.Reader, error) {
	trace, _, err := c.Client.Jobs.GetTraceFile(repositoryName, jobID)
	if err != nil {
		return nil, fmt.Errorf("Failed get trace file. %s", err.Error())
	}
	return trace, nil
}

type MockLabJobClient struct {
	Job
	MockGetProjectJobs func(opt *gitlab.ListJobsOptions, repositoryName string) ([]*gitlab.Job, error)
}

func (m *MockLabJobClient) GetProjectJobs(opt *gitlab.ListJobsOptions, repositoryName string) ([]*gitlab.Job, error) {
	return m.MockGetProjectJobs(opt, repositoryName)
}
