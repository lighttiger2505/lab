package main

import (
	"reflect"
	"testing"
)

type newGitRemoteTest struct {
	url       string
	gitRemote *GitRemote
}

var newRemoteTests = []newGitRemoteTest{
	{
		url: "ssh://git@gitlab.ssl.domain.jp/namespace/repository.git",
		gitRemote: &GitRemote{
			Domain:     "gitlab.ssl.domain.jp",
			User:       "namespace",
			Repository: "repository",
		},
	},
	{
		url: "git@gitlab.ssl.domain.jp:namespace/repository.git",
		gitRemote: &GitRemote{
			Domain:     "gitlab.ssl.domain.jp",
			User:       "namespace",
			Repository: "repository",
		},
	},
	{
		url: "https://gitlab.ssl.domain.jp/namespace/repository",
		gitRemote: &GitRemote{
			Domain:     "gitlab.ssl.domain.jp",
			User:       "namespace",
			Repository: "repository",
		},
	},
}

func TestNewGitRemote(t *testing.T) {
	for i, test := range newRemoteTests {
		got, _ := NewRemoteUrl(test.url)
		if !reflect.DeepEqual(test.gitRemote, got) {
			t.Errorf("#%d: bad return value want %#v got %#v", i, test.gitRemote, got)
		}
	}
}
