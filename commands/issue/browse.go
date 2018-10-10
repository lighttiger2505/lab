package issue

import (
	"fmt"

	"github.com/lighttiger2505/lab/cmd"
	"github.com/lighttiger2505/lab/commands/internal"
	"github.com/lighttiger2505/lab/git"
)

type browseMethod struct {
	internal.Method
	opener cmd.URLOpener
	remote *git.RemoteInfo
	id     int
}

func (m *browseMethod) Process() (string, error) {
	var subpage string
	if m.id > 0 {
		subpage = fmt.Sprintf("issues/%d", m.id)
	} else {
		subpage = "issues"
	}

	url := m.remote.Subpage(subpage)
	if err := m.opener.Open(url); err != nil {
		return "", err
	}

	// Return empty value
	return "", nil
}
