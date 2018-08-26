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

const ISSUE_TEMPLATE = ".gitlab/issue_templates"

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
	parceArgs, err := projectCommandParser.ParseArgs(args)
	if err != nil {
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

	if len(parceArgs) > 0 {
		filename := ISSUE_TEMPLATE + "/" + parceArgs[0]
		res, err := client.GetFile(
			gitlabRemote.RepositoryFullName(),
			filename,
			makeShowIssueTemplateOption(),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		c.UI.Message(res)
	} else {
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
	}

	return ExitCodeOK
}

func makeIssueTemplateOption() *gitlab.ListTreeOptions {
	opt := &gitlab.ListTreeOptions{
		Path: gitlab.String(ISSUE_TEMPLATE),
		Ref:  gitlab.String("master"),
	}
	return opt
}

func makeShowIssueTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
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
