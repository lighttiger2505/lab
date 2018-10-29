package pipeline

import (
	"fmt"

	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/git"
)

type browseMethod struct {
	opener cmd.URLOpener
	remote *git.RemoteInfo
	id     int
}

func (m *browseMethod) Process() (string, error) {
	var subpage string
	if m.id > 0 {
		subpage = fmt.Sprintf("pipelines/%d", m.id)
	} else {
		subpage = "pipelines"
	}

	url := m.remote.Subpage(subpage)
	if err := m.opener.Open(url); err != nil {
		return "", err
	}

	// Return empty value
	return "", nil
}
