package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	flags "github.com/jessevdk/go-flags"
	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/git"
	gitpath "github.com/lighttiger2505/lab/git/path"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
)

type BrowseOption struct {
	Subpage string `short:"s" long:"subpage" description:"open project sub page"`
	URL     bool   `short:"u" long:"url" description:"show project url"`
	Copy    bool   `short:"c" long:"copy" description:"copy project url to clipboard"`
	Project string `short:"p" long:"project" description:"command target specific project"`
}

func newBrowseOption() *BrowseOption {
	browse := flags.NewNamedParser("lab", flags.Default)
	browse.AddGroup("Browse Options", "", &BrowseOption{})
	return &BrowseOption{}
}

func (g *BrowseOption) NameSpaceAndProject() (group, subgroup, project string) {
	splited := strings.Split(g.Project, "/")
	if len(splited) == 3 {
		group = splited[0]
		subgroup = splited[1]
		project = splited[2]
		return
	}
	group = splited[0]
	project = splited[1]
	return
}

type BrowseCommandOption struct {
	BrowseOption *BrowseOption `group:"Global Options"`
}

func newBrowseOptionParser(opt *BrowseCommandOption) *flags.Parser {
	opt.BrowseOption = newBrowseOption()
	parser := flags.NewParser(opt, flags.Default)
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
	Ui        ui.Ui
	Provider  lab.Provider
	GitClient git.Client
	Opener    cmd.URLOpener
}

func (c *BrowseCommand) Synopsis() string {
	return "Browse project page"
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

	browseOption := browseCommnadOption.BrowseOption

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
		group, subgroup, project := browseOption.NameSpaceAndProject()
		gitlabRemote.Group = group
		gitlabRemote.SubGroup = subgroup
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

	if browseOption.Copy {
		if err := clipboard.WriteAll(url); err != nil {
			c.Ui.Error(fmt.Sprintf("Error copying %s to clipboard:\n%s\n", url, err))
		}
		return ExitCodeOK
	}

	if err := c.Opener.Open(url); err != nil {
		c.Ui.Error(err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *BrowseCommand) getURL(args []string, remote *git.RemoteInfo, branch string, opt *BrowseOption) (string, error) {
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
			return remote.BranchFileWithLine(branch, gitAbsPath, opt.Subpage), nil
		}
		return remote.BranchPath(branch, gitAbsPath), nil
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
