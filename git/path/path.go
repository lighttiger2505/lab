package path

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lighttiger2505/lab/git"
)

func Abs(path string) (string, error) {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Can't get os absolute path. %s", err)
	}

	if !isFileExist(fullpath) {
		return "", fmt.Errorf("Not found file or path. Path:%s", fullpath)
	}

	gitroot, err := git.Root()
	if err != nil {
		return "", err
	}

	gitAbsPath := strings.Replace(strings.Replace(fullpath, gitroot, "", -1), "/", "", 1)
	return gitAbsPath, nil
}

func Current() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	gitroot, err := git.Root()
	if err != nil {
		return "", err
	}

	gitAbsPath := strings.Replace(strings.Replace(currentDir, gitroot, "", -1), "/", "", 1)
	return gitAbsPath, nil
}

func isFileExist(fPath string) bool {
	_, err := os.Stat(fPath)
	return err == nil || !os.IsNotExist(err)
}
