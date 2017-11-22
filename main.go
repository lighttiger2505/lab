package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/lighttiger2505/lab/commands"
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
			return &commands.BrowseCommand{Ui: ui}, nil
		},
		"issue": func() (cli.Command, error) {
			return &IssueCommand{Ui: ui}, nil
		},
		"merge-request": func() (cli.Command, error) {
			return &MergeRequestCommand{Ui: ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error(err.Error())
	}

	os.Exit(exitStatus)
}
