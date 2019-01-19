package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/lighttiger2505/lab/commands"
	configcmd "github.com/lighttiger2505/lab/commands/config"
	"github.com/lighttiger2505/lab/commands/issue"
	"github.com/lighttiger2505/lab/commands/mr"
	"github.com/lighttiger2505/lab/commands/pipeline"
	"github.com/lighttiger2505/lab/commands/runner"
	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/config"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
	"github.com/mitchellh/cli"
)

var (
	version  string
	revision string
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

func main() {
	os.Exit(realMain(os.Stdout, version, revision))
}

func realMain(writer io.Writer, ver, rev string) int {
	c := cli.NewCLI("lab", fmt.Sprintf("ver: %s rev: %s", ver, rev))
	c.Args = os.Args[1:]
	c.HelpWriter = writer

	// Determine where logs should go in general (requested by the user)
	logWriter, err := logOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't setup log output: %s", err)
	}
	if logWriter == nil {
		logWriter = ioutil.Discard
	}

	// Disable logging here
	log.SetOutput(ioutil.Discard)

	ui := ui.NewBasicUi()
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config, %s", err)
	}
	remoteCollecter := gitutil.NewRemoteCollecter(ui, cfg, git.NewGitClient())

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &commands.BrowseCommand{
				Ui:              ui,
				RemoteCollecter: remoteCollecter,
				GitClient:       &git.GitClient{},
				Opener:          &browse.Browser{},
			}, nil
		},
		"issue": func() (cli.Command, error) {
			return &issue.IssueCommand{
				Ui:              ui,
				RemoteCollecter: remoteCollecter,
				MethodFactory:   &issue.IssueMethodFactory{},
			}, nil
		},
		"merge-request": func() (cli.Command, error) {
			return &mr.MergeRequestCommand{
				Ui:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"mr": func() (cli.Command, error) {
			return &mr.MergeRequestCommand{
				Ui:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"project": func() (cli.Command, error) {
			return &commands.ProjectCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"pipeline": func() (cli.Command, error) {
			return &pipeline.PipelineCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				MethodFactory:   &pipeline.PipelineMethodFacotry{},
			}, nil
		},
		"job": func() (cli.Command, error) {
			return &commands.JobCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"lint": func() (cli.Command, error) {
			return &commands.LintCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"user": func() (cli.Command, error) {
			return &commands.UserCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"project-variable": func() (cli.Command, error) {
			return &commands.ProjectVariableCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"issue-template": func() (cli.Command, error) {
			return &commands.IssueTemplateCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"merge-request-template": func() (cli.Command, error) {
			return &commands.MergeRequestTemplateCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"runner": func() (cli.Command, error) {
			return &runner.RunnerCommand{
				UI:              ui,
				RemoteCollecter: remoteCollecter,
				ClientFactory:   &lab.GitlabClientFactory{},
			}, nil
		},
		"config": func() (cli.Command, error) {
			return &configcmd.ConfigCommand{
				UI:     ui,
				Config: cfg,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}
	return exitStatus
}
