package commands

import (
	"bytes"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type IssueTemplateCommnadOption struct {
}

func newIssueTemplateCommandParser(opt *IssueTemplateCommnadOption) *flags.Parser {
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "issue-template [options]"
	return parser
}

type IssueTemplateCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *IssueTemplateCommand) Synopsis() string {
	return "List issue template"
}

func (c *IssueTemplateCommand) Help() string {
	buf := &bytes.Buffer{}
	var projectCommandOption IssueTemplateCommnadOption
	projectCommandParser := newIssueTemplateCommandParser(&projectCommandOption)
	projectCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueTemplateCommand) Run(args []string) int {
	// Parse flags
	var projectCommandOption IssueTemplateCommnadOption
	projectCommandParser := newIssueTemplateCommandParser(&projectCommandOption)
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

	client, err := c.Provider.GetRepositoryClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	treeNode, err := client.GetTree(
		gitlabRemote.RepositoryFullName(),
		makeIssueTemplateOption(),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	result := columnize.SimpleFormat(issueTemplateOutput(treeNode))
	c.UI.Message(result)

	return ExitCodeOK
}

func makeIssueTemplateOption() *gitlab.ListTreeOptions {
	opt := &gitlab.ListTreeOptions{
		Path: gitlab.String(".gitlab/issue_templates"),
		Ref:  gitlab.String("master"),
	}
	return opt
}

func issueTemplateOutput(treeNode []*gitlab.TreeNode) []string {
	var outputs []string
	for _, node := range treeNode {
		output := strings.Join([]string{
			node.Name,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
