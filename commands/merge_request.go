package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type CreateUpdateMergeRequestOption struct {
	Edit         bool   `short:"e" long:"edit" description:"Edit the merge request on editor. Start the editor with the contents in the given title and message options."`
	Title        string `short:"i" long:"title" value-name:"<title>" description:"The title of an merge request"`
	Message      string `short:"m" long:"message" value-name:"<message>" description:"The message of an merge request"`
	SourceBranch string `short:"s" long:"source" description:"The source branch"`
	TargetBranch string `short:"t" long:"target" default:"master" default-mask:"master" description:"The target branch"`
}

func newCreateUpdateMergeRequestOption() *CreateUpdateMergeRequestOption {
	return &CreateUpdateMergeRequestOption{}
}

type ListMergeRequestOption struct {
	Num        int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of merge request to output."`
	State      string `long:"state" value-name:"<state>" default:"all" default-mask:"all" description:"Print only merge request of the state just those that are \"opened\", \"closed\", \"merged\" or \"all\""`
	Scope      string `long:"scope" value-name:"<scope>" default:"all" default-mask:"all" description:"Print only given scope. \"created-by-me\", \"assigned-to-me\" or \"all\"."`
	OrderBy    string `long:"orderby" value-name:"<orderby>" default:"updated_at" default-mask:"updated_at" description:"Print merge request ordered by \"created_at\" or \"updated_at\" fields."`
	Sort       string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print merge request ordered in \"asc\" or \"desc\" order."`
	Opened     bool   `short:"o" long:"opened" description:"Shorthand of the state option for \"--state=opened\"."`
	Closed     bool   `short:"c" long:"closed" description:"Shorthand of the state option for \"--state=closed\"."`
	Merged     bool   `short:"g" long:"merged" description:"Shorthand of the state option for \"--state=merged\"."`
	CreatedMe  bool   `short:"r" long:"created-me" description:"Shorthand of the scope option for \"--scope=created-by-me\"."`
	AssignedMe bool   `short:"a" long:"assigned-me" description:"Shorthand of the scope option for \"--scope=assigned-by-me\"."`
	AllProject bool   `short:"A" long:"all-project" description:"Print the merge request of all projects"`
}

func (l *ListMergeRequestOption) GetState() string {
	if l.Opened {
		return "opened"
	}
	if l.Closed {
		return "closed"
	}
	if l.Merged {
		return "merged"
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

type MergeRequestOperation int

const (
	CreateMergeRequest MergeRequestOperation = iota
	CreateMergeRequestOnEditor
	UpdateMergeRequest
	UpdateMergeRequestOnEditor
	ShowMergeRequest
	ListMergeRequest
	ListMergeRequestAllProject
)

func newListMergeRequestOption() *ListMergeRequestOption {
	return &ListMergeRequestOption{}
}

type MergeRequestCommandOption struct {
	CreateUpdateOption *CreateUpdateMergeRequestOption `group:"Create, Update Options"`
	ListOption         *ListMergeRequestOption         `group:"List Options"`
}

func newMergeRequestOptionParser(opt *MergeRequestCommandOption) *flags.Parser {
	opt.CreateUpdateOption = newCreateUpdateMergeRequestOption()
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
	EditFunc func(program, file string) error
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

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	iid, err := validMergeRequestIID(parseArgs)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	switch mergeRequestOperation(mergeRequestCommandOption, parseArgs) {
	case UpdateMergeRequest:
		// Getting exist merge request
		mergeRequest, err := client.GetMergeRequest(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Create new title or description
		createUpdateOption := mergeRequestCommandOption.CreateUpdateOption
		updatedTitle := mergeRequest.Title
		updatedMessage := mergeRequest.Description
		if createUpdateOption.Title != "" {
			updatedTitle = createUpdateOption.Title
		}
		if createUpdateOption.Message != "" {
			updatedMessage = createUpdateOption.Message
		}

		// Do update merge request
		updatedMergeRequest, err := client.UpdateMergeRequest(
			makeUpdateMergeRequestOption(createUpdateOption, updatedTitle, updatedMessage),
			iid,
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print update Issue IID
		c.Ui.Message(fmt.Sprintf("#%d", updatedMergeRequest.IID))

	case UpdateMergeRequestOnEditor:
		createUpdateOption := mergeRequestCommandOption.CreateUpdateOption

		// Getting exist merge request
		mergeRequest, err := client.GetMergeRequest(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Starting editor for edit title and description
		template := editMergeRequestTemplate(mergeRequest.Title, mergeRequest.Description)
		title, message, err := editIssueTitleAndDesc(template, c.EditFunc)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Do update merge request
		updatedMergeRequest, err := client.UpdateMergeRequest(
			makeUpdateMergeRequestOption(createUpdateOption, title, message),
			iid,
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print update Issue IID
		c.Ui.Message(fmt.Sprintf("#%d", updatedMergeRequest.IID))

	case CreateMergeRequest:
		// Get source branch. current branch from local repository when non specific flags
		createUpdateOption := mergeRequestCommandOption.CreateUpdateOption
		currentBranch, err := git.GitCurrentBranch()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		if createUpdateOption.SourceBranch != "" {
			// TODO Checking source branch exitst
			currentBranch = createUpdateOption.SourceBranch
		}

		// Do create merge request
		mergeRequest, err := client.CreateMergeRequest(
			makeCreateMergeRequestOption(createUpdateOption, createUpdateOption.Title, createUpdateOption.Message, currentBranch),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print created merge request IID
		c.Ui.Message(fmt.Sprintf("!%d", mergeRequest.IID))

	case CreateMergeRequestOnEditor:
		// Starting editor for edit title and description
		createUpdateOption := mergeRequestCommandOption.CreateUpdateOption
		template := editMergeRequestTemplate(createUpdateOption.Title, createUpdateOption.Message)
		title, message, err := editIssueTitleAndDesc(template, c.EditFunc)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Get source branch. current branch from local repository when non specific flags
		currentBranch, err := git.GitCurrentBranch()
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		if createUpdateOption.SourceBranch != "" {
			// TODO Checking source branch exitst
			currentBranch = createUpdateOption.SourceBranch
		}

		// Do create merge request
		mergeRequest, err := client.CreateMergeRequest(
			makeCreateMergeRequestOption(createUpdateOption, title, message, currentBranch),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print created merge request IID
		c.Ui.Message(fmt.Sprintf("!%d", mergeRequest.IID))

	case ShowMergeRequest:
		// Do get merge request
		mergeRequest, err := client.GetMergeRequest(iid, gitlabRemote.RepositoryFullName())
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		output := outMergeRequestDetail(mergeRequest)
		c.Ui.Message(output)

	case ListMergeRequest:
		listOption := mergeRequestCommandOption.ListOption
		mergeRequests, err := client.ProjectMergeRequest(
			makeProjectMergeRequestOption(listOption),
			gitlabRemote.RepositoryFullName(),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}
		outputs := outProjectMergeRequest(mergeRequests)
		result := columnize.SimpleFormat(outputs)
		c.Ui.Message(result)

	case ListMergeRequestAllProject:
		// Do get merge request list
		listOption := mergeRequestCommandOption.ListOption
		mergeRequests, err := client.MergeRequest(
			makeMergeRequestOption(listOption),
		)
		if err != nil {
			c.Ui.Error(err.Error())
			return ExitCodeError
		}

		// Print merge request list
		outputs := outMergeRequest(mergeRequests)
		result := columnize.SimpleFormat(outputs)
		c.Ui.Message(result)

	default:
		c.Ui.Error("Invalid merge request operation")
		return ExitCodeError
	}

	return ExitCodeOK
}

func mergeRequestOperation(opt MergeRequestCommandOption, args []string) MergeRequestOperation {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption

	// Case of getting Merge Request IID
	if len(args) > 0 {
		if createUpdateOption.Edit {
			return UpdateMergeRequestOnEditor
		}
		if createUpdateOption.Title != "" || createUpdateOption.Message != "" {
			return UpdateMergeRequest
		}
		return ShowMergeRequest
	}

	// Case of nothing MergeRequest IID
	if createUpdateOption.Edit {
		return CreateMergeRequestOnEditor
	}
	if createUpdateOption.Title != "" {
		return CreateMergeRequest
	}
	if listOption.AllProject {
		return ListMergeRequestAllProject
	}

	return ListMergeRequest
}

func validMergeRequestIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid Issue IID. IID: %s", args[0])
	}
	return iid, nil
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

func makeCreateMergeRequestOption(opt *CreateUpdateMergeRequestOption, title, description, branch string) *gitlabc.CreateMergeRequestOptions {
	createMergeRequestOption := &gitlabc.CreateMergeRequestOptions{
		Title:           gitlabc.String(title),
		Description:     gitlabc.String(description),
		SourceBranch:    gitlabc.String(branch),
		TargetBranch:    gitlabc.String(opt.TargetBranch),
		AssigneeID:      nil,
		TargetProjectID: nil,
	}
	return createMergeRequestOption
}

func makeUpdateMergeRequestOption(opt *CreateUpdateMergeRequestOption, title, description string) *gitlabc.UpdateMergeRequestOptions {
	updateMergeRequestOptions := &gitlabc.UpdateMergeRequestOptions{
		Title:        gitlabc.String(title),
		Description:  gitlabc.String(description),
		TargetBranch: gitlabc.String(opt.TargetBranch),
		AssigneeID:   nil,
	}
	return updateMergeRequestOptions
}

func outMergeRequest(mergeRequsets []*gitlabc.MergeRequest) []string {
	outputs := []string{}
	for _, mergeRequest := range mergeRequsets {
		output := strings.Join([]string{
			gitlab.ParceRepositoryFullName(mergeRequest.WebURL),
			fmt.Sprintf("!%d", mergeRequest.IID),
			mergeRequest.Title,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}

func outMergeRequestDetail(mergeRequest *gitlabc.MergeRequest) string {
	base := `#%d
Title: %s
Assignee: %s
Author: %s
CreatedAt: %s
UpdatedAt: %s

%s`
	detial := fmt.Sprintf(
		base,
		mergeRequest.IID,
		mergeRequest.Title,
		mergeRequest.Assignee.Name,
		mergeRequest.Author.Name,
		mergeRequest.CreatedAt.String(),
		mergeRequest.UpdatedAt.String(),
		mergeRequest.Description,
	)
	return detial
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

func editMergeRequestTemplate(title, description string) string {
	message := `<!-- Write a message for this merge request. The first block of text is the title -->
%s

<!-- the rest is the description.  -->
%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}
