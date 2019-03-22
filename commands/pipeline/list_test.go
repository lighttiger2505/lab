package pipeline

import (
	"testing"

	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/config"
	"github.com/lighttiger2505/lab/internal/gitutil"
	gitlab "github.com/xanzy/go-gitlab"
)

func Test_listMethod_Process(t *testing.T) {

	pipelines := gitlab.PipelineList{
		struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
		}{
			ID:     1,
			Status: "status1",
			Ref:    "ref1",
			Sha:    "sha1",
		},
		struct {
			ID     int    `json:"id"`
			Status string `json:"status"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
		}{
			ID:     2,
			Status: "status2",
			Ref:    "ref2",
			Sha:    "sha2",
		},
	}
	type fields struct {
		client api.Pipeline
		opt    *ListOption
		pInfo  *gitutil.GitLabProjectInfo
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "nomal",
			fields: fields{
				client: &api.MockPipelineClient{
					MockProjectPipelines: func(repositoryName string, opt *gitlab.ListProjectPipelinesOptions) (gitlab.PipelineList, error) {
						return pipelines, nil
					},
				},
				opt: &ListOption{},
				pInfo: &gitutil.GitLabProjectInfo{
					Project: "group/project",
					Profile: &config.Profile{},
				},
			},
			want:    "1  status1  ref1  sha1\n2  status2  ref2  sha2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &listMethod{
				client: tt.fields.client,
				opt:    tt.fields.opt,
				pInfo: &gitutil.GitLabProjectInfo{
					Project: "group/project",
					Profile: &config.Profile{},
				},
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("listMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unmatch output\ngot: %#v\nwant %#v", got, tt.want)
			}
		})
	}
}

func Test_listJobMethod_Process(t *testing.T) {
	jobs := []*gitlab.Job{
		&gitlab.Job{
			ID:     1,
			Status: "Status1",
			Ref:    "Ref1",
			Commit: &gitlab.Commit{ShortID: "ShortID1"},
			User:   &gitlab.User{Username: "Username1"},
			Stage:  "Stage1",
			Name:   "Name1",
		},
		&gitlab.Job{
			ID:     2,
			Status: "Status2",
			Ref:    "Ref2",
			Commit: &gitlab.Commit{ShortID: "ShortID2"},
			User:   &gitlab.User{Username: "Username2"},
			Stage:  "Stage2",
			Name:   "Name2",
		},
	}
	type fields struct {
		client  api.Pipeline
		opt     *ListOption
		project string
		id      int
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "nomal",
			fields: fields{
				client: &api.MockPipelineClient{
					MockProjectPipelineJobs: func(repositoryName string, opt *gitlab.ListJobsOptions, pid int) ([]*gitlab.Job, error) {
						return jobs, nil
					},
				},
				opt:     &ListOption{},
				project: "group/project",
				id:      12,
			},
			want:    "1  Status1  Ref1  ShortID1  Username1  Stage1  Name1\n2  Status2  Ref2  ShortID2  Username2  Stage2  Name2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &listJobMethod{
				client:  tt.fields.client,
				opt:     tt.fields.opt,
				project: tt.fields.project,
				id:      tt.fields.id,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("listJobMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unmatch output\ngot: %#v\nwant %#v", got, tt.want)
			}
		})
	}
}
