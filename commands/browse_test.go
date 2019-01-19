package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
)

var mockGitClient = &git.MockClient{
	MockRemoteInfos: func() ([]*git.RemoteInfo, error) {
		return []*git.RemoteInfo{
			&git.RemoteInfo{
				Domain:     "gitlab.ssl.domain1.jp",
				Group:      "group",
				Repository: "project",
			},
		}, nil
	},
	MockCurrentRemoteBranch: func() (string, error) {
		return "currentBranch", nil
	},
}

type MockURLOpener struct{}

func (m *MockURLOpener) Open(url string) error {
	return nil
}

func TestBrowseCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := BrowseCommand{
		UI:              mockUI,
		RemoteCollecter: &gitutil.MockCollecter{},
		GitClient:       mockGitClient,
		Opener:          &MockURLOpener{},
	}
	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
}
