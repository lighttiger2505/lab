package issue

import (
	"testing"

	"github.com/lighttiger2505/lab/internal/api"
	gitlab "github.com/xanzy/go-gitlab"
)

func Test_listMethod_Process(t *testing.T) {
	issues := []*gitlab.Issue{
		&gitlab.Issue{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo/issues/12"},
		&gitlab.Issue{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo/issues/13"},
	}

	type fields struct {
		client  api.Issue
		opt     *ListOption
		project string
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
				client: &api.MockLabIssueClient{
					MockGetProjectIssues: func(opt *gitlab.ListProjectIssuesOptions, repositoryName string) ([]*gitlab.Issue, error) {
						return issues, nil
					},
				},
				project: "group/project",
				opt:     &ListOption{},
			},
			want:    "12  Title12\n13  Title13",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &listMethod{
				client:  tt.fields.client,
				opt:     tt.fields.opt,
				project: tt.fields.project,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("listMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unmatch output \ngot: %#v\nwant:%#v", got, tt.want)
			}
		})
	}
}

func Test_listAllMethod_Process(t *testing.T) {
	issues := []*gitlab.Issue{
		&gitlab.Issue{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo/issues/12"},
		&gitlab.Issue{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo/issues/13"},
	}

	type fields struct {
		client api.Issue
		opt    *ListOption
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
				client: &api.MockLabIssueClient{
					MockGetAllProjectIssues: func(opt *gitlab.ListIssuesOptions) ([]*gitlab.Issue, error) {
						return issues, nil
					},
				},
				opt: &ListOption{},
			},
			want:    "namespace/repo  12  Title12\nnamespace/repo  13  Title13",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &listAllMethod{
				client: tt.fields.client,
				opt:    tt.fields.opt,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("listAllMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unmatch output \ngot: %#v\nwant:%#v", got, tt.want)
			}
		})
	}
}
