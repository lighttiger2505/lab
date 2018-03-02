package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestConfigManagerLoad(t *testing.T) {
	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	defer os.Remove(tmppath)

	configData := ConfigDataTest
	f.Write([]byte(configData))

	conf := NewConfigManager()
	conf.Path = tmppath
	c, err := conf.Load()
	if err != nil {
		t.Fatalf("wrong error. errors: \n%s", err.Error())
	}

	wantTokens := TokensTest
	wantPreferredDomains := PreferredDomainTest
	if !reflect.DeepEqual(wantTokens, c.Tokens) || !reflect.DeepEqual(wantPreferredDomains, c.PreferredDomains) {
		t.Fatalf("bad return value \nwant %#v %#v \ngot  %#v %#v", wantTokens, wantPreferredDomains, c.Tokens, c.PreferredDomains)
	}
}

func TestConfigManagerRead(t *testing.T) {
	configData := ConfigDataTest
	r := strings.NewReader(configData)

	conf := NewConfigManager()
	c, err := conf.read(r)
	if err != nil {
		t.Fatalf("wrong error. errors: \n%s", err.Error())
	}

	wantTokens := TokensTest
	wantPreferredDomains := PreferredDomainTest
	if !reflect.DeepEqual(wantTokens, c.Tokens) || !reflect.DeepEqual(wantPreferredDomains, c.PreferredDomains) {
		t.Fatalf("bad return value \nwant %#v %#v \ngot  %#v %#v", wantTokens, wantPreferredDomains, c.Tokens, c.PreferredDomains)
	}
}

func TestConfigManagerSave(t *testing.T) {
	f, _ := ioutil.TempFile("", "test")
	tmppath := f.Name()
	f.Close()
	defer os.Remove(tmppath)

	conf := NewConfigManagerPath(tmppath)
	conf.Config = &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	if err := conf.Save(); err != nil {
		t.Fatalf("wrong error. errors: \n%s", err.Error())
	}

	read, _ := ioutil.ReadFile(tmppath)
	got := string(read)
	want := ConfigDataTest
	if want != got {
		t.Fatalf("bad write value \nwant %s \ngot  %s", want, got)
	}
}

func TestConfigManagerWrite(t *testing.T) {
	w := bytes.NewBufferString("")
	conf := NewConfigManager()
	conf.Config = &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}

	err := conf.write(w)
	if err != nil {
		t.Fatalf("wrong error. errors: \n%s", err.Error())
	}

	want := ConfigDataTest
	got := w.String()
	if want != got {
		t.Fatalf("bad return value \nwant %s \ngot  %s", want, got)
	}
}

func TestConfigGetToken(t *testing.T) {
	conf := &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	got := conf.getToken("gitlab.ssl.domain1.jp")
	want := "token1"
	if want != got {
		t.Fatalf("bad return value \nwant %s \ngot  %s", want, got)
	}
}

func TestConfigHasDomain(t *testing.T) {
	conf := &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	got := conf.hasDomain("gitlab.ssl.domain2.jp")
	want := true
	if want != got {
		t.Fatalf("bad return value \nwant %v \ngot  %v", want, got)
	}

	got = conf.hasDomain("unknown")
	want = false
	if want != got {
		t.Fatalf("bad return value \nwant %v \ngot  %v", want, got)
	}
}

func TestConfigGetTopDomain(t *testing.T) {
	conf := &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	got := conf.getTopDomain()
	want := "gitlab.ssl.domain1.jp"
	if want != got {
		t.Fatalf("bad return value \nwant %s \ngot  %s", want, got)
	}
}

func TestConfigAddToken(t *testing.T) {
	conf := &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	conf.addToken("NewToken", "NewDomain")
	got := conf.Tokens
	want := yaml.MapSlice{
		yaml.MapItem{
			Key:   "gitlab.ssl.domain1.jp",
			Value: "token1",
		},
		yaml.MapItem{
			Key:   "gitlab.ssl.domain2.jp",
			Value: "token2",
		},
		yaml.MapItem{
			Key:   "NewDomain",
			Value: "NewToken",
		},
	}
	if reflect.DeepEqual(want, got) {
		t.Fatalf("bad return value \nwant %v \ngot  %v", want, got)
	}
}

func TestConfigAddRepository(t *testing.T) {

	conf := &Config{
		Tokens:           TokensTest,
		PreferredDomains: PreferredDomainTest,
	}
	conf.AddRepository("NewDomain")
	got := conf.PreferredDomains
	want := []string{
		"gitlab.ssl.domain1.jp",
		"gitlab.ssl.domain2.jp",
	}
	if reflect.DeepEqual(want, got) {
		t.Fatalf("bad return value \nwant %v \ngot  %v", want, got)
	}
}
