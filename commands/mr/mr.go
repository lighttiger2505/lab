package mr

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/api"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
)

const (
	ExitCodeOK        int = iota //0
	ExitCodeError     int = iota //1
	ExitCodeFileError int = iota //2
)

type MergeRequestCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	GitClient       git.Client
	ClientFactory   api.APIClientFactory
	EditFunc        func(program, file string) error
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Create and Edit, list a merge request"
}

func (c *MergeRequestCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt Option
	mergeRequestCommandParser := newOptionParser(&opt)
	mergeRequestCommandParser.WriteHelp(buf)
	return buf.String()
}

func (c *MergeRequestCommand) Run(args []string) int {
	var opt Option
	mergeRequestCommandParser := newOptionParser(&opt)
	parseArgs, err := mergeRequestCommandParser.ParseArgs(args)
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

	method, err := c.getMethod(opt, parseArgs, pInfo, c.ClientFactory)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

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

func (c *MergeRequestCommand) getMethod(opt Option, args []string, pInfo *gitutil.GitLabProjectInfo, clientFactory api.APIClientFactory) (internal.Method, error) {
	createUpdateOption := opt.CreateUpdateOption
	listOption := opt.ListOption
	browseOption := opt.BrowseOption
	showOption := opt.ShowOption

	mrClient := clientFactory.GetMergeRequestClient()
	repositoryClient := clientFactory.GetRepositoryClient()
	noteClient := clientFactory.GetNoteClient()

	iid, err := validMergeRequestIID(args)
	if err != nil {
		return nil, err
	}

	if browseOption.HasBrowse() {
		return &internal.BrowseMethod{
			Opener: &browse.Browser{},
			Opt:    browseOption,
			URL:    pInfo.SubpageUrl("merge_requests"),
			ID:     iid,
		}, nil
	}

	// Case of getting Merge Request id
	if len(args) > 0 {
		if createUpdateOption.hasEdit() {
			return &updateOnEditorMethod{
				client:   mrClient,
				opt:      createUpdateOption,
				project:  pInfo.Project,
				id:       iid,
				editFunc: c.EditFunc,
			}, nil
		}
		if createUpdateOption.hasUpdate() {
			return &updateMethod{
				client:  mrClient,
				opt:     createUpdateOption,
				project: pInfo.Project,
				id:      iid,
			}, nil
		}

		return &detailMethod{
			mrClient:   mrClient,
			noteClient: noteClient,
			opt:        showOption,
			project:    pInfo.Project,
			id:         iid,
		}, nil
	}

	// Case of nothing MergeRequest id
	if createUpdateOption.hasEdit() {
		return &createOnEditorMethod{
			client:           mrClient,
			repositoryClient: repositoryClient,
			opt:              createUpdateOption,
			pInfo:            pInfo,
			editFunc:         c.EditFunc,
		}, nil

	}
	if createUpdateOption.hasCreate() {
		return &createMethod{
			client: mrClient,
			opt:    createUpdateOption,
			pInfo:  pInfo,
		}, nil
	}

	if listOption.AllProject {
		return &listAllMethod{
			client: mrClient,
			opt:    listOption,
		}, nil

	}
	return &listMethod{
		client:  mrClient,
		opt:     listOption,
		project: pInfo.Project,
	}, nil
}

func validMergeRequestIID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, nil
	}

	iid, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("Invalid args, please input merge request id")
	}
	return iid, nil
}
