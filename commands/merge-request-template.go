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

type MergeRequestTemplateCommnadOption struct {
}

func newMergeRequestTemplateCommandParser(opt *MergeRequestTemplateCommnadOption) *flags.Parser {
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = "merge-request-template [options]"
	return parser
}

type MergeRequestTemplateCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *MergeRequestTemplateCommand) Synopsis() string {
	return "List merge request template"
}

func (c *MergeRequestTemplateCommand) Help() string {
	buf := &bytes.Buffer{}
	var projectCommandOption MergeRequestTemplateCommnadOption
	projectCommandParser := newMergeRequestTemplateCommandParser(&projectCommandOption)
	projectCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestTemplateCommand) Run(args []string) int {
	// Parse flags
	var projectCommandOption MergeRequestTemplateCommnadOption
	projectCommandParser := newMergeRequestTemplateCommandParser(&projectCommandOption)
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
		filename := MergeRequestTemplateDir + "/" + parceArgs[0]
		res, err := client.GetFile(
			gitlabRemote.RepositoryFullName(),
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
			gitlabRemote.RepositoryFullName(),
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
