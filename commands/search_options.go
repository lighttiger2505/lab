package commands

import (
	flags "github.com/jessevdk/go-flags"
)

var searchOptions SearchOptons
var searchParser = flags.NewParser(&searchOptions, flags.Default)

type SearchOptons struct {
	Line       int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	State      string `short:"t" long:"state" default:"all" default-mask:"all" description:"just those that are opened, closed or all"`
	Scope      string `short:"c" long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy    string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Sort       string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
	Opened     bool   `short:"e" long:"opened" description:"search state opened"`
	Closed     bool   `short:"l" long:"closed" description:"search scope closed"`
	CreatedMe  bool   `short:"r" long:"created-me" description:"search scope created-by-me"`
	AssignedMe bool   `short:"g" long:"assigned-me" description:"search scope assigned-to-me"`
	AllProject bool   `short:"a" long:"all-project" description:"search all project"`
}

func (s *SearchOptons) GetState() string {
	if s.Opened {
		return "opened"
	}
	if s.Closed {
		return "closed"
	}
	return s.State
}

func (s *SearchOptons) GetScope() string {
	if s.CreatedMe {
		return "created-by-me"
	}
	if s.AssignedMe {
		return "assigned-to-me"
	}
	return s.Scope
}
