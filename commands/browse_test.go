package commands

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
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

func TestBrowseCommandRun_Path(t *testing.T) {
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
	args := []string{"--path=hoge"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
	got := mockUI.Writer.String()
	want := "browse path\n"
	if got != want {
		t.Fatalf("Invalid output. \n got:%q\nwant:%q", got, want)
	}
}

func TestBrowseCommandRun_Path2(t *testing.T) {
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
	args := []string{"--path"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}
	got := mockUI.Writer.String()
	want := "browse path\n"
	if got != want {
		t.Fatalf("Invalid output. \n got:%q\nwant:%q", got, want)
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

var getUrlByUserSpecificTests = []getUrlByUserSpecificTest{
	{gitlabRemote: gitlabRemoteTest, args: []string{"#12", "#13"}, domain: "specific", url: "https://domain/namespace/project/issues/12", err: nil},
	{gitlabRemote: gitlabRemoteTest, args: []string{}, domain: "specific", url: "https://domain/namespace/project", err: nil},
	{gitlabRemote: nil, args: []string{}, domain: "specific", url: "https://specific", err: nil},
}

func TestGetUrlByUserSpecific(t *testing.T) {
	for i, test := range getUrlByUserSpecificTests {
		url, err := getUrlByUserSpecific(test.gitlabRemote, test.args, test.domain)
		if test.url != url || !reflect.DeepEqual(test.err, err) {
			t.Errorf("#%d: bad return value \nwant %#v %#v \ngot  %#v %#v", i, test.url, test.err, url, err)
		}
	}
}

type getUrlByRemoteTest struct {
	gitlabRemote *git.RemoteInfo
	args         []string
	branch       string
	url          string
	err          error
}

var getUrlByRemoteTests = []getUrlByRemoteTest{
	{gitlabRemote: gitlabRemoteTest, args: []string{"#12", "#13"}, branch: "", url: "https://domain/namespace/project/issues/12", err: nil},
	{gitlabRemote: gitlabRemoteTest, args: []string{}, branch: "master", url: "https://domain/namespace/project", err: nil},
	{gitlabRemote: gitlabRemoteTest, args: []string{}, branch: "develop", url: "https://domain/namespace/project/tree/develop", err: nil},
}

func TestGetUrlByRemote(t *testing.T) {
	for i, test := range getUrlByRemoteTests {
		url, err := getUrlByRemote(test.gitlabRemote, test.args, test.branch)
		if test.url != url || !reflect.DeepEqual(test.err, err) {
			t.Errorf("#%d: bad return value \nwant %#v %#v \ngot  %#v %#v", i, test.url, test.err, url, err)
		}
	}
}

type makeGitlabResourceUrlTest struct {
	gitlabRemote *git.RemoteInfo
	browseType   BrowseType
	number       int
	url          string
}

var makeGitlabResourceUrlTests = []makeGitlabResourceUrlTest{
	{gitlabRemote: gitlabRemoteTest, browseType: Issue, number: 12, url: gitlabRemoteTest.IssueDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: MergeRequest, number: 12, url: gitlabRemoteTest.MergeRequestDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: PipeLine, number: 12, url: gitlabRemoteTest.PipeLineDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: Issue, number: 0, url: gitlabRemoteTest.IssueUrl()},
	{gitlabRemote: gitlabRemoteTest, browseType: MergeRequest, number: 0, url: gitlabRemoteTest.MergeRequestUrl()},
	{gitlabRemote: gitlabRemoteTest, browseType: PipeLine, number: 0, url: gitlabRemoteTest.PipeLineUrl()},
}

func TestBrowseUrl(t *testing.T) {
	for i, test := range makeGitlabResourceUrlTests {
		url := makeGitlabResourceUrl(test.gitlabRemote, test.browseType, test.number)
		if test.url != url {
			t.Errorf("#%d: bad return value want %#v got %#v", i, test.url, url)
		}
	}
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

type splitPrefixAndNumberTest struct {
	arg        string
	browseType BrowseType
	number     int
	err        error
}

var splitPrefixAndNumberTests = []splitPrefixAndNumberTest{
	{arg: "#12", browseType: Issue, number: 12, err: nil},
	{arg: "I12", browseType: Issue, number: 12, err: nil},
	{arg: "i12", browseType: Issue, number: 12, err: nil},
	{arg: "!12", browseType: MergeRequest, number: 12, err: nil},
	{arg: "M12", browseType: MergeRequest, number: 12, err: nil},
	{arg: "m12", browseType: MergeRequest, number: 12, err: nil},
	{arg: "P12", browseType: PipeLine, number: 12, err: nil},
	{arg: "p12", browseType: PipeLine, number: 12, err: nil},
	{arg: "I", browseType: Issue, number: 0, err: nil},
	{arg: "M", browseType: MergeRequest, number: 0, err: nil},
	{arg: "P", browseType: PipeLine, number: 0, err: nil},
	{arg: "Iunknown", browseType: 0, number: 0, err: errors.New("Invalid browse number. \"unknown\"")},
	{arg: "Unknown", browseType: 0, number: 0, err: errors.New("Invalid arg. Unknown")},
}

func TestSplitPrefixAndNumber(t *testing.T) {
	for i, test := range splitPrefixAndNumberTests {
		browseType, number, err := splitPrefixAndNumber(test.arg)
		if test.browseType != browseType || test.number != number || !reflect.DeepEqual(test.err, err) {
			t.Errorf(
				"#%d: bad return value \nwant %#v %#v %#v \ngot  %#v %#v %#v",
				i,
				test.browseType,
				test.number,
				test.err,
				browseType,
				number,
				err,
			)
		}
	}
}
