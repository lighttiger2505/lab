package git

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadTitleAndBody(t *testing.T) {
	message := `A title
A title continues

A body
A body continues
# comment
`
	r := strings.NewReader(message)
	reader := bufio.NewReader(r)
	title, body, err := readTitleAndBody(reader, "#")

	var want string
	if err != nil {
		t.Errorf("except %#v", err)
	}
	want = "A title A title continues"
	if want != title {
		t.Errorf("bad return value want %#v got %#v", want, title)
	}
	want = "A body\nA body continues"
	if want != body {
		t.Errorf("bad return value want %#v got %#v", want, body)
	}

	message = `# Dat title

/ This line is commented out.

Dem body.
`
	r = strings.NewReader(message)
	reader = bufio.NewReader(r)
	title, body, err = readTitleAndBody(reader, "/")

	if err != nil {
		t.Errorf("except %#v", err)
	}
	want = "# Dat title"
	if want != title {
		t.Errorf("bad return value want %#v got %#v", want, title)
	}
	want = "Dem body."
	if want != body {
		t.Errorf("bad return value want %#v got %#v", want, body)
	}
}
