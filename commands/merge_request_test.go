package commands

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlabc "github.com/xanzy/go-gitlab"
)

var mergeRequests []*gitlabc.MergeRequest = []*gitlabc.MergeRequest{
	&gitlabc.MergeRequest{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlabc.MergeRequest{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockLabMergeRequestClient *gitlab.MockLabClient = &gitlab.MockLabClient{
	MockMergeRequest: func(baseurl, token string, opt *gitlabc.ListMergeRequestsOptions) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockProjectMergeRequest: func(baseurl, token string, opt *gitlabc.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlabc.MergeRequest, error) {
		return mergeRequests, nil
	},
}

func TestMergeRequestCommandRun(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")

	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	conf := config.NewConfigManagerPath(tmppath)

	c := MergeRequestCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		LabClient:    mockLabMergeRequestClient,
		Config:       conf,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "!12  Title12\n!13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_AllProjectOption(t *testing.T) {
	mockUi := ui.NewMockUi()
	mockUi.Reader = bytes.NewBufferString("token\n")

	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Write([]byte(config.ConfigDataTest))
	f.Close()
	defer os.Remove(tmppath)
	conf := config.NewConfigManagerPath(tmppath)

	c := MergeRequestCommand{
		Ui:           mockUi,
		RemoteFilter: gitlab.NewRemoteFilter(),
		LabClient:    mockLabMergeRequestClient,
		Config:       conf,
	}

	args := []string{"-a"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUi.ErrorWriter.String())
	}

	got := mockUi.Writer.String()
	want := "!12  namespace/repo12  Title12\n!13  namespace/repo13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}
