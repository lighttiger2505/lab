package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GitRemotes() ([]GitRemote, error) {
	// Get remote repositorys
	remotes := gitOutputs("git", []string{"remote"})

	// Remote repository is not registered
	if len(remotes) == 0 {
		return nil, errors.New("No remote setting in this repository")
	}

	var gitRemotes []GitRemote
	for _, remote := range remotes {
		url := gitOutput("git", []string{"remote", "get-url", remote})

		gitRemote, err := NewRemoteUrl(url)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed serialize remote url. %s", url))
		}

		gitRemotes = append(gitRemotes, *gitRemote)
	}

	return gitRemotes, nil
}

type BrowseArg struct {
	Type string
	No   int
}

func NewBrowseArg(arg string) (*BrowseArg, error) {
	var browseArg BrowseArg
	if strings.HasPrefix(arg, "#") {
		number, err := strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil {
			return nil, errors.New("Invalid number")
		}
		browseArg = BrowseArg{
			Type: "MergeRequest",
			No:   number,
		}
	} else {
		return nil, errors.New("Invalid args")
	}
	return &browseArg, nil
}
