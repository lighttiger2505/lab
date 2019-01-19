package commands

import (
	"bytes"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	gitlab "github.com/xanzy/go-gitlab"
)

type MergeRequestTemplateCommnadOption struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
}

func newMergeRequestTemplateCommandParser(opt *MergeRequestTemplateCommnadOption) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "merge-request-template [options]"
	return parser
}

type MergeRequestTemplateCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	ClientFactory   lab.APIClientFactory
}

func (c *MergeRequestTemplateCommand) Synopsis() string {
	return "List merge request template"
}

func (c *MergeRequestTemplateCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt MergeRequestTemplateCommnadOption
	projectCommandParser := newMergeRequestTemplateCommandParser(&opt)
	projectCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestTemplateCommand) Run(args []string) int {
	var opt MergeRequestTemplateCommnadOption
	projectCommandParser := newMergeRequestTemplateCommandParser(&opt)
	parceArgs, err := projectCommandParser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	pInfo, err := c.RemoteCollecter.CollectTarget(
		opt.ProjectProfileOption.Project,
		opt.ProjectProfileOption.Profile,
	)
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
		filename := MergeRequestTemplateDir + "/" + parceArgs[0]
		res, err := client.GetFile(
			pInfo.Project,
			filename,
			makeShowMergeRequestTemplateOption(),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		c.UI.Message(res)
	} else {
		treeNode, err := client.GetTree(
			pInfo.Project,
			makeMergeRequestTemplateOption(),
		)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
		result := columnize.SimpleFormat(mergeRequestTemplateOutput(treeNode))
		c.UI.Message(result)
	}

	return ExitCodeOK
}

func makeMergeRequestTemplateOption() *gitlab.ListTreeOptions {
	opt := &gitlab.ListTreeOptions{
		Path: gitlab.String(MergeRequestTemplateDir),
		Ref:  gitlab.String("master"),
	}
	return opt
}

func makeShowMergeRequestTemplateOption() *gitlab.GetRawFileOptions {
	opt := &gitlab.GetRawFileOptions{
		Ref: gitlab.String("master"),
	}
	return opt
}

func mergeRequestTemplateOutput(treeNode []*gitlab.TreeNode) []string {
	var outputs []string
	for _, node := range treeNode {
		output := strings.Join([]string{
			node.Name,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
