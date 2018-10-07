package internal

import "regexp"

type Method interface {
	Process() (string, error)
}

func SweepMarkdownComment(text string) string {
	r := regexp.MustCompile("<!--[\\s\\S]*?-->[\\n]*")
	return r.ReplaceAllString(text, "")
}
