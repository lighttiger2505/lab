package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
	gitpath "github.com/lighttiger2505/lab/git/path"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/gitutil"
	"github.com/lighttiger2505/lab/internal/ui"
)

type BrowseCommandOption struct {
	ProjectProfileOption *internal.ProjectProfileOption `group:"Project, Profile Options"`
	BrowseOption         *BrowseOption                  `group:"Browse Options"`
}

type BrowseOption struct {
	Subpage string `short:"s" long:"subpage" description:"open project sub page"`
	URL     bool   `short:"u" long:"url" description:"show project url"`
	Copy    bool   `short:"c" long:"copy" description:"copy project url to clipboard"`
}

func newBrowseOptionParser(opt *BrowseCommandOption) *flags.Parser {
	opt.ProjectProfileOption = &internal.ProjectProfileOption{}
	opt.BrowseOption = &BrowseOption{}
	parser := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	parser.Usage = `browse - Browse project page

Synopsis:
  # Browse project page (Show url or Copy url to clipboard by using option)
  lab browse [-u | -c]

  # Browse project file
  lab browse ./README.md

  # Browse issue page
  lab browse -s issues

  # Browse specific project page
  lab browse -p <namespace>/<project>`
	return parser
}

type BrowseCommand struct {
	UI              ui.UI
	RemoteCollecter gitutil.Collecter
	GitClient       git.Client
	Opener          browse.URLOpener
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse project page"
}

func (c *BrowseCommand) Help() string {
	buf := &bytes.Buffer{}
	var opt BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&opt)
	browseOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *BrowseCommand) Run(args []string) int {
	var opt BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&opt)
	parseArgs, err := browseOptionParser.ParseArgs(args)
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

	var branch = "master"
	isGitDir, err := git.IsGitDirReverseTop()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	if isGitDir {
		branch, err = c.GitClient.CurrentRemoteBranch()
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeError
		}
	}

	browseOption := opt.BrowseOption
	url, err := c.getURL(parseArgs, pInfo, branch, browseOption)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}

	if browseOption.URL {
		c.UI.Message(url)
		return ExitCodeOK
	}

	if browseOption.Copy {
		if err := clipboard.WriteAll(url); err != nil {
			c.UI.Error(fmt.Sprintf("Error copying %s to clipboard:\n%s\n", url, err))
		}
		return ExitCodeOK
	}

	if err := c.Opener.Open(url); err != nil {
		c.UI.Error(err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *BrowseCommand) getURL(args []string, pInfo *gitutil.GitLabProjectInfo, branch string, opt *BrowseOption) (string, error) {
	if len(args) > 0 {
		arg := args[0]
		if !isFilePath(arg) {
			return "", fmt.Errorf("Invalid path")
		}

		// TODO In order to display an appropriate error message, it is necessary to check whether the argument is a file path
		// result, err := isDir(arg)
		// if err != nil {
		// 	return "", err
		// }
		// if result {
		// 	gitAbsPath, err := gitpath.Current()
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	return remote.BranchPath(branch, gitAbsPath), nil
		// }

		gitAbsPath, err := gitpath.Abs(arg)
		if err != nil {
			return "", err
		}

		if opt.Subpage != "" {
			return pInfo.BranchFileWithLine(branch, gitAbsPath, opt.Subpage), nil
		}
		return pInfo.BranchPath(branch, gitAbsPath), nil
	}

	if opt.Subpage != "" {
		return pInfo.Subpage(opt.Subpage), nil
	}

	// TODO You need to ignore the branch when the project is specified as an option
	if branch == "master" {
		return pInfo.RepositoryUrl(), nil
	}
	return pInfo.BranchUrl(branch), nil
}

func isFilePath(value string) bool {
	absPath, _ := filepath.Abs(value)
	if isFileExist(absPath) {
		return true
	}
	return false
}

func isDir(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}

	fi, err := file.Stat()
	switch {
	case err != nil:
		return false, err
	case fi.IsDir():
		return true, nil
	default:
		return false, nil
	}
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}
