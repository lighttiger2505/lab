package issue

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/config"
	"github.com/lighttiger2505/lab/internal/gitutil"
	gitlab "github.com/xanzy/go-gitlab"
)

var createdAt, _ = time.Parse("2006-01-02", "2018-02-14")
var updatedAt, _ = time.Parse("2006-01-02", "2018-03-14")
var issue = &gitlab.Issue{
	IID:   12,
	Title: "Title12",
	State: "State12",
	Assignee: struct {
		ID        int    `json:"id"`
		State     string `json:"state"`
		WebURL    string `json:"web_url"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	}{
		Name: "AssigneeName",
	},
	Author: struct {
		ID        int    `json:"id"`
		State     string `json:"state"`
		WebURL    string `json:"web_url"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Username  string `json:"username"`
	}{
		Name: "AuthorName",
	},
	CreatedAt:   &createdAt,
	UpdatedAt:   &updatedAt,
	Description: "Description",
}

func Test_createMethod_Process(t *testing.T) {
	type fields struct {
		client api.Issue
		opt    *CreateUpdateOption
		pInfo  *gitutil.GitLabProjectInfo
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "create input all issue value",
			fields: fields{
				client: &api.MockLabIssueClient{
					MockCreateIssue: func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.CreateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("desc"),
							AssigneeIDs: []int{13},
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "title",
					Message:    "desc",
					AssigneeID: 13,
				},
				pInfo: &gitutil.GitLabProjectInfo{
					Project: "group/project",
					Profile: &config.Profile{},
				},
			},
			want:    "12",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &createMethod{
				client: tt.fields.client,
				opt:    tt.fields.opt,
				pInfo:  tt.fields.pInfo,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("createMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createMethod.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createOnEditorMethod_Process(t *testing.T) {
	type fields struct {
		issueClient      api.Issue
		repositoryClient api.Repository
		opt              *CreateUpdateOption
		editFunc         func(program, file string) error
		pInfo            *gitutil.GitLabProjectInfo
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "create input all issue value",
			fields: fields{
				issueClient: &api.MockLabIssueClient{
					MockCreateIssue: func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.CreateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("desc"),
							AssigneeIDs: []int{13},
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				repositoryClient: &api.MockRepositoryClient{
					MockGetFile: func(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error) {
						return "template", nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "title",
					Message:    "desc",
					AssigneeID: 13,
				},
				pInfo: &gitutil.GitLabProjectInfo{
					Project: "group/project",
					Profile: &config.Profile{},
				},
				editFunc: func(program, file string) error { return nil },
			},
			want:    "12",
			wantErr: false,
		},
		{
			name: "use template",
			fields: fields{
				issueClient: &api.MockLabIssueClient{
					MockCreateIssue: func(opt *gitlab.CreateIssueOptions, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.CreateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("desc"),
							AssigneeIDs: []int{13},
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				repositoryClient: &api.MockRepositoryClient{
					MockGetFile: func(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error) {
						return "template", nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "title",
					Message:    "desc",
					Template:   "template",
					AssigneeID: 13,
				},
				pInfo: &gitutil.GitLabProjectInfo{
					Project: "group/project",
					Profile: &config.Profile{},
				},
				editFunc: func(program, file string) error { return nil },
			},
			want:    "12",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &createOnEditorMethod{
				issueClient:      tt.fields.issueClient,
				repositoryClient: tt.fields.repositoryClient,
				opt:              tt.fields.opt,
				editFunc:         tt.fields.editFunc,
				pInfo:            tt.fields.pInfo,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("createOnEditorMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createOnEditorMethod.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
