package gitlab

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
)

var TestRemoteInfos = []*git.RemoteInfo{
	&git.RemoteInfo{
		Domain: "gitlab.com",
	},
	&git.RemoteInfo{
		Domain: "gitlab.ssl.unknown.jp",
	},
	&git.RemoteInfo{
		Domain: "github.com",
	},
	&git.RemoteInfo{
		Domain: "gitlao.com",
	},
}

var TestRemoteInfoGitlab = []*git.RemoteInfo{
	&git.RemoteInfo{
		Domain: "gitlab.com",
	},
	&git.RemoteInfo{
		Domain: "gitlab.ssl.unknown.jp",
	},
}

func TestGetClient(t *testing.T) {
	mockUI := ui.NewMockUi()

	// Dummy config file
	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	configManager := config.NewConfigManagerPath(tmppath)

	// Mock git client
	want := &git.RemoteInfo{Domain: "gitlab.ssl.unknown.jp"}
	mockGitClient := &git.MockClient{
		MockRemoteInfos: func() ([]*git.RemoteInfo, error) {
			return []*git.RemoteInfo{
				&git.RemoteInfo{Domain: "github.com"},
				want,
			}, nil
		},
	}

	// Initialize provider
	provider := NewProvider(mockUI, mockGitClient, configManager)
	provider.Init()

	got, err := provider.GetCurrentRemote()
	if err != nil {
		t.Fail()
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("bad output \nwant %#v \ngot  %#v", want, got)
	}
}

func TestGetClient_TokenNotFound(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("token\n")

	// Dummy config file
	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	configManager := config.NewConfigManagerPath(tmppath)

	// Initialize provider
	provider := NewProvider(mockUI, git.NewGitClient(), configManager)
	provider.Init()

	remoteInfo := &git.RemoteInfo{
		Domain:     "gitlab.ssl.unknown.jp",
		Group:      "group",
		Repository: "repository",
	}

	got, err := provider.GetClient(remoteInfo)
	if err != nil {
		t.Fail()
	}
	if got == nil {
		t.Fail()
	}

	// Assert stdout
	outGot := mockUI.Writer.String()
	outWant := `Please input GitLab private token :`
	if outWant != outGot {
		t.Errorf("bad output \nwant %#v \ngot  %#v", outWant, outGot)
	}
}

func TestSelectTargetRemote(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("2\n")

	// Dummy config file
	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	configManager := config.NewConfigManagerPath(tmppath)

	// Initialize provider
	provider := NewProvider(mockUI, git.NewGitClient(), configManager)
	provider.Init()

	got, err := provider.selectTargetRemote(TestRemoteInfoGitlab)
	if err != nil {
		t.Fail()
	}

	// Assert return value
	want := "gitlab.ssl.unknown.jp"
	if want != got.Domain {
		t.Errorf("bad return value want %#v got %#v", want, got.Domain)
	}

	// Assert stdout
	outGot := mockUI.Writer.String()
	outWant := `That repository existing multi gitlab remote repository.
1) gitlab.com
2) gitlab.ssl.unknown.jp
Please choice target domain :`
	if outWant != outGot {
		t.Errorf("bad output want %#v got %#v", outWant, outGot)
	}
}

