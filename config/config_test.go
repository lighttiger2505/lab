package config

import (
	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	// "reflect"
	"testing"
	// "gopkg.in/yaml.v2"
)

// func TestConfig(t *testing.T) {
// 	want := Config{
// 		Tokens: &yaml.MapSlice{
// 			yaml.MapItem{
// 				Key:   "gitlab.ssl.iridge.jp",
// 				Value: "SkLaKmzYDYVsD2bQ2TR",
// 			},
// 		},
// 		Repositorys: []string{"gitlab.ssl.iridge.jp"},
// 		Line:        30,
// 		Scope:       "assigned-to-me",
// 		State:       "closed",
// 		Orderby:     "updated_at",
// 		Sort:        "asc",
// 	}
// 	got, err := NewConfig()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		t.Fail()
// 	}
//
// 	fmt.Println(fmt.Sprintf("%v", got))
// 	if !reflect.DeepEqual(want, got) {
// 		t.Errorf("bad return value want %#v got %#v", want, got)
// 	}
// }

func TestConfig(t *testing.T) {
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err.Error())
	}

	filePath := fmt.Sprintf("%s/labtest.yml", dir)
	err1 := createConfig(filePath)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
}
