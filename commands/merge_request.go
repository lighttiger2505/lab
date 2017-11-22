package commands

import (
	"bytes"
	"fmt"

	"github.com/lighttiger2505/lab/config"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type MergeRequestCommand struct {
	Ui ui.Ui
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	searchParser.Usage = "merge-request [options]"
	searchParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	if _, err := searchParser.Parse(); err != nil {
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

	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: searchOptions.Line,
	}
	listMergeRequestsOptions := &gitlabc.ListProjectMergeRequestsOptions{
		State:       gitlabc.String(searchOptions.State),
		Scope:       gitlabc.String(searchOptions.Scope),
		OrderBy:     gitlabc.String(searchOptions.OrderBy),
		Sort:        gitlabc.String(searchOptions.Sort),
		ListOptions: *listOption,
	}
	mergeRequests, _, err := client.MergeRequests.ListProjectMergeRequests(
		gitlabRemote.RepositoryFullName(),
		listMergeRequestsOptions,
	)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	for _, mergeRequest := range mergeRequests {
		data := fmt.Sprintf("!%d", mergeRequest.IID) + "|" + mergeRequest.Title
		datas = append(datas, data)
	}

	result := columnize.SimpleFormat(datas)
	c.Ui.Message(result)

	return ExitCodeOK
}