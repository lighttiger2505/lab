package commands

import (
	"testing"
)

func TestRemoveLineBreak(t *testing.T) {
	got := removeLineBreak("123\r\n456\r789\n")
	want := "123456789"
	if got != want {
		t.Fatalf("want %q, but %q:", want, got)
	}
}
