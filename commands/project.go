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

type ProjectOpt struct {
	Line       int    `short:"n" long:"line" default:"20" default-mask:"20" description:"output the NUM lines"`
	OrderBy    string `short:"o" long:"orderby" default:"updated_at" default-mask:"updated_at" description:"ordered by id, name, path, created_at, updated_at, or last_activity_at fields"`
	Sort       string `short:"s" long:"sort" default:"desc" default-mask:"desc" description:"sorted in asc or desc order"`
	Owned      bool   `short:"w" long:"owned" description:"Limit by projects owned by the current user"`
	Membership bool   `short:"m" long:"member-ship" description:"Limit by projects that the current user is a member of"`
}

var projectOptions ProjectOpt
var projectParser = flags.NewParser(&projectOptions, flags.Default)

type ProjectCommand struct {
	UI       ui.Ui
	Provider gitlab.Provider
}

func (c *ProjectCommand) Synopsis() string {
	return "Show project"
}

func (c *ProjectCommand) Help() string {
	buf := &bytes.Buffer{}
	projectParser.Usage = "project [options]"
	projectParser.WriteHelp(buf)
	return buf.String()
}

func (c *ProjectCommand) Run(args []string) int {
	// Parse flags
	if _, err := projectParser.Parse(); err != nil {
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

	projects, err := client.Projects(makeProjectOptions(projectOptions))
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(projectOutput(projects))
	c.UI.Message(result)

	return ExitCodeOK
}

func makeProjectOptions(opt ProjectOpt) *gitlabc.ListProjectsOptions {
	fmt.Println(opt.OrderBy)
	listProjectsOptions := &gitlabc.ListProjectsOptions{
		Archived:   gitlabc.Bool(false),
		OrderBy:    gitlabc.String(opt.OrderBy),
		Sort:       gitlabc.String(opt.Sort),
		Search:     gitlabc.String(""),
		Simple:     gitlabc.Bool(false),
		Owned:      gitlabc.Bool(opt.Owned),
		Membership: gitlabc.Bool(opt.Membership),
		Starred:    gitlabc.Bool(false),
		Statistics: gitlabc.Bool(false),
		Visibility: gitlabc.Visibility("private"),
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
