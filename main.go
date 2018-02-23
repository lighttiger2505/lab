package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/commands"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/mitchellh/cli"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]

	// ui := &cli.BasicUi{Writer: os.Stdout}
	ui := ui.NewBasicUi()

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

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &commands.BrowseCommand{
				Ui:           ui,
				RemoteFilter: &gitlab.GitlabRemoteFilter{},
				GitClient:    &git.GitClient{},
				Cmd:          cmd.NewBasicCmd(""),
			}, nil
		},
		"issue": func() (cli.Command, error) {
			return &commands.IssueCommand{
				Ui:           ui,
				RemoteFilter: &gitlab.GitlabRemoteFilter{},
				GitClient:    &git.GitClient{},
				LabClient:    &gitlab.LabClient{},
			}, nil
		},
		"add-issue": func() (cli.Command, error) {
			return &commands.AddIssueCommand{Ui: ui}, nil
		},
		"merge-request": func() (cli.Command, error) {
			return &commands.MergeRequestCommand{
				Ui:           ui,
				RemoteFilter: &gitlab.GitlabRemoteFilter{},
				GitClient:    &git.GitClient{},
				LabClient:    &gitlab.LabClient{},
			}, nil
		},
		"add-merge-request": func() (cli.Command, error) {
			return &commands.AddMergeReqeustCommand{Ui: ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
