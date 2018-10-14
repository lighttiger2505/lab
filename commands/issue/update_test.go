package issue

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	gitlab "github.com/xanzy/go-gitlab"
)

func Test_updateMethod_Process(t *testing.T) {
	var issue = &gitlab.Issue{
		IID:   12,
		Title: "title",
		State: "state",
		Assignee: struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		}{
			ID: 24,
		},
		Description: "desc",
	}

	tests := []struct {
		name    string
		method  internal.Method
		want    string
		wantErr bool
	}{
		{
			name: "update all",
			method: &updateMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("newtitle"),
							Description: gitlab.String("newmessage"),
							StateEvent:  gitlab.String("newstate"),
							AssigneeIDs: []int{13},
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "newtitle",
					Message:    "newmessage",
					StateEvent: "newstate",
					AssigneeID: 13,
				},
				project: "group/project",
				id:      12,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "update title only",
			method: &updateMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("newtitle"),
							Description: gitlab.String("desc"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "newtitle",
					Message:    "",
					StateEvent: "",
					AssigneeID: 0,
				},
				project: "group/project",
				id:      12,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "update message only",
			method: &updateMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("newmessage"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "",
					Message:    "newmessage",
					StateEvent: "",
					AssigneeID: 0,
				},
				project: "group/project",
				id:      12,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "update state only",
			method: &updateMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("desc"),
							StateEvent:  gitlab.String("newstate"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "",
					Message:    "",
					StateEvent: "newstate",
					AssigneeID: 0,
				},
				project: "group/project",
				id:      12,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "update assignee only",
			method: &updateMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
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
					Title:      "",
					Message:    "",
					StateEvent: "",
					AssigneeID: 13,
				},
				project: "group/project",
				id:      12,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.method
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("updateMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("updateMethod.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateOnEditorMethod_Process(t *testing.T) {
	var issue = &gitlab.Issue{
		IID:   12,
		Title: "title",
		State: "state",
		Assignee: struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Username  string `json:"username"`
			State     string `json:"state"`
			AvatarURL string `json:"avatar_url"`
			WebURL    string `json:"web_url"`
		}{
			ID: 24,
		},
		Description: "desc",
	}

	tests := []struct {
		name    string
		method  internal.Method
		want    string
		wantErr bool
	}{
		{
			name: "update all",
			method: &updateOnEditorMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("newtitle"),
							Description: gitlab.String("newmessage"),
							StateEvent:  gitlab.String("newstate"),
							AssigneeIDs: []int{13},
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "newtitle",
					Message:    "newmessage",
					StateEvent: "newstate",
					AssigneeID: 13,
				},
				project:  "group/project",
				id:       12,
				editFunc: func(program, file string) error { return nil },
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "change title only",
			method: &updateOnEditorMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("newtitle"),
							Description: gitlab.String("desc"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "newtitle",
					Message:    "",
					StateEvent: "",
					AssigneeID: 0,
				},
				project:  "group/project",
				id:       12,
				editFunc: func(program, file string) error { return nil },
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "change message only",
			method: &updateOnEditorMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("newmessage"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "",
					Message:    "newmessage",
					StateEvent: "",
					AssigneeID: 0,
				},
				project:  "group/project",
				id:       12,
				editFunc: func(program, file string) error { return nil },
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "change state only",
			method: &updateOnEditorMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
							Title:       gitlab.String("title"),
							Description: gitlab.String("desc"),
							StateEvent:  gitlab.String("newstate"),
						}
						if diff := cmp.Diff(got, want); diff != "" {
							t.Errorf("invalide arg (-got +want)\n%s", diff)
						}
						return issue, nil
					},
				},
				opt: &CreateUpdateOption{
					Title:      "",
					Message:    "",
					StateEvent: "newstate",
					AssigneeID: 0,
				},
				project:  "group/project",
				id:       12,
				editFunc: func(program, file string) error { return nil },
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "change assignee only",
			method: &updateOnEditorMethod{
				client: &lab.MockLabIssueClient{
					MockGetIssue: func(pid int, repositoryName string) (*gitlab.Issue, error) {
						return issue, nil
					},
					MockUpdateIssue: func(opt *gitlab.UpdateIssueOptions, pid int, repositoryName string) (*gitlab.Issue, error) {
						got := opt
						want := &gitlab.UpdateIssueOptions{
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
					Title:      "",
					Message:    "",
					StateEvent: "",
					AssigneeID: 13,
				},
				project:  "group/project",
				id:       12,
				editFunc: func(program, file string) error { return nil },
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.method
			got, err := m.Process()
			if (err != nil) != tt.wantErr {
				t.Errorf("updateMethod.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("updateMethod.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
