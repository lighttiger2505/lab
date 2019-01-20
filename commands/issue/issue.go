package issue

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type IssueCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	MethodFactory   MethodFactory
}

func (c *IssueCommand) Synopsis() string {
	return "Create and Edit, list a issue"
}

func (c *IssueCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt Option
	parser := newOptionParser(&opt)
	parser.WriteHelp(buf)
	return buf.String()
}

func (c *IssueCommand) Run(args []string) int {
	var opt Option
	parser := newOptionParser(&opt)
	parseArgs, err := parser.ParseArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	iid, err := validIssueIID(parseArgs)
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

	clientFacotry, err := lab.NewGitlabClientFactory(pInfo.ApiUrl(), pInfo.Token)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	method := c.MethodFactory.CreateMethod(opt, pInfo, iid, clientFacotry)
	res, err := method.Process()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if res != "" {
		c.UI.Message(res)
	}

	return ExitCodeOK
}

func validIssueIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid args, please intput issue IID.")
	}
	return iid, nil
}

func editIssueMessage(title, description string) string {
	message := `%s

%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}

func editIssueTitleAndDesc(template string, editFunc func(program, file string) error) (string, string, error) {
	editor, err := git.NewEditor("ISSUE", "issue", template, editFunc)
	if err != nil {
		return "", "", err
	}

	title, description, err := editor.EditTitleAndDescription()
	if err != nil {
		return "", "", err
	}

	if editor != nil {
		defer editor.DeleteFile()
	}

	return title, description, nil
}
