package issue

import (
	"testing"
	"time"

	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

func Test_detailMethod_Process(t *testing.T) {
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
	notes := []*gitlab.Note{
		&gitlab.Note{
			ID:    1,
			Body:  "body",
			Title: "title",
			Author: struct {
				ID        int    `json:"id"`
				Username  string `json:"username"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				State     string `json:"state"`
				AvatarURL string `json:"avatar_url"`
				WebURL    string `json:"web_url"`
			}{
				Name: "author1",
			},
			UpdatedAt: &updatedAt,
			CreatedAt: &createdAt,
		},
		&gitlab.Note{
			ID:    2,
			Body:  "body",
			Title: "title",
			Author: struct {
				ID        int    `json:"id"`
				Username  string `json:"username"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				State     string `json:"state"`
				AvatarURL string `json:"avatar_url"`
				WebURL    string `json:"web_url"`
			}{
				Name: "author2",
			},
			UpdatedAt: &updatedAt,
			CreatedAt: &createdAt,
		},
	}

	// Define sub tests
	type fields struct {
		issueClient lab.Issue
		noteClient  lab.Note
		id          int
		project     string
		opt         *ShowOption
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "show issue",
			fields: fields{
				issueClient: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
				},
				noteClient: &lab.MockNoteClient{
					MockGetIssueNotes: func(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
						return notes, nil
					},
				},
				project: "group/project",
				id:      12,
				opt: &ShowOption{
					NoComment: true,
				},
			},
			want: `12 Title12 [State12] (created by @AuthorName, 2018-02-14 00:00:00 +0000 UTC)
Assignee: AssigneeName
Milestone: 
Labels: 

Description`,
			wantErr: false,
		},
		{
			name: "show issue with note",
			fields: fields{
				issueClient: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
				},
				noteClient: &lab.MockNoteClient{
					MockGetIssueNotes: func(repositoryName string, iid int, opt *gitlab.ListIssueNotesOptions) ([]*gitlab.Note, error) {
						return notes, nil
					},
				},
				project: "group/project",
				id:      12,
				opt: &ShowOption{
					NoComment: false,
				},
			},
			want: `12 Title12 [State12] (created by @AuthorName, 2018-02-14 00:00:00 +0000 UTC)
Assignee: AssigneeName
Milestone: 
Labels: 

Description

comment 1 (created by @author1, 2018-02-14 00:00:00 +0000 UTC)

body

comment 2 (created by @author2, 2018-02-14 00:00:00 +0000 UTC)

body`,
			wantErr: false,
		},
	}

	// Do tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &detailMethod{
				issueClient: tt.fields.issueClient,
				noteClient:  tt.fields.noteClient,
				id:          tt.fields.id,
				project:     tt.fields.project,
				opt:         tt.fields.opt,
			}
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("detailMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("detailMethod.Process() = \ngot: %#v\nwant:%#v", got, tt.want)
			}
		})
	}
}
