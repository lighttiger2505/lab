package gitutil

import (
	"reflect"
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/config"
)

func Test_filterHasGitlabDomain(t *testing.T) {

	gitlabRemoteInfo := &git.RemoteInfo{

		Domain: "gitlab.com",
	}
	githubRemoteInfo := &git.RemoteInfo{

		Domain: "github.com",
	}
	type args struct {
		remoteInfos []*git.RemoteInfo
		cfg         *config.Config
	}
	tests := []struct {
		name string
		args args
		want []*git.RemoteInfo
	}{
		{
			name: "find the gitlab domain",
			args: args{
				remoteInfos: []*git.RemoteInfo{
					gitlabRemoteInfo,
					githubRemoteInfo,
				},
				cfg: &config.Config{
					Profiles: map[string]config.Profile{},
				},
			},
			want: []*git.RemoteInfo{gitlabRemoteInfo},
		},
		{
			name: "find the gitlab domain specified in the config file",
			args: args{
				remoteInfos: []*git.RemoteInfo{
					gitlabRemoteInfo,
					githubRemoteInfo,
				},
				cfg: &config.Config{
					Profiles: map[string]config.Profile{
						"gitlab.com": {},
					},
				},
			},
			want: []*git.RemoteInfo{gitlabRemoteInfo},
		},
		{
			name: "find the other domain specified in the config file",
			args: args{
				remoteInfos: []*git.RemoteInfo{
					gitlabRemoteInfo,
					githubRemoteInfo,
				},
				cfg: &config.Config{
					Profiles: map[string]config.Profile{
						"gitlab.com": {},
						"github.com": {},
					},
				},
			},
			want: []*git.RemoteInfo{gitlabRemoteInfo, githubRemoteInfo},
		},
		{
			name: "not found",
			args: args{
				remoteInfos: []*git.RemoteInfo{
					githubRemoteInfo,
				},
				cfg: &config.Config{
					Profiles: map[string]config.Profile{
						"gitlab.com": {},
					},
				},
			},
			want: []*git.RemoteInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterHasGitlabDomain(tt.args.remoteInfos, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterHasGitlabDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
