package hoge

import (
	"testing"
)

func TestSimple(t *testing.T) {
	got := 2
	want := 2
	if got != want {
		t.Fatalf("want %v, but %v:", want, got)
	}
}

// 指定したプロジェクトとプロファイルのみを使う
// ロカールにあるリポジトリの情報のみを使う
// プロジェクトは指定したものその他はローカルリポジトリと
