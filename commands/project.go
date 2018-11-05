package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type ProjectCommnadOption struct {
	OutputOption *ListProjectOption `group:"List Options"`
}

func newProjectCommandParser(opt *ProjectCommnadOption) *flags.Parser {
	opt.OutputOption = newListProjectOption()
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "project [options]"
	return parser
}

type ListProjectOption struct {
	Num        int    `short:"n" long:"num" value-name:"<num>" default:"20" default-mask:"20" description:"Limit the number of project to output."`
	Sort       string `long:"sort"  value-name:"<sort>" default:"desc" default-mask:"desc" description:"Print project ordered in \"asc\" or \"desc\" order."`
	OrderBy    string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by id, name, path, created_at, updated_at, or last_activity_at fields"`
	Owned      bool   `short:"w" long:"owned" description:"Limit by projects owned by the current user"`
	Membership bool   `short:"m" long:"member-ship" description:"Limit by projects that the current user is a member of"`
}

func newListProjectOption() *ListProjectOption {
	return &ListProjectOption{}
}

type ProjectCommand struct {
	UI            ui.Ui
	Provider      lab.Provider
	ClientFactory lab.APIClientFactory
}

func (c *ProjectCommand) Synopsis() string {
	return "List project"
}

func (c *ProjectCommand) Help() string {
	buf := &bytes.Buffer{}
	var projectCommandOption ProjectCommnadOption
	projectCommandParser := newProjectCommandParser(&projectCommandOption)
	projectCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *ProjectCommand) Run(args []string) int {
	var projectCommandOption ProjectCommnadOption
	projectCommandParser := newProjectCommandParser(&projectCommandOption)
	if _, err := projectCommandParser.ParseArgs(args); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if err := c.Provider.Init(); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	token, err := c.Provider.GetAPIToken(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if err := c.ClientFactory.Init(gitlabRemote.ApiUrl(), token); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	client := c.ClientFactory.GetProjectClient()

	projects, err := client.Projects(
		makeProjectOptions(projectCommandOption.OutputOption),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(projectOutput(projects))
	c.UI.Message(result)

	return ExitCodeOK
}

func makeProjectOptions(listProjectOption *ListProjectOption) *gitlab.ListProjectsOptions {
	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: listProjectOption.Num,
	}
	listProjectsOptions := &gitlab.ListProjectsOptions{
		Archived:    gitlab.Bool(false),
		OrderBy:     gitlab.String(listProjectOption.OrderBy),
		Sort:        gitlab.String(listProjectOption.Sort),
		Search:      gitlab.String(""),
		Simple:      gitlab.Bool(false),
		Owned:       gitlab.Bool(listProjectOption.Owned),
		Membership:  gitlab.Bool(listProjectOption.Membership),
		Starred:     gitlab.Bool(false),
		Statistics:  gitlab.Bool(false),
		Visibility:  gitlab.Visibility("private"),
		ListOptions: *listOption,
	}
	return listProjectsOptions
}

func removeLineBreak(value string) string {
	value = strings.Replace(value, "\r\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)
	value = strings.Replace(value, "\n", "", -1)
	return value
}

func projectOutput(projects []*gitlab.Project) []string {
	var outputs []string
	for _, project := range projects {
		output := strings.Join([]string{
			fmt.Sprintf("%s/%s", project.Namespace.Name, project.Name),
			removeLineBreak(project.Description),
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
