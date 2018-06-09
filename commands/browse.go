package commands

import (
	"bytes"
	"os"
	"path/filepath"

	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/git"
	gitpath "github.com/lighttiger2505/lab/git/path"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

type BrowseCommandOption struct {
	BrowseOption *BrowseOption `group:"Global Options"`
}

func newBrowseOptionParser(opt *BrowseCommandOption) *flags.Parser {
	opt.BrowseOption = newBrowseOption()
	parser := flags.NewParser(opt, flags.Default)
	parser.Usage = "browse [options]"
	return parser
}

type BrowseCommand struct {
	Ui        ui.Ui
	Provider  lab.Provider
	GitClient git.Client
	Opener    cmd.URLOpener
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse repository page"
}

func (c *BrowseCommand) Help() string {
	buf := &bytes.Buffer{}
	var browseCommnadOption BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&browseCommnadOption)
	browseOptionParser.WriteHelp(buf)
	return buf.String()
}

func (c *BrowseCommand) Run(args []string) int {
	var browseCommnadOption BrowseCommandOption
	browseOptionParser := newBrowseOptionParser(&browseCommnadOption)
	// Parse option
	parseArgs, err := browseOptionParser.ParseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Validate option
	browseOption := browseCommnadOption.BrowseOption
	if err := browseOption.IsValid(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Initialize provider
	if err := c.Provider.Init(); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// Getting git remote info
	gitlabRemote, err := c.Provider.GetCurrentRemote()
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	if browseOption.Project != "" {
		namespace, project := browseOption.NameSpaceAndProject()
		gitlabRemote.NameSpace = namespace
		gitlabRemote.Repository = project
	}

	branch, err := c.GitClient.CurrentRemoteBranch(gitlabRemote)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	url, err := c.getURL(parseArgs, gitlabRemote, branch, browseOption)
	if err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	if browseOption.URL {
		c.Ui.Message(url)
		return ExitCodeOK
	}

	if err := c.Opener.Open(url); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}

	// if err := c.doBrowse(url); err != nil {
	// 	c.Ui.Error(err.Error())
	// 	return ExitCodeError
	// }
	return ExitCodeOK
}

func (c *BrowseCommand) getURL(args []string, remote *git.RemoteInfo, branch string, opt *BrowseOption) (string, error) {
	if len(args) > 0 {
		arg := args[0]
		// TODO In order to display an appropriate error message, it is necessary to check whether the argument is a file path
		if isFilePath(arg) {
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
				return remote.BranchFileWithLine(branch, gitAbsPath, opt.Subpage), nil
			}
			return remote.BranchPath(branch, gitAbsPath), nil
		}
	}

	if opt.Subpage != "" {
		return remote.Subpage(opt.Subpage), nil
	}

	// TODO You need to ignore the branch when the project is specified as an option
	if branch == "master" {
		return remote.RepositoryUrl(), nil
	}
	return remote.BranchUrl(branch), nil
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
