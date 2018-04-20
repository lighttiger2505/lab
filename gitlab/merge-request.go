package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

type MergeRequest interface {
	GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error)
	GetAllProjectMergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	GetProjectMargeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error)
}

type MergeRequestClient struct {
	MergeRequest
	Client *gitlab.Client
}

func NewMergeRequestClient(client *gitlab.Client) *MergeRequestClient {
	return &MergeRequestClient{Client: client}
}

func (l *MergeRequestClient) GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.GetMergeRequest(repositoryName, pid)
	if err != nil {
		return nil, fmt.Errorf("Failed get merge request. %s", err.Error())
	}
	return mergeRequest, nil
}

func (l *MergeRequestClient) GetAllProjectMergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	mergeRequests, _, err := l.Client.MergeRequests.ListMergeRequests(opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (l *MergeRequestClient) GetProjectMargeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	mergeRequests, _, err := l.Client.MergeRequests.ListProjectMergeRequests(repositoryName, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func (l *MergeRequestClient) CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.CreateMergeRequest(
		repositoryName,
		opt,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}
	return mergeRequest, nil
}

func (l *MergeRequestClient) UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	mergeRequest, _, err := l.Client.MergeRequests.UpdateMergeRequest(repositoryName, pid, opt)
	if err != nil {
		return nil, fmt.Errorf("Failed get merge request. %s", err.Error())
	}
	return mergeRequest, nil
}

type MockLabMergeRequestClient struct {
	MergeRequest
	MockGetMergeRequest           func(pid int, repositoryName string) (*gitlab.MergeRequest, error)
	MockGetAllProjectMergeRequest func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error)
	MockGetProjectMargeRequest    func(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error)
	MockCreateMergeRequest        func(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error)
	MockUpdateMergeRequest        func(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error)
}

func (m *MockLabMergeRequestClient) GetMergeRequest(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockGetMergeRequest(pid, repositoryName)
}

func (m *MockLabMergeRequestClient) GetAllProjectMergeRequest(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
	return m.MockGetAllProjectMergeRequest(opt)
}

func (m *MockLabMergeRequestClient) GetProjectMargeRequest(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
	return m.MockGetProjectMargeRequest(opt, repositoryName)
}

func (m *MockLabMergeRequestClient) CreateMergeRequest(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockCreateMergeRequest(opt, repositoryName)
}

func (m *MockLabMergeRequestClient) UpdateMergeRequest(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
	return m.MockUpdateMergeRequest(opt, pid, repositoryName)
}
