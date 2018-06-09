package commands

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/lighttiger2505/lab/cmd"
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
		Cmd:       cmd.NewMockCmd("browse"),
	}
	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
}

var gitlabRemoteTest = &git.RemoteInfo{
	Domain:     "domain",
	NameSpace:  "namespace",
	Repository: "project",
}

type getUrlByUserSpecificTest struct {
	gitlabRemote *git.RemoteInfo
	args         []string
	domain       string
	url          string
	err          error
}

type getUrlByRemoteTest struct {
	gitlabRemote *git.RemoteInfo
	args         []string
	branch       string
	url          string
	err          error
}

type searchBrowserLauncherTest struct {
	goos    string
	browser string
}

var searchBrowserLauncherTests = []searchBrowserLauncherTest{
	{goos: "darwin", browser: "open"},
	{goos: "windows", browser: "cmd /c start"},
}

func TestSearchBrowserLauncher(t *testing.T) {
	for i, test := range searchBrowserLauncherTests {
		browser := searchBrowserLauncher(test.goos)
		if test.browser != browser {
			t.Errorf("#%d: bad return value \nwant %#v \ngot  %#v", i, test.browser, browser)
		}
	}
}
