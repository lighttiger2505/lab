package commands

import (
	"errors"
	"github.com/lighttiger2505/lab/git"
	"reflect"
	"testing"
)

func TestBrowseCommandRun(t *testing.T) {

}

type browseUrlTest struct {
	gitlabRemote *git.RemoteInfo
	browseType   BrowseType
	number       int
	url          string
}

var gitlabRemoteTest = &git.RemoteInfo{
	Domain:     "domain",
	NameSpace:  "namespace",
	Repository: "project",
}

var browseUrlTests = []browseUrlTest{
	{gitlabRemote: gitlabRemoteTest, browseType: Issue, number: 12, url: gitlabRemoteTest.IssueDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: MergeRequest, number: 12, url: gitlabRemoteTest.MergeRequestDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: PipeLine, number: 12, url: gitlabRemoteTest.PipeLineDetailUrl(12)},
	{gitlabRemote: gitlabRemoteTest, browseType: Issue, number: 0, url: gitlabRemoteTest.IssueUrl()},
	{gitlabRemote: gitlabRemoteTest, browseType: MergeRequest, number: 0, url: gitlabRemoteTest.MergeRequestUrl()},
	{gitlabRemote: gitlabRemoteTest, browseType: PipeLine, number: 0, url: gitlabRemoteTest.PipeLineUrl()},
}

func TestBrowseUrl(t *testing.T) {
	for i, test := range browseUrlTests {
		url := browseUrl(test.gitlabRemote, test.browseType, test.number)
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
