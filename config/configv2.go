package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"
)

var configFilePath = getXDGConfigPath(runtime.GOOS)

type ConfigV2 struct {
	Profiles       map[string]Profile `yaml:"profiles"`
	DefalutProfile string             `yaml:"default_profile"`
}

type Profile struct {
	Token          string `yaml:"token"`
	DefaultGroup   string `yaml:"default_group"`
	DefaultProject string `yaml:"default_project"`
}

func NewConfig() *ConfigV2 {
	cfg := &ConfigV2{
		Profiles: map[string]Profile{},
	}
	return cfg
}

func GetConfig() (*ConfigV2, error) {
	cfg := NewConfig()
	if err := cfg.Load(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *ConfigV2) Load() error {
	if err := os.MkdirAll(filepath.Dir(configFilePath), 0700); err != nil {
		return fmt.Errorf("cannot create directory, %s", err)
	}

	if !fileExists(configFilePath) {
		_, err := os.Create(configFilePath)
		if err != nil {
			return fmt.Errorf("cannot create config, %s", err.Error())
		}
	}

	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot open config, %s", err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read config, %s", err)
	}

	if err = yaml.Unmarshal(b, c); err != nil {
		return fmt.Errorf("failed unmarshal yaml. \nError: %s \nBuffer: %s", err, string(b))
	}
	return nil
}

func (c *ConfigV2) Save() error {
	file, err := os.OpenFile(configFilePath, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("cannot open file, %s", err)
	}
	defer file.Close()

	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("Failed marshal config. Error: %v", err)
	}

	if _, err = io.WriteString(file, string(out)); err != nil {
		return fmt.Errorf("Failed write config file. Error: %s", err)
	}
	return nil
}

func (c *ConfigV2) GetProfile(domain string) (*Profile, error) {
	profile, ok := c.Profiles[domain]
	if !ok {
		return nil, fmt.Errorf("not found profile, [%s]. Please check config", domain)
	}
	return &profile, nil
}

func (c *ConfigV2) GetDefaultProfile() *Profile {
	profile, _ := c.GetProfile(c.DefalutProfile)
	return profile
}

func (c *ConfigV2) SetProfile(domain string, profile Profile) {
	c.Profiles[domain] = profile
}

func (c *ConfigV2) HasDomain(domain string) bool {
	_, ok := c.Profiles[domain]
	if !ok {
		return false
	}
	return true
}

func (c *ConfigV2) GetToken(domain string) string {
	profile, _ := c.GetProfile(domain)
	return profile.Token
}

func (c *ConfigV2) SetToken(domain, token string) {
	profile, _ := c.GetProfile(domain)
	profile.Token = token
	c.SetProfile(domain, *profile)
}

func getXDGConfigPath(goos string) string {
	var dir string
	if goos == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "lab")
		}
		dir = filepath.Join(dir, "lab")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "lab")
	}
	return filepath.Join(dir, "config.yml")
}