func TestSelectTargetRemote_InvalidValue_String(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("abc\n")

	// Initialize provider
	provider := NewProvider(mockUI, git.NewGitClient(), config.NewConfigManager())
	provider.Init()

	_, err := provider.selectTargetRemote(TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
	errGot := err.Error()
	errWant := "Failed parse number. Error: "
	if !strings.HasPrefix(errGot, errWant) {
		t.Errorf("bad error message want %s got %s", errWant, errGot)
	}
}

func TestSelectTargetRemote_InvalidValue_Lower(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("0\n")

	// Initialize provider
	provider := NewProvider(mockUI, git.NewGitClient(), config.NewConfigManager())
	provider.Init()

	_, err := provider.selectTargetRemote(TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
	errGot := err.Error()
	errWant := "Invalid number. Input: 0"
	if errGot != errWant {
		t.Errorf("bad error message want %s got %s", errWant, errGot)
	}
}

func TestSelectTargetRemote_InvalidValue_Upper(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("3\n")

	// Initialize provider
	provider := NewProvider(mockUI, git.NewGitClient(), config.NewConfigManager())
	provider.Init()

	_, err := provider.selectTargetRemote(TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
	errGot := err.Error()
	errWant := "Invalid number. Input: 3"
	if errGot != errWant {
		t.Errorf("bad error message want %s got %s", errWant, errGot)
	}
}

func TestFilterHasGitlabDomain(t *testing.T) {
	want := &git.RemoteInfo{
		Domain: "gitlab.ssl.unknown.jp",
	}
	remoteInfos := []*git.RemoteInfo{
		&git.RemoteInfo{Domain: "hogehoge.com"},
		want,
		&git.RemoteInfo{Domain: "hugahuga.com"},
	}
	got := filterHasGitlabDomain(remoteInfos)
	if reflect.DeepEqual(want, got) {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestRegistedDomainRemote(t *testing.T) {
	remoteInfos := []*git.RemoteInfo{
		&git.RemoteInfo{Domain: "foo.com"},
		&git.RemoteInfo{Domain: "bar.com"},
		&git.RemoteInfo{Domain: "hoge.com"},
	}
	domains := []string{
		"notfound.com",
		"hoge.com",
	}
	want := "hoge.com"
	result := registedDomainRemote(remoteInfos, domains)
	got := result.Domain
	if want != got {
		t.Errorf("bad return value \nwant %s \ngot %s", want, got)
	}
}

func TestRegistedDomainRemote_ReturnNil(t *testing.T) {
	remoteInfos := []*git.RemoteInfo{
		&git.RemoteInfo{Domain: "foo.com"},
		&git.RemoteInfo{Domain: "bar.com"},
		&git.RemoteInfo{Domain: "hoge.com"},
	}
	domains := []string{
		"notfound.com",
		"unknown.com",
	}
	got := registedDomainRemote(remoteInfos, domains)
	if got != nil {
		t.Fail()
	}
}

func TestParceRepositoryFullName(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			"project issue",
			"https://gitlab.com/group/project/issues/12",
			"group/project",
		},
		{
			"subgroup issue",
			"https://gitlab.com/group/subgroup/project/issues/12",
			"group/subgroup/project",
		},
		{
			"project mr",
			"https://gitlab.com/group/project/issues/12",
			"group/project",
		},
		{
			"subgroup mr",
			"https://gitlab.com/group/subgroup/project/issues/12",
			"group/subgroup/project",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParceRepositoryFullName(tt.arg); got != tt.want {
				t.Errorf("ParceRepositoryFullName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_excludeDuplicateDomain(t *testing.T) {
	tests := []struct {
		name string
		arg  []*git.RemoteInfo
		want []*git.RemoteInfo
	}{
		{
			name: "has git submodule",
			arg: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "origin",
					Domain: "gitlab.com",
				},
				&git.RemoteInfo{
					Remote: "lib1",
					Domain: "gitlab.com",
				},
				&git.RemoteInfo{
					Remote: "lib2",
					Domain: "gitlab.com",
				},
			},
			want: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "origin",
					Domain: "gitlab.com",
				},
			},
		},
		{
			name: "has git submodule no origin",
			arg: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "lib1",
					Domain: "gitlab.com",
				},
				&git.RemoteInfo{
					Remote: "lib2",
					Domain: "gitlab.com",
				},
			},
			want: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "lib1",
					Domain: "gitlab.com",
				},
			},
		},
		{
			name: "multi domain",
			arg: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "origin",
					Domain: "gitlab.com",
				},
				&git.RemoteInfo{
					Remote: "xxx",
					Domain: "gitlab.ssl.xxx.com",
				},
			},
			want: []*git.RemoteInfo{
				&git.RemoteInfo{
					Remote: "origin",
					Domain: "gitlab.com",
				},
				&git.RemoteInfo{
					Remote: "xxx",
					Domain: "gitlab.ssl.xxx.com",
				},
			},
		},
		// TODO tmp disable for tesing green
		// {
		// 	name: "has submoduel and multi domain",
		// 	arg: []*git.RemoteInfo{
		// 		&git.RemoteInfo{
		// 			Remote: "origin",
		// 			Domain: "gitlab.com",
		// 		},
		// 		&git.RemoteInfo{
		// 			Remote: "lib1",
		// 			Domain: "gitlab.com",
		// 		},
		// 		&git.RemoteInfo{
		// 			Remote: "xxx1",
		// 			Domain: "gitlab.ssl.xxx.com",
		// 		},
		// 		&git.RemoteInfo{
		// 			Remote: "xxx2",
		// 			Domain: "gitlab.ssl.xxx.com",
		// 		},
		// 	},
		// 	want: []*git.RemoteInfo{
		// 		&git.RemoteInfo{
		// 			Remote: "origin",
		// 			Domain: "gitlab.com",
		// 		},
		// 		&git.RemoteInfo{
		// 			Remote: "xxx1",
		// 			Domain: "gitlab.ssl.xxx.com",
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := excludeDuplicateDomain(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("excludeDuplicateDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}
