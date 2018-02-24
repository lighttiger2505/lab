package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

var issueOpt IssueOpt
var issueOptionParser *flags.Parser = newIssueOptionParser(&issueOpt)

type IssueOpt struct {
	GlobalOpt *GlobalOpt `group:"Global Options"`
	SearchOpt *SearchOpt `group:"Search Options"`
}

func newIssueOptionParser(issueOpt *IssueOpt) *flags.Parser {
	globalParser := flags.NewParser(&globalOpt, flags.Default)
	globalParser.AddGroup("Global Options", "", &GlobalOpt{})

	searchParser := flags.NewParser(&searchOptions, flags.Default)
	searchParser.AddGroup("Search Options", "", &GlobalOpt{})

	parser := flags.NewParser(issueOpt, flags.Default)
	parser.Usage = "issue [options]"
	return parser
}

type IssueCommand struct {
	Ui           ui.Ui
	RemoteFilter gitlab.RemoteFilter
	GitClient    git.Client
	LabClient    gitlab.Client
	Config       *config.ConfigManager
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	issueOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	if _, err := issueOptionParser.ParseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	opt := issueOpt.GlobalOpt
	if err := opt.IsValid(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Load config
	if err := c.Config.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	conf, err := c.Config.Load()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	if err := c.RemoteFilter.Collect(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	gitlabRemote, err := c.RemoteFilter.Filter(c.Ui, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Replace specific repository
	domain := conf.MustDomain()
	if issueOpt.GlobalOpt.Repository != "" {
		namespace, project := issueOpt.GlobalOpt.NameSpaceAndProject()
		gitlabRemote.Domain = domain
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	fmt.Println(gitlabRemote.Domain)
	token, err := gitlab.GetPrivateToken(c.Ui, gitlabRemote.Domain, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	if issueOpt.SearchOpt.AllRepository {
		issues, err := c.LabClient.Issues(
			gitlabRemote.ApiUrl(),
			token,
			makeIssueOption(issueOpt.SearchOpt),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = issueOutput(issues)

	} else {
		issues, err := c.LabClient.ProjectIssues(
			gitlabRemote.ApiUrl(),
			token,
			makeProjectIssueOption(issueOpt.SearchOpt),
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

func makeProjectIssueOption(opt *SearchOpt) *gitlabc.ListProjectIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: opt.Line,
	}
	listProjectIssuesOptions := &gitlabc.ListProjectIssuesOptions{
		State:       gitlabc.String(opt.GetState()),
		Scope:       gitlabc.String(opt.GetScope()),
		OrderBy:     gitlabc.String(opt.OrderBy),
		Sort:        gitlabc.String(opt.Sort),
		ListOptions: *listOption,
	}
	return listProjectIssuesOptions
}

func makeIssueOption(opt *SearchOpt) *gitlabc.ListIssuesOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: opt.Line,
	}
	listIssuesOptions := &gitlabc.ListIssuesOptions{
		State:       gitlabc.String(opt.GetState()),
		Scope:       gitlabc.String(opt.GetScope()),
		OrderBy:     gitlabc.String(opt.OrderBy),
		Sort:        gitlabc.String(opt.Sort),
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
