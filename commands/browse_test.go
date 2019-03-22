package commands

import (
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/clipboard"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
	gitlab "github.com/xanzy/go-gitlab"
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

var mockAPIClientFactory = &api.MockAPIClientFactory{
	MockGetBranchClient: func() api.Branch {
		return &api.MockBranchClient{
			MockGetBranch: func(project string, branch string) (*gitlab.Branch, error) {
				return &gitlab.Branch{}, nil
			},
		}
	},
}

func TestBrowseCommandRun(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := BrowseCommand{
		UI:              mockUI,
		RemoteCollecter: &gitutil.MockCollecter{},
		GitClient:       mockGitClient,
		Clipboard:       &clipboard.MockClipboardRW{},
		Opener:          &MockURLOpener{},
		ClientFactory:   mockAPIClientFactory,
	}
	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
}
