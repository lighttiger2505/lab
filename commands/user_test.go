package commands

import (
	"testing"
)

func TestSimple(t *testing.T) {
	got := 1
	want := 2
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}
