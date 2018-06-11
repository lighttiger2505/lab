package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

var mockGitClient = &git.MockClient{
	MockRemoteInfos: func() ([]*git.RemoteInfo, error) {
		return []*git.RemoteInfo{
			&git.RemoteInfo{
				Domain:     "gitlab.ssl.domain1.jp",
				NameSpace:  "namespace",
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

	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	configManager := config.NewConfigManagerPath(tmppath)

	// Initialize provider
	provider := gitlab.NewProvider(mockUI, mockGitClient, configManager)
	provider.Init()

	c := BrowseCommand{
		Ui:        mockUI,
		Provider:  provider,
		GitClient: mockGitClient,
		Opener:    &MockURLOpener{},
	}
	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
}
