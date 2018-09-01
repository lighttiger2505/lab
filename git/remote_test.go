package git

import (
	"reflect"
	"testing"
)

type newGitRemoteTest struct {
	url        string
	remoteInfo *RemoteInfo
}

var newRemoteTests = []newGitRemoteTest{
	{
		url: "ssh://git@gitlab.ssl.domain.jp/namespace/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:namespace/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/namespace/repository",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
	{
		url: "ssh://git@gitlab.ssl.domain.jp/namespace/subgroup/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:namespace/subgroup/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/namespace/subgroup/repository",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
}

func TestNewGitRemote(t *testing.T) {
	for i, test := range newRemoteTests {
		got := NewRemoteInfo("origin", test.url)
		if !reflect.DeepEqual(test.remoteInfo, got) {
			t.Errorf("#%d: bad return value \nwant %#v \ngot  %#v", i, test.remoteInfo, got)
		}
	}
}

var testRemoteInfo = &RemoteInfo{
	Domain:     "gitlab.ssl.domain.jp",
	NameSpace:  "Namespace",
	Repository: "Repository",
}

var testRemoteInfoWithSubgroup = &RemoteInfo{
	Domain:     "gitlab.ssl.domain.jp",
	NameSpace:  "Namespace",
	SubGroup:   "subgroup",
	Repository: "Repository",
}

func TestRepositoryFullName(t *testing.T) {
	got := testRemoteInfoWithSubgroup.RepositoryUrl()
	want := "https://gitlab.ssl.domain.jp/Namespace/subgroup/Repository"
	if want != got {
		t.Errorf("bad return value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestRepositoryUrl(t *testing.T) {
	got := testRemoteInfo.RepositoryUrl()
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestBranchUrl(t *testing.T) {
	got := testRemoteInfo.BranchUrl("Branch")
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/tree/Branch"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestIssueUrl(t *testing.T) {
	got := testRemoteInfo.IssueUrl()
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/issues"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestIssueDetailUrl(t *testing.T) {
	got := testRemoteInfo.IssueDetailUrl(12)
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/issues/12"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestMergeRequestUrl(t *testing.T) {
	got := testRemoteInfo.MergeRequestUrl()
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/merge_requests"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestMergeRequestDetailUrl(t *testing.T) {
	got := testRemoteInfo.MergeRequestDetailUrl(12)
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/merge_requests/12"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestPipeLineUrl(t *testing.T) {
	got := testRemoteInfo.PipeLineUrl()
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/pipelines"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestPipeLineDetailUrl(t *testing.T) {
	got := testRemoteInfo.PipeLineDetailUrl(12)
	want := "https://gitlab.ssl.domain.jp/Namespace/Repository/pipelines/12"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestBaseUrl(t *testing.T) {
	got := testRemoteInfo.BaseUrl()
	want := "https://gitlab.ssl.domain.jp"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestApiUrl(t *testing.T) {
	got := testRemoteInfo.ApiUrl()
	want := "https://gitlab.ssl.domain.jp/api/v4"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}
