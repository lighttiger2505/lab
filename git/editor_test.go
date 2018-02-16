package git

import (
	"bufio"
	"strings"
	"testing"
)

func TestSweepMarkdownCommnet(t *testing.T) {
	text := `A title
<!-- A comment -->
A title continues

<!-- A comment -->
A body
<!--
A comment continues -->
<!-- A comment continues
-->
A body continues
<!--
A comment continues
-->
`
	want := `A title
A title continues

A body
A body continues
`
	got := sweepMarkdownComment(text)
	if want != got {
		t.Errorf("bad return value want %#v got %#v", want, got)
	}
}

func TestParceTitleAndBody(t *testing.T) {
	text := `A title
A title continues

A body
A body continues
`
	title, body := parceTitleAndBody(text)
	var want string
	want = "A title A title continues"
	if want != title {
		t.Errorf("bad return value want %#v got %#v", want, title)
	}
	want = "A body\nA body continues"
	if want != body {
		t.Errorf("bad return value want %#v got %#v", want, body)
	}
}

func TestReadTitleAndBody(t *testing.T) {
	message := `A title
A title continues

A body
A body continues
<!-- comment -->
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
}
