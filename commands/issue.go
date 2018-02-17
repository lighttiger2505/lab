package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

var issueOpt IssueOpt

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
	Ui ui.Ui
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	newIssueOptionParser(&issueOpt).WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	parser := newIssueOptionParser(&issueOpt)
	if _, err := parser.Parse(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	conf, err := config.NewConfig()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := gitlab.GitlabRemote(c.Ui, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Replace specific repository
	domain := conf.MustDomain()
	if issueOpt.GlobalOpt.Repository != "" {
		namespace, project, err := issueOpt.GlobalOpt.ValidRepository()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		gitlabRemote.Domain = domain
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	client, err := gitlab.GitlabClient(c.Ui, gitlabRemote, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	if issueOpt.SearchOpt.AllRepository {
		issues, err := getIssues(client, issueOpt.SearchOpt)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		datas = issueOutput(issues)

	} else {
		issues, err := getProjectIssues(client, issueOpt.SearchOpt, gitlabRemote.RepositoryFullName())
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

func getIssues(client *gitlabc.Client, opt *SearchOpt) ([]*gitlabc.Issue, error) {
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

	issues, _, err := client.Issues.ListIssues(
		listIssuesOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list issue. %s", err.Error())
	}

	return issues, nil
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

func getProjectIssues(client *gitlabc.Client, opt *SearchOpt, repositoryName string) ([]*gitlabc.Issue, error) {
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

	issues, _, err := client.Issues.ListProjectIssues(
		repositoryName,
		listProjectIssuesOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project issue. %s", err.Error())
	}

	return issues, nil
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
