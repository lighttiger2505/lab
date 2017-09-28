package main

import (
	"flag"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

type BrowseCommand struct {
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse repository"
}

func (c *BrowseCommand) Help() string {
	return "Usage: lab brewse [option]"
}

func (c *BrowseCommand) Run(args []string) int {
	var debug bool

	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.BoolVar(&debug, "debug", false, "Run as DEBUG mode")
	return 0
}

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"browse": func() (cli.Command, error) {
			return &BrowseCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
