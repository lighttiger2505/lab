package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

var mergeRequestOpt MergeRequestOpt
var mergeRequestParser *flags.Parser = newMergeRequestOptionParser(&mergeRequestOpt)

type MergeRequestOpt struct {
	GlobalOpt *GlobalOpt `group:"Global Options"`
	SearchOpt *SearchOpt `group:"Search Options"`
}

func newMergeRequestOptionParser(mrOpt *MergeRequestOpt) *flags.Parser {
	globalParser := flags.NewParser(&globalOpt, flags.Default)
	globalParser.AddGroup("Global Options", "", &GlobalOpt{})

	searchParser := flags.NewParser(&searchOptions, flags.Default)
	searchParser.AddGroup("Search Options", "", &GlobalOpt{})

	parser := flags.NewParser(mrOpt, flags.Default)
	parser.Usage = "merge-request [options]"
	return parser
}

type MergeRequestCommand struct {
	Ui       ui.Ui
	Provider gitlab.Provider
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	mergeRequestParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	if _, err := mergeRequestParser.ParseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	globalOpt := browseOpt.GlobalOpt
	if err := globalOpt.IsValid(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	var gitlabRemote *git.RemoteInfo
	if globalOpt.Repository != "" {
		namespace, project := globalOpt.NameSpaceAndProject()
		gitlabRemote = c.Provider.GetSpecificRemote(namespace, project)
	} else {
		var err error
		gitlabRemote, err = c.Provider.GetCurrentRemote()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
	}

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var outputs []string
	if mergeRequestOpt.SearchOpt.AllRepository {
		mergeRequests, err := client.MergeRequest(
			makeMergeRequestOption(mergeRequestOpt.SearchOpt),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		outputs = outMergeRequest(mergeRequests)
	} else {
		mergeRequests, err := client.ProjectMergeRequest(
			makeProjectMergeRequestOption(mergeRequestOpt.SearchOpt),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		outputs = outProjectMergeRequest(mergeRequests)
	}

	result := columnize.SimpleFormat(outputs)
	c.Ui.Message(result)

	return ExitCodeOK
}

func makeMergeRequestOption(opt *SearchOpt) *gitlabc.ListMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: opt.Line,
	}
	listRequestsOptions := &gitlabc.ListMergeRequestsOptions{
		State:       gitlabc.String(opt.GetState()),
		Scope:       gitlabc.String(opt.GetScope()),
		OrderBy:     gitlabc.String(opt.OrderBy),
		Sort:        gitlabc.String(opt.Sort),
		ListOptions: *listOption,
	}
	return listRequestsOptions
}

func makeProjectMergeRequestOption(opt *SearchOpt) *gitlabc.ListProjectMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: opt.Line,
	}
	listMergeRequestsOptions := &gitlabc.ListProjectMergeRequestsOptions{
		State:       gitlabc.String(opt.GetState()),
		Scope:       gitlabc.String(opt.GetScope()),
		OrderBy:     gitlabc.String(opt.OrderBy),
		Sort:        gitlabc.String(opt.Sort),
		ListOptions: *listOption,
	}
	return listMergeRequestsOptions
}

func outMergeRequest(mergeRequsets []*gitlabc.MergeRequest) []string {
	outputs := []string{}
	for _, mergeRequest := range mergeRequsets {
		output := strings.Join([]string{
			fmt.Sprintf("!%d", mergeRequest.IID),
			gitlab.ParceRepositoryFullName(mergeRequest.WebURL),
			mergeRequest.Title,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}

func outProjectMergeRequest(mergeRequsets []*gitlabc.MergeRequest) []string {
	outputs := []string{}
	for _, mergeRequest := range mergeRequsets {
		output := strings.Join([]string{
			fmt.Sprintf("!%d", mergeRequest.IID),
			mergeRequest.Title,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
