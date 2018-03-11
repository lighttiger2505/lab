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

var mergeRequestCommandOption MergeRequestCommandOption
var mergeRequestCommandParser *flags.Parser = newMergeRequestOptionParser(&mergeRequestCommandOption)

type MergeRequestCommandOption struct {
	GlobalOption *GlobalOption `group:"Global Options"`
	SearchOption *SearchOption `group:"Search Options"`
	OutputOption *OutputOption `group:"Output Options"`
}

func newMergeRequestOptionParser(opt *MergeRequestCommandOption) *flags.Parser {
	opt.GlobalOption = newGlobalOption()
	opt.SearchOption = newSearchOption()
	opt.OutputOption = newOutputOption()
	parser := flags.NewParser(opt, flags.Default)
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
	mergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	if _, err := mergeRequestCommandParser.ParseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	globalOption := mergeRequestCommandOption.GlobalOption
	if err := globalOption.IsValid(); err != nil {
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
	if globalOption.Project != "" {
		namespace, project := globalOption.NameSpaceAndProject()
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
	searchOption := mergeRequestCommandOption.SearchOption
	outputOption := mergeRequestCommandOption.OutputOption
	if searchOption.AllProject {
		mergeRequests, err := client.MergeRequest(
			makeMergeRequestOption(searchOption, outputOption),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		outputs = outMergeRequest(mergeRequests)
	} else {
		mergeRequests, err := client.ProjectMergeRequest(
			makeProjectMergeRequestOption(searchOption, outputOption),
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

func makeMergeRequestOption(searchOption *SearchOption, outputOption *OutputOption) *gitlabc.ListMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: outputOption.Line,
	}
	listRequestsOptions := &gitlabc.ListMergeRequestsOptions{
		State:       gitlabc.String(searchOption.GetState()),
		Scope:       gitlabc.String(searchOption.GetScope()),
		OrderBy:     gitlabc.String(searchOption.OrderBy),
		Sort:        gitlabc.String(outputOption.Sort),
		ListOptions: *listOption,
	}
	return listRequestsOptions
}

func makeProjectMergeRequestOption(searchOption *SearchOption, outputOption *OutputOption) *gitlabc.ListProjectMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: outputOption.Line,
	}
	listMergeRequestsOptions := &gitlabc.ListProjectMergeRequestsOptions{
		State:       gitlabc.String(searchOption.GetState()),
		Scope:       gitlabc.String(searchOption.GetScope()),
		OrderBy:     gitlabc.String(searchOption.OrderBy),
		Sort:        gitlabc.String(outputOption.Sort),
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
