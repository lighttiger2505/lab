package cmd

import (
	"testing"
)

type searchBrowserLauncherTest struct {
	goos    string
	browser string
}

var searchBrowserLauncherTests = []searchBrowserLauncherTest{
	{goos: "darwin", browser: "open"},
	{goos: "windows", browser: "cmd /c start"},
}

func TestSearchBrowserLauncher(t *testing.T) {
	for i, test := range searchBrowserLauncherTests {
		browser := searchBrowserLauncher(test.goos)
		if test.browser != browser {
			t.Errorf("#%d: bad return value \nwant %#v \ngot  %#v", i, test.browser, browser)
		}
	}
}
