package commands

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lighttiger2505/lab/git"
	// "github.com/lighttiger2505/lab/ui"
)

// func TestBrowseCommandRun(t *testing.T) {
// 	ui := ui.NewMockUi()
// 	cmd := BrowseCommand{Ui: ui}
// 	args := []string{}
// 	want := 0
// 	got := cmd.Run(args)
// 	if want != got {
// 		t.Errorf("bad return value want %#v got %#v", want, got)
// 	}
// }

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
