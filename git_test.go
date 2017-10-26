package main

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
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:namespace/repository.git",
		remoteInfo: &RemoteInfo{
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/namespace/repository",
		remoteInfo: &RemoteInfo{
			Domain:     "gitlab.ssl.domain.jp",
			NameSpace:  "namespace",
			Repository: "repository",
		},
	},
}

func TestNewGitRemote(t *testing.T) {
	for i, test := range newRemoteTests {
		got := NewRemoteInfo(test.url)
		if !reflect.DeepEqual(test.remoteInfo, got) {
			t.Errorf("#%d: bad return value want %#v got %#v", i, test.remoteInfo, got)
		}
	}
}
