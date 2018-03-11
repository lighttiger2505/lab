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

var issueCommandOption IssueCommnadOption
var issueCommnadParser *flags.Parser = newIssueOptionParser(&issueCommandOption)

type IssueCommnadOption struct {
	GlobalOption *GlobalOption `group:"Global Options"`
	SearchOption *SearchOption `group:"Search Options"`
	OutputOption *OutputOption `group:"Output Options"`
}

func newIssueOptionParser(opt *IssueCommnadOption) *flags.Parser {
	opt.GlobalOption = newGlobalOption()
	opt.SearchOption = newSearchOption()
	opt.OutputOption = newOutputOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "add-issue [options]"
	return parser
}

type IssueCommand struct {
	Ui       ui.Ui
	Provider gitlab.Provider
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	issueCommnadParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	if _, err := issueCommnadParser.ParseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	globalOption := issueCommandOption.GlobalOption
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
	if globalOption.Repository != "" {
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

	var datas []string
	searchOption := issueCommandOption.SearchOption
	outputOption := issueCommandOption.OutputOption
	if searchOption.AllProject {
		issues, err := client.Issues(makeIssueOption(searchOption, outputOption))
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = issueOutput(issues)

	} else {
		issues, err := client.ProjectIssues(
			makeProjectIssueOption(searchOption, outputOption),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = projectIssueOutput(issues)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Message(result)

	return ExitCodeOK
}

func makeProjectIssueOption(searchOption *SearchOption, outputOption *OutputOption) *gitlabc.ListProjectIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: outputOption.Line,
	}
	listProjectIssuesOptions := &gitlabc.ListProjectIssuesOptions{
		State:       gitlabc.String(searchOption.GetState()),
		Scope:       gitlabc.String(searchOption.GetScope()),
		OrderBy:     gitlabc.String(searchOption.OrderBy),
		Sort:        gitlabc.String(outputOption.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func makeIssueOption(searchOption *SearchOption, outputOption *OutputOption) *gitlabc.ListIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: outputOption.Line,
	}
	listIssuesOptions := &gitlabc.ListIssuesOptions{
		State:       gitlabc.String(searchOption.GetState()),
		Scope:       gitlabc.String(searchOption.GetScope()),
		OrderBy:     gitlabc.String(searchOption.OrderBy),
		Sort:        gitlabc.String(outputOption.Sort),
		ListOptions: *listOption,
	}
	return listIssuesOptions
}

func issueOutput(issues []*gitlabc.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("#%d", issue.IID),
			gitlab.ParceRepositoryFullName(issue.WebURL),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}

func projectIssueOutput(issues []*gitlabc.Issue) []string {
	var datas []string
	for _, issue := range issues {
		data := strings.Join([]string{
			fmt.Sprintf("#%d", issue.IID),
			issue.Title,
		}, "|")
		datas = append(datas, data)
	}
	return datas
}
