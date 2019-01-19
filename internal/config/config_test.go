package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupTestConfig(content string) string {
	tmpfile, err := ioutil.TempFile("", "config")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		panic(err)
	}
	return tmpfile.Name()
}

func getTestConfigContent(fpath string) string {
	file, err := os.OpenFile(fpath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TestConfigV2_Load(t *testing.T) {
	tests := []struct {
		name           string
		configContents string
		want           *Config
		wantErr        bool
	}{
		{
			name:           "empty",
			configContents: "",
			want: &Config{
				Profiles:       map[string]Profile{},
				DefalutProfile: "",
			},
			wantErr: false,
		},
		{
			name: "nomal",
			configContents: `profiles:
  hoge1.com:
    token: token1
    default_group: default_group1
    default_project: default_project1
  hoge2.com:
    token: token2
    default_group: default_group2
    default_project: default_project2
default_profile: default_profile
`,
			want: &Config{
				Profiles: map[string]Profile{
					"hoge1.com": Profile{
						Token:          "token1",
						DefaultGroup:   "default_group1",
						DefaultProject: "default_project1",
					},
					"hoge2.com": Profile{
						Token:          "token2",
						DefaultGroup:   "default_group2",
						DefaultProject: "default_project2",
					},
				},
				DefalutProfile: "default_profile",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configFilePath = setupTestConfig(tt.configContents)
			defer os.Remove(configFilePath)

			c := NewConfig()
			if err := c.Load(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigV2.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := c
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConfigV2.Load() differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func TestConfigV2_Save(t *testing.T) {
	type fields struct {
		Profiles       map[string]Profile
		DefalutProfile string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Profiles:       map[string]Profile{},
				DefalutProfile: "",
			},
			want: `profiles: {}
default_profile: ""
`,
			wantErr: false,
		},
		{
			name: "nomal",
			fields: fields{
				Profiles: map[string]Profile{
					"hoge1.com": Profile{
						Token:          "token1",
						DefaultGroup:   "default_group1",
						DefaultProject: "default_project1",
					},
					"hoge2.com": Profile{
						Token:          "token2",
						DefaultGroup:   "default_group2",
						DefaultProject: "default_project2",
					},
				},
				DefalutProfile: "default_profile",
			},
			want: `profiles:
  hoge1.com:
    token: token1
    default_group: default_group1
    default_project: default_project1
  hoge2.com:
    token: token2
    default_group: default_group2
    default_project: default_project2
default_profile: default_profile
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configFilePath = setupTestConfig("")
			defer os.Remove(configFilePath)

			c := &Config{
				Profiles:       tt.fields.Profiles,
				DefalutProfile: tt.fields.DefalutProfile,
			}

			if err := c.Save(); (err != nil) != tt.wantErr {
				t.Errorf("ConfigV2.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := getTestConfigContent(configFilePath)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("ConfigV2.Save() differs: (-got +want)\n%s", diff)
			}
		})
	}
}

func Test_getXDGConfigPath(t *testing.T) {
	os.Setenv("APPDATA", "appdata")
	os.Setenv("HOME", "home")
	tests := []struct {
		name string
		goos string
		want string
	}{
		{
			name: "windows",
			goos: "windows",
			want: "appdata/lab/config.yml",
		},
		{
			name: "other windows",
			goos: "linux",
			want: "home/.config/lab/config.yml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getXDGConfigPath(tt.goos); got != tt.want {
				t.Errorf("getXDGConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
