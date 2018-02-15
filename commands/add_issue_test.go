package commands

import (
	"testing"
)

func TestCreateMessage(t *testing.T) {
	got := createMessage("title", "description", "#")
	want := `<!-- Creating an issue -->
title

<!-- 
Write a message for this issue. The first block of
text is the title and the rest is the description.
-->
description
`
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}
