package commands

import (
	"bytes"
	"strconv"
	"strings"

	flags "github.com/jessevdk/go-flags"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	"github.com/ryanuber/columnize"
	"github.com/xanzy/go-gitlab"
)

type LabelCommnadOption struct {
}

func newLabelOptionParser(opt *LabelCommnadOption) *flags.Parser {
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "label [options]"
	return parser
}

type LabelCommand struct {
	UI       ui.Ui
	Provider lab.Provider
}

func (c *LabelCommand) Synopsis() string {
	return "List label"
}

func (c *LabelCommand) Help() string {
	var labelCommnadOption LabelCommnadOption
	labelCommnadOptionParser := newLabelOptionParser(&labelCommnadOption)
	buf := &bytes.Buffer{}
	labelCommnadOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *LabelCommand) Run(args []string) int {
	// Parse flags
	var labelCommnadOption LabelCommnadOption
	labelCommnadOptionParser := newLabelOptionParser(&labelCommnadOption)
	if _, err := labelCommnadOptionParser.ParseArgs(args); err != nil {
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

	client, err := c.Provider.GetLabelClient(gitlabRemote)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	labels, err := client.GetLabels(
		gitlabRemote.RepositoryFullName(),
		makeListLabelsOptions(),
	)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	result := columnize.SimpleFormat(labelOutput(labels))
	c.UI.Message(result)

	return ExitCodeOK
}

func makeListLabelsOptions() *gitlab.ListLabelsOptions {
	return &gitlab.ListLabelsOptions{}
}

func labelOutput(labels []*gitlab.Label) []string {
	var outputs []string
	for _, label := range labels {
		output := strings.Join([]string{
			strconv.Itoa(label.ID),
			label.Name,
			label.Description,
		}, "|")
		outputs = append(outputs, output)
	}
	return outputs
}
