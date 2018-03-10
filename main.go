package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/commands"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
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
	c := cli.NewCLI("app", fmt.Sprintf("ver: %s rev: %s", ver, rev))
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
	configManager := config.NewConfigManager()
	provider := gitlab.NewProvider(ui, git.NewGitClient(), configManager)

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &commands.BrowseCommand{
				Ui:        ui,
				Provider:  provider,
				GitClient: &git.GitClient{},
				Cmd:       cmd.NewBasicCmd(""),
			}, nil
		},
		"issue": func() (cli.Command, error) {
			return &commands.IssueCommand{
				Ui:       ui,
				Provider: provider,
			}, nil
		},
		"add-issue": func() (cli.Command, error) {
			return &commands.AddIssueCommand{
				Ui:       ui,
				Provider: provider,
			}, nil
		},
		"merge-request": func() (cli.Command, error) {
			return &commands.MergeRequestCommand{
				Ui:       ui,
				Provider: provider,
			}, nil
		},
		"add-merge-request": func() (cli.Command, error) {
			return &commands.AddMergeReqeustCommand{
				Ui:       ui,
				Provider: provider,
			}, nil
		},
		"project": func() (cli.Command, error) {
			return &commands.ProjectCommand{
				UI:       ui,
				Provider: provider,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}
	return exitStatus
}
