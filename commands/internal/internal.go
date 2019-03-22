package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/git"
	"github.com/lighttiger2505/lab/internal/browse"
	"github.com/lighttiger2505/lab/internal/clipboard"
)

type Method interface {
	Process() (string, error)
}

type BrowseMethod struct {
	Opener    browse.URLOpener
	Clipboard clipboard.Clipboard
	Opt       *BrowseOption
	URL       string
	ID        int
}

func (m *BrowseMethod) Process() (string, error) {
	url := m.URL
	if m.ID > 0 {
		url = strings.Join([]string{url, strconv.Itoa(m.ID)}, "/")
	}

	if m.Opt.Browse {
		if err := m.Opener.Open(url); err != nil {
			return "", err
		}
	}

	if m.Opt.Copy {
		if err := m.Clipboard.Write(url); err != nil {
			return "", err
		}
	}

	if m.Opt.URL {
		return url, nil
	}

	// Return empty value
	return "", nil
}

type MockMethod struct{}

func (m *MockMethod) Process() (string, error) {
	return "result", nil
}

func SweepMarkdownComment(text string) string {
	r := regexp.MustCompile("<!--[\\s\\S]*?-->[\\n]*")
	return r.ReplaceAllString(text, "")
}

func ParceRepositoryFullName(webURL string) string {
	splitURL := strings.Split(webURL, "/")[3:]

	subPageWords := []string{
		"issues",
		"merge_requests",
	}
	var subPageIndex int
	for i, word := range splitURL {
		for _, subPageWord := range subPageWords {
			if word == subPageWord {
				subPageIndex = i
			}
		}
	}

	return strings.Join(splitURL[:subPageIndex], "/")
}

func EditContents(title, description string) string {
	message := `%s

%s
`
	message = fmt.Sprintf(message, title, description)
	return message
}

func EditTitleAndDesc(prefix, template string, editFunc func(program, file string) error) (string, string, error) {
	editor, err := git.NewEditor(prefix, template, editFunc)
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
