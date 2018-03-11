package commands

import (
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type OptionValidator interface {
	IsValid(error)
}

type GlobalOption struct {
	Repository string `short:"p" long:"repository" description:"target specific repository"`
}

func newGlobalOption() *GlobalOption {
	global := flags.NewNamedParser("lab", flags.Default)
	global.AddGroup("Global Options", "", &GlobalOption{})
	return &GlobalOption{}
}

func (g *GlobalOption) IsValid() error {
	var errMsg []string
	var tmpErr error

	tmpErr = validRepository(g.Repository)
	if tmpErr != nil {
		errMsg = append(errMsg, tmpErr.Error())
	}

	if len(errMsg) > 0 {
		return fmt.Errorf("Invalid value in global option. %v", errMsg)
	}
	return nil
}

func validRepository(value string) error {
	splited := strings.Split(value, "/")
	if value != "" && len(splited) != 2 {
		return fmt.Errorf("Invalid repository \"%s\". Assumed input style is \"Namespace/Project\".", value)
	}
	return nil
}

func (g *GlobalOption) NameSpaceAndProject() (namespace, project string) {
	splited := strings.Split(g.Repository, "/")
	namespace = splited[0]
	project = splited[1]
	return
}

type SearchOption struct {
	State      string `short:"t" long:"state" default:"all" default-mask:"all" description:"just those that are opened, closed or all"`
	Scope      string `short:"c" long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy    string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Opened     bool   `short:"e" long:"opened" description:"search state opened"`
	Closed     bool   `short:"l" long:"closed" description:"search scope closed"`
	CreatedMe  bool   `short:"r" long:"created-me" description:"search scope created-by-me"`
	AssignedMe bool   `short:"g" long:"assigned-me" description:"search scope assigned-to-me"`
	AllProject bool   `short:"a" long:"all-project" description:"search target all project"`
}

func newSearchOption() *SearchOption {
	search := flags.NewNamedParser("lab", flags.Default)
	search.AddGroup("Search Options", "", &SearchOption{})
	return &SearchOption{}
}

func (s *SearchOption) GetState() string {
	if s.Opened {
		return "opened"
	}
	if s.Closed {
		return "closed"
	}
	return s.State
}

func (s *SearchOption) GetScope() string {
	if s.CreatedMe {
		return "created-by-me"
	}
	if s.AssignedMe {
		return "assigned-to-me"
	}
	return s.Scope
}

type OutputOption struct {
	Line int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	Sort string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
}

func newOutputOption() *OutputOption {
	output := flags.NewNamedParser("lab", flags.Default)
	output.AddGroup("Output Options", "", &OutputOption{})
	return &OutputOption{}
}
