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
	AllProject bool   `short:"a" long:"all-project" description:"search all project"`
}
