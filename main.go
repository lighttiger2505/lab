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

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &BrowseCommand{}, nil
		},
		"issue": func() (cli.Command, error) {
			return &IssueCommand{Ui: &cli.BasicUi{Writer: os.Stdout}}, nil
		},
		"merge-request": func() (cli.Command, error) {
			return &MergeRequestCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
