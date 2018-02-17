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

var mergeRequestOpt MergeRequestOpt

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
	Ui ui.Ui
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	newMergeRequestOptionParser(&mergeRequestOpt).WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	parser := newMergeRequestOptionParser(&mergeRequestOpt)
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

	client, err := gitlab.GitlabClient(c.Ui, gitlabRemote, conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Replace specific repository
	domain := conf.MustDomain()
	if mergeRequestOpt.GlobalOpt.Repository != "" {
		namespace, project, err := mergeRequestOpt.GlobalOpt.ValidRepository()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		gitlabRemote.Domain = domain
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	var datas []string
	if mergeRequestOpt.SearchOpt.AllRepository {
		mergeRequests, err := getMergeRequest(client, mergeRequestOpt.SearchOpt)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		for _, mergeRequest := range mergeRequests {
			data := strings.Join([]string{
				fmt.Sprintf("!%d", mergeRequest.IID),
				gitlab.ParceRepositoryFullName(mergeRequest.WebURL),
				mergeRequest.Title,
			}, "|")
			datas = append(datas, data)
		}
	} else {
		mergeRequests, err := getProjectMergeRequest(client, mergeRequestOpt.SearchOpt, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		for _, mergeRequest := range mergeRequests {
			data := strings.Join([]string{
				fmt.Sprintf("!%d", mergeRequest.IID),
				mergeRequest.Title,
			}, "|")
			datas = append(datas, data)
		}
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Message(result)

	return ExitCodeOK
}

func getMergeRequest(client *gitlabc.Client, opt *SearchOpt) ([]*gitlabc.MergeRequest, error) {
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

	mergeRequests, _, err := client.MergeRequests.ListMergeRequests(
		listRequestsOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}

func getProjectMergeRequest(client *gitlabc.Client, opt *SearchOpt, repositoryName string) ([]*gitlabc.MergeRequest, error) {
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

	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(
		repositoryName,
		listMergeRequestsOptions,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed list project merge requests. %s", err.Error())
	}

	return mergeRequests, nil
}
