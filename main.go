package main

import (
	"log"
	"os"

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

	ui := &cli.BasicUi{Writer: os.Stdout}

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &BrowseCommand{Ui: ui}, nil
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
		log.Println(err)
	}

	os.Exit(exitStatus)
}
