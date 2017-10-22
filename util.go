package main

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

func SearchBrowserLauncher(goos string) (browser string) {
	switch goos {
	case "darwin":
		browser = "open"
	case "windows":
		browser = "cmd /c start"
	default:
		candidates := []string{
			"xdg-open",
			"cygstart",
			"x-www-browser",
			"firefox",
			"opera",
			"mozilla",
			"netscape",
		}
		for _, b := range candidates {
			path, err := exec.LookPath(b)
			if err == nil {
				browser = path
				break
			}
		}
	}
	return browser
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
			Type: "Issue",
			No:   number,
		}
	} else if strings.HasPrefix(arg, "!") {
		number, err := strconv.Atoi(strings.TrimPrefix(arg, "!"))
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
