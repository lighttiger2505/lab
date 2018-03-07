package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlabc "github.com/xanzy/go-gitlab"
)

type ProjectCommand struct {
	UI       ui.Ui
	Provider gitlab.Provider
}

func (c *ProjectCommand) Synopsis() string {
	return "Show project"
}

func (c *ProjectCommand) Help() string {
	buf := &bytes.Buffer{}
	return buf.String()
}

func (c *ProjectCommand) Run(args []string) int {
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

	opt := &gitlabc.ListProjectsOptions{
		Archived:   gitlabc.Bool(false),
		OrderBy:    gitlabc.String("name"),
		Sort:       gitlabc.String("desc"),
		Search:     gitlabc.String(""),
		Simple:     gitlabc.Bool(false),
		Owned:      gitlabc.Bool(false),
		Membership: gitlabc.Bool(false),
		Starred:    gitlabc.Bool(false),
		Statistics: gitlabc.Bool(false),
		Visibility: gitlabc.Visibility("private"),
	}

	projects, err := client.Projects(opt)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	var datas []string
	for _, project := range projects {
		data := strings.Join([]string{
			fmt.Sprintf("%s/%s", project.Namespace.Name, project.Name),
			removeLineBreak(project.Description),
		}, "|")
		datas = append(datas, data)
	}

	result := columnize.SimpleFormat(datas)
	c.UI.Message(result)

	return ExitCodeOK
}

func removeLineBreak(value string) string {
	value = strings.Replace(value, "\r\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)
	value = strings.Replace(value, "\n", "", -1)
	return value
}
