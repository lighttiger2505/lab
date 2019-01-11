package commands

import (
	"bytes"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type IssueTemplateCommnadOption struct {
}

func newIssueTemplateCommandParser(opt *IssueTemplateCommnadOption) *flags.Parser {
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "issue-template [options]"
	return parser
}

type IssueTemplateCommand struct {
	UI              ui.Ui
	RemoteCollecter gitutil.Collecter
	ClientFactory   lab.APIClientFactory
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
	var projectCommandOption IssueTemplateCommnadOption
	projectCommandParser := newIssueTemplateCommandParser(&projectCommandOption)
	parceArgs, err := projectCommandParser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pInfo, err := c.RemoteCollecter.CollectTarget("", "")
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if err := c.ClientFactory.Init(pInfo.ApiUrl(), pInfo.Token); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	client := c.ClientFactory.GetRepositoryClient()

	if len(parceArgs) > 0 {
		filename := IssueTemplateDir + "/" + parceArgs[0]
		res, err := client.GetFile(
			pInfo.Project,
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
			pInfo.Project,
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
		Path: gitlab.String(IssueTemplateDir),
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
