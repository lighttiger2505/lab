package issue

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	lab "github.com/lighttiger2505/lab/gitlab"
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
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarURL string `json:"avatar_url"`
		WebURL    string `json:"web_url"`
	}{
		Name: "AssigneeName",
	},
	Author: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Email     string     `json:"email"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AuthorName",
	},
	CreatedAt:   &createdAt,
	UpdatedAt:   &updatedAt,
	Description: "Description",
}

func Test_createMethod_Process(t *testing.T) {
	type fields struct {
		client  lab.Issue
		opt     *CreateUpdateOption
		project string
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
				client: &lab.MockLabIssueClient{
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
				project: "group/project",
			},
			want:    "12",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &createMethod{
				client:  tt.fields.client,
				opt:     tt.fields.opt,
				project: tt.fields.project,
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
		issueClient      lab.Issue
		repositoryClient lab.Repository
		opt              *CreateUpdateOption
		editFunc         func(program, file string) error
		project          string
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
				issueClient: &lab.MockLabIssueClient{
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
				repositoryClient: &lab.MockRepositoryClient{
					MockGetFile: func(repositoryName string, filename string, opt *gitlab.GetRawFileOptions) (string, error) {
						return "template", nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "title",
					Message:    "desc",
					AssigneeID: 13,
				},
				project:  "group/project",
				editFunc: func(program, file string) error { return nil },
			},
			want:    "12",
			wantErr: false,
		},
		{
			name: "use template",
			fields: fields{
				issueClient: &lab.MockLabIssueClient{
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
				repositoryClient: &lab.MockRepositoryClient{
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
				project:  "group/project",
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
				project:          tt.fields.project,
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
