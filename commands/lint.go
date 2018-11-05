package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/lighttiger2505/lab/gitlab"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

type LintCommand struct {
	UI            ui.Ui
	Provider      gitlab.Provider
	ClientFactory lab.APIClientFactory
}

func (c *LintCommand) Synopsis() string {
	return "validate .gitlab-ci.yml"
}

func (c *LintCommand) Help() string {
	buf := &bytes.Buffer{}
	buf.WriteString(`Usage:
    lab lint [gitlab-ci.yal file path]`)
	return buf.String()
}

func (c *LintCommand) Run(args []string) int {
	if err := c.Provider.Init(); err != nil {
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
	client := c.ClientFactory.GetLintClient()

	var content string = ""
	if len(args) > 0 {
		b, err := ioutil.ReadFile(args[0])
		content = string(b)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Failed read validate file. \nError: %s", err.Error()))
			return ExitCodeError
		}
	} else {
		c.UI.Error("Required validate file.")
		return ExitCodeError
	}

	result, err := client.Lint(content)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if result.Status == "invalid" {
		for _, msg := range result.Errors {
			c.UI.Message(msg)
			return ExitCodeError
		}
	}

	return ExitCodeOK
}
