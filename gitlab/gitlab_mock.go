package gitlab

import (
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
)

type MockRemoteFilter struct {
}

func NewRemoteFilter() *MockRemoteFilter {
	return &MockRemoteFilter{}
}

func (g *MockRemoteFilter) Collect() error {
	return nil
}

func (g *MockRemoteFilter) Filter(ui ui.Ui, conf *config.Config) (*git.RemoteInfo, error) {
	gitlabRemote := &git.RemoteInfo{
		Domain:     "domain",
		NameSpace:  "namespace",
		Repository: "project",
	}
	return gitlabRemote, nil
}
