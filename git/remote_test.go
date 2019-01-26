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
		url: "ssh://git@gitlab.ssl.domain.jp/group/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:group/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/group/repository",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			Repository: "repository",
		},
	},
	{
		url: "ssh://git@gitlab.ssl.domain.jp/group/subgroup/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:group/subgroup/repository.git",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/group/subgroup/repository",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			SubGroup:   "subgroup",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/group/subgroup1/subgroup2/subgroup3/repository",
		remoteInfo: &RemoteInfo{
			Remote:     "origin",
			Domain:     "gitlab.ssl.domain.jp",
			Group:      "group",
			SubGroup:   "subgroup1/subgroup2/subgroup3",
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
	Group:      "group",
	Repository: "Repository",
}

var testRemoteInfoWithSubgroup = &RemoteInfo{
	Domain:     "gitlab.ssl.domain.jp",
	Group:      "group",
	SubGroup:   "subgroup",
	Repository: "Repository",
}

func TestRepositoryFullName(t *testing.T) {
	got := testRemoteInfoWithSubgroup.RepositoryUrl()
	want := "https://gitlab.ssl.domain.jp/group/subgroup/Repository"
	if want != got {
		t.Errorf("bad return value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestRepositoryUrl(t *testing.T) {
	got := testRemoteInfo.RepositoryUrl()
	want := "https://gitlab.ssl.domain.jp/group/Repository"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestBranchUrl(t *testing.T) {
	got := testRemoteInfo.BranchUrl("Branch")
	want := "https://gitlab.ssl.domain.jp/group/Repository/tree/Branch"
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
