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

type AddMergeRequestOption struct {
	Add          bool   `short:"a" long:"add" description:"add issue"`
	Message      string `short:"m" long:"message" description:"issue description"`
	SourceBranch string `short:"s" long:"source" description:"The source branch"`
	TargetBranch string `short:"g" long:"target" default:"master" default-mask:"master" description:"The target branch"`
}

func newAddMergeRequestOption() *AddMergeRequestOption {
	return &AddMergeRequestOption{}
}

type ListMergeRequestOption struct {
	List       bool   `short:"l" long:"line" description:"show list"`
	Num        int    `short:"n" long:"num"  default:"20" default-mask:"20" description:"show issue num"`
	State      string `long:"state" default:"all" default-mask:"all" description:"those that are opened, closed, or merged"`
	Scope      string `long:"scope" default:"all" default-mask:"all" description:"given scope: created-by-me, assigned-to-me or all."`
	OrderBy    string `long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by created_at or updated_at fields."`
	Sort       string `long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order."`
	Opened     bool   `short:"o" long:"opened" description:"search state opened"`
	Closed     bool   `short:"c" long:"closed" description:"search scope closed"`
	CreatedMe  bool   `short:"r" long:"created-me" description:"search scope created-by-me"`
	AssignedMe bool   `long:"s" long:"assigned-me" description:"search scope assigned-to-me"`
	AllProject bool   `long:"a" long:"all-project" description:"search target all project"`
}

func (l *ListMergeRequestOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	return l.State
}

func (l *ListMergeRequestOption) GetScope() string {
	if l.CreatedMe {
		return "created-by-me"
	}
	if l.AssignedMe {
		return "assigned-to-me"
	}
	return l.Scope
}

func newListMergeRequestOption() *ListMergeRequestOption {
	return &ListMergeRequestOption{}
}

type MergeRequestCommandOption struct {
	GlobalOption *GlobalOption           `group:"Global Options"`
	AddOption    *AddMergeRequestOption  `group:"Create Options"`
	ListOption   *ListMergeRequestOption `group:"List Options"`
}

func newMergeRequestOptionParser(opt *MergeRequestCommandOption) *flags.Parser {
	opt.GlobalOption = newGlobalOption()
	opt.AddOption = newAddMergeRequestOption()
	opt.ListOption = newListMergeRequestOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = `merge-request [options]

Synopsis:
    lab merge-request -a <title> [-d <message>]
    lab merge-request [-n <num>] -l [--state <state>] [--scope <scope>]
                      [--orderby <orderby>] [--sort <sort>] -o -c 
                      -cm -am -al
	lab merge-request [-t <title>] [-d <description>] <id>
`
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
	var mergeRequestCommandOption MergeRequestCommandOption
	mergeRequestCommandParser := newMergeRequestOptionParser(&mergeRequestCommandOption)
	mergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	var mergeRequestCommandOption MergeRequestCommandOption
	mergeRequestCommandParser := newMergeRequestOptionParser(&mergeRequestCommandOption)
	parseArgs, err := mergeRequestCommandParser.ParseArgs(args)
	if err != nil {
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

	addOption := mergeRequestCommandOption.AddOption
	if addOption.Add {
		argsTitle := ""
		if len(parseArgs) < 0 {
			argsTitle = parseArgs[0]
		}

		title, message, err := getMergeRequestTitleAndDesc(argsTitle, addOption.Message)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Get source branch
		// current branch from local repository when non specific flags
		currentBranch, err := git.GitCurrentBranch()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		if addOption.SourceBranch != "" {
			currentBranch = addOption.SourceBranch
		}

		mergeRequest, err := client.CreateMergeRequest(
			makeCreateMergeRequestOption(addOption, title, message, currentBranch),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		c.Ui.Message(fmt.Sprintf("!%d", mergeRequest.IID))
		return ExitCodeOK
	}

	var outputs []string
	listOption := mergeRequestCommandOption.ListOption
	if listOption.AllProject {
		mergeRequests, err := client.MergeRequest(
			makeMergeRequestOption(listOption),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		outputs = outMergeRequest(mergeRequests)
	} else {
		mergeRequests, err := client.ProjectMergeRequest(
			makeProjectMergeRequestOption(listOption),
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

func makeMergeRequestOption(listMergeRequestsOption *ListMergeRequestOption) *gitlabc.ListMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: listMergeRequestsOption.Num,
	}
	listRequestsOptions := &gitlabc.ListMergeRequestsOptions{
		State:       gitlabc.String(listMergeRequestsOption.GetState()),
		Scope:       gitlabc.String(listMergeRequestsOption.GetScope()),
		OrderBy:     gitlabc.String(listMergeRequestsOption.OrderBy),
		Sort:        gitlabc.String(listMergeRequestsOption.Sort),
		ListOptions: *listOption,
	}
	return listRequestsOptions
}

func makeProjectMergeRequestOption(listMergeRequestsOption *ListMergeRequestOption) *gitlabc.ListProjectMergeRequestsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: listMergeRequestsOption.Num,
	}
	listMergeRequestsOptions := &gitlabc.ListProjectMergeRequestsOptions{
		State:       gitlabc.String(listMergeRequestsOption.GetState()),
		Scope:       gitlabc.String(listMergeRequestsOption.GetScope()),
		OrderBy:     gitlabc.String(listMergeRequestsOption.OrderBy),
		Sort:        gitlabc.String(listMergeRequestsOption.Sort),
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
