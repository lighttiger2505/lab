package commands

import (
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

var globalOpt GlobalOpt
var globalParser = flags.NewParser(&globalOpt, flags.Default)

type GlobalOpt struct {
	Repository string `short:"p" long:"repository" description:"target specific repository"`
}

var searchOptions SearchOpt
var searchParser = flags.NewParser(&searchOptions, flags.Default)

type SearchOpt struct {
	Line          int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	State         string `short:"t" long:"state" default:"all" default-mask:"all" description:"just those that are opened, closed or all"`
	Scope         string `short:"c" long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy       string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Sort          string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
	Opened        bool   `short:"e" long:"opened" description:"search state opened"`
	Closed        bool   `short:"l" long:"closed" description:"search scope closed"`
	CreatedMe     bool   `short:"r" long:"created-me" description:"search scope created-by-me"`
	AssignedMe    bool   `short:"g" long:"assigned-me" description:"search scope assigned-to-me"`
	AllRepository bool   `short:"a" long:"all-repository" description:"search target all repository"`
}

func (s *SearchOpt) GetState() string {
	if s.Opened {
		return "opened"
	}
	if s.Closed {
		return "closed"
	}
	return s.State
}

func (s *SearchOpt) GetScope() string {
	if s.CreatedMe {
		return "created-by-me"
	}
	if s.AssignedMe {
		return "assigned-to-me"
	}
	return s.Scope
}

func (g *GlobalOpt) ValidRepository() (string, string, error) {
	value := g.Repository
	splited := strings.Split(value, "/")
	if value != "" && len(splited) != 2 {
		return "", "", fmt.Errorf("Invalid repository \"%s\". Assumed input style is \"Namespace/Project\".", value)
	}
	return splited[0], splited[1], nil
}
