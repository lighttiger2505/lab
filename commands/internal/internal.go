package internal

import "regexp"

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
