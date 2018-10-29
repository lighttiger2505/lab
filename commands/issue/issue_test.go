package issue

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

var mockProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			Group:      "group",
			Repository: "repository",
		}, nil
	},
}
var mockAPIClientFactory = &lab.MockAPIClientFactory{}
var mockMethodFactory = &MockMethodFactory{}

func TestIssueCommand_Run(t *testing.T) {
	type fields struct {
		Provider      lab.Provider
		MethodFactory MethodFactory
	}
	tests := []struct {
		name     string
		fields   fields
		args     []string
		wantCode int
		wantOut  string
		wantErr  string
	}{
		{
			name: "normal",
			fields: fields{
				Provider:      mockProvider,
				MethodFactory: mockMethodFactory,
			},
			args:     []string{},
			wantCode: 0,
			wantOut:  "result\n",
			wantErr:  "",
		},
		{
			name: "unknown flag",
			fields: fields{
				Provider:      mockProvider,
				MethodFactory: mockMethodFactory,
			},
			args:     []string{"--hogehoge"},
			wantCode: 1,
			wantOut:  "",
			wantErr:  "unknown flag `hogehoge'\n",
		},
		{
			name: "nomal args",
			fields: fields{
				Provider:      mockProvider,
				MethodFactory: mockMethodFactory,
			},
			args:     []string{"12"},
			wantCode: 0,
			wantOut:  "result\n",
			wantErr:  "",
		},
		{
			name: "multipul args",
			fields: fields{
				Provider:      mockProvider,
				MethodFactory: mockMethodFactory,
			},
			args:     []string{"12", "13"},
			wantCode: 0,
			wantOut:  "result\n",
			wantErr:  "",
		},
		{
			name: "invalid args",
			fields: fields{
				Provider:      mockProvider,
				MethodFactory: mockMethodFactory,
			},
			args:     []string{"aa"},
			wantCode: 1,
			wantOut:  "",
			wantErr:  "Invalid args, please intput issue IID.\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUI := ui.NewMockUi()
			c := &IssueCommand{
				Ui:            mockUI,
				Provider:      tt.fields.Provider,
				MethodFactory: tt.fields.MethodFactory,
			}
			if got := c.Run(tt.args); got != tt.wantCode {
				t.Errorf("failed issue command run.\ngot: %v\nwant:%v", got, tt.wantCode)
			}
			if got := mockUI.Writer.String(); got != tt.wantOut {
				t.Errorf("unmatch want stdout.\ngot: %#v\nwant:%#v", got, tt.wantOut)
			}
			if got := mockUI.ErrorWriter.String(); got != tt.wantErr {
				t.Errorf("unmatch want stderr.\ngot: %#v\nwant:%#v", got, tt.wantErr)
			}
		})
	}
}
