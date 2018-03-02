package commands

import (
	"testing"
)

func TestCreateMergeRequestMessage(t *testing.T) {
	got := createMergeRequestMessage("title", "description")
	want := `<!-- Write a message for this merge request. The first block of text is the title -->
title

<!-- the rest is the description.  -->
description
`
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}
