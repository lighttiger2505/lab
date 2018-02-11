package gitlab

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/ui"
)

var TestRemoteInfos = []git.RemoteInfo{
	git.RemoteInfo{
		Domain: "gitlab.com",
	},
	git.RemoteInfo{
		Domain: "gitlab.ssl.unknown.jp",
	},
	git.RemoteInfo{
		Domain: "github.com",
	},
	git.RemoteInfo{
		Domain: "gitlao.com",
	},
}

var TestRemoteInfoGitlab = []git.RemoteInfo{
	git.RemoteInfo{
		Domain: "gitlab.com",
	},
	git.RemoteInfo{
		Domain: "gitlab.ssl.unknown.jp",
	},
}

func TestFilterHasGitlabDomain(t *testing.T) {
	got := filterHasGitlabDomain(TestRemoteInfos)
	want := []git.RemoteInfo{
		git.RemoteInfo{
			Domain: "gitlab.com",
		},
		git.RemoteInfo{
			Domain: "gitlab.ssl.unknown.jp",
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestHasPriorityRemote(t *testing.T) {
	domains := []string{
		"gitlab.ssl.unknown.jp",
		"gitlab.com",
	}

	got := hasPriorityRemote(TestRemoteInfoGitlab, domains).Domain
	want := "gitlab.ssl.unknown.jp"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestHasPriorityRemote_NotFound(t *testing.T) {
	domains := []string{
		"gitlao.ssl.unknown.jp",
		"gitlao.com",
	}

	got := hasPriorityRemote(TestRemoteInfoGitlab, domains)
	if nil != got {
		t.Errorf("bad return value want %#v got %#v", nil, got)
	}
}

func TestInputUseRemote(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("2\n")
	got, err := inputUseRemote(mockUi, TestRemoteInfoGitlab)
	if err != nil {
		t.Fail()
	}

	want := "gitlab.ssl.unknown.jp"
	if want != got.Domain {
		t.Errorf("bad return value want %#v got %#v", want, got.Domain)
	}

	outGot := mockUi.Writer.String()
	outWant := `That repository existing multi gitlab remote repository.
1) gitlab.com
2) gitlab.ssl.unknown.jp
Please choice target domain :`
	if outWant != outGot {
		t.Errorf("bad output want %#v got %#v", outWant, outGot)
	}
}

func TestInputUseRemote_InvalidValue_String(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("abc\n")
	_, err := inputUseRemote(mockUI, TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
}

func TestInputUseRemote_InvalidValue_Lower(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("0\n")
	_, err := inputUseRemote(mockUI, TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
}

func TestInputUseRemote_InvalidValue_Upper(t *testing.T) {
	mockUI := ui.NewMockUi()
	mockUI.Reader = bytes.NewBufferString("3\n")
	_, err := inputUseRemote(mockUI, TestRemoteInfoGitlab)
	if err == nil {
		t.Fail()
	}
}

func TestParceRepositoryFullName(t *testing.T) {
	got := ParceRepositoryFullName("https://gitlab.ssl.iridge.jp/proj/repo/issues/12")
	want := "proj/repo"
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}
