package internal

import (
	"regexp"
	"strings"
)

type Method interface {
	Process() (string, error)
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
