package mr

import (
	"strconv"
	"strings"

	"github.com/lighttiger2505/lab/internal/browse"
)

type browseMethod struct {
	opener browse.URLOpener
	url    string
	id     int
}

func (m *browseMethod) Process() (string, error) {
	url := m.url
	if m.id > 0 {
		url = strings.Join([]string{url, strconv.Itoa(m.id)}, "/")
	}

	if err := m.opener.Open(url); err != nil {
		return "", err
	}

	// Return empty value
	return "", nil
}
