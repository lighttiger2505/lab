package commands

import (
	"bytes"
	"fmt"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type ProjectCommnadOption struct {
	ProjectOption *ProjectOption `group:"Project Options"`
	OutputOption  *OutputOption  `group:"Output Options"`
}

func newIssueCommandParser(opt *ProjectCommnadOption) *flags.Parser {
	opt.ProjectOption = newProjectOption()
	opt.OutputOption = newOutputOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "project [options]"
	return parser
}

type ProjectOption struct {
	OrderBy    string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by id, name, path, created_at, updated_at, or last_activity_at fields"`
	Owned      bool   `short:"w" long:"owned" description:"Limit by projects owned by the current user"`
	Membership bool   `short:"m" long:"member-ship" description:"Limit by projects that the current user is a member of"`
}

func newProjectOption() *ProjectOption {
	project := flags.NewNamedParser("lab", flags.Default)
	project.AddGroup("Project Options", "", &ProjectOption{})
	return &ProjectOption{}
}

type ProjectCommand struct {
	UI       ui.Ui
	Provider gitlab.Provider
}

func (c *ProjectCommand) Synopsis() string {
	return "Show project"
}

func (c *ProjectCommand) Help() string {
	buf := &bytes.Buffer{}
	var projectCommandOption ProjectCommnadOption
	projectCommandParser := newIssueCommandParser(&projectCommandOption)
	projectCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *ProjectCommand) Run(args []string) int {
	// Parse flags
	var projectCommandOption ProjectCommnadOption
	projectCommandParser := newIssueCommandParser(&projectCommandOption)
	if _, err := projectCommandParser.ParseArgs(args); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	client, err := c.Provider.GetClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	projects, err := client.Projects(
		makeProjectOptions(projectCommandOption.ProjectOption, projectCommandOption.OutputOption),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(projectOutput(projects))
	c.UI.Message(result)

	return ExitCodeOK
}

func makeProjectOptions(projectOption *ProjectOption, outputOption *OutputOption) *gitlabc.ListProjectsOptions {
	listOption := &gitlabc.ListOptions{
		Page:    1,
		PerPage: outputOption.Line,
	}
	listProjectsOptions := &gitlabc.ListProjectsOptions{
		Archived:    gitlabc.Bool(false),
		OrderBy:     gitlabc.String(projectOption.OrderBy),
		Sort:        gitlabc.String(outputOption.Sort),
		Search:      gitlabc.String(""),
		Simple:      gitlabc.Bool(false),
		Owned:       gitlabc.Bool(projectOption.Owned),
		Membership:  gitlabc.Bool(projectOption.Membership),
		Starred:     gitlabc.Bool(false),
		Statistics:  gitlabc.Bool(false),
		Visibility:  gitlabc.Visibility("private"),
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

func projectOutput(projects []*gitlabc.Project) []string {
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
