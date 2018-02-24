package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/mitchellh/go-homedir"
)

type ConfigManager struct {
	Path   string
	Config *Config
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		Path:   "",
		Config: nil,
	}
}

func (c *ConfigManager) Init() error {
	filepath := c.Path
	if filepath == "" {
		filepath := getConfigPath()
		if !fileExists(filepath) {
			if err := createConfig(filepath); err != nil {
				return fmt.Errorf("Not exist config: %s", filepath)
			}
		}
		c.Path = filepath
	}
	return nil
}

func (c *ConfigManager) Load() (*Config, error) {
	if !fileExists(c.Path) {
		return nil, fmt.Errorf("Not exist config: %s", c.Path)
	}

	configData, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed read config file: %s", err.Error())
	}

	conf := Config{}
	if err := yaml.Unmarshal(configData, &conf); err != nil {
		return nil, fmt.Errorf("Failed unmarshal yaml: %s", err.Error())
	}
	c.Config = &conf
	return &conf, nil
}

type Config struct {
	Tokens           yaml.MapSlice
	PreferredDomains []string
}

func NewConfig() (*Config, error) {
	filepath := getConfigPath()
	if !fileExists(filepath) {
		err := createConfig(filepath)
		if err != nil {
			return nil, fmt.Errorf("Not exist config: %s", filepath)
		}
	}
	return NewConfigWithFile(filepath)
}

func NewConfigWithFile(filepath string) (*Config, error) {
	if !fileExists(filepath) {
		return nil, fmt.Errorf("Not exist config: %s", filepath)
	}

	configData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Failed read config file: %s", err.Error())
	}

	c := Config{}
	if err := yaml.Unmarshal(configData, &c); err != nil {
		return nil, fmt.Errorf("Failed unmarshal yaml: %s", err.Error())
	}
	return &c, nil
}

func getConfigPath() string {
	dir, _ := homedir.Dir()
	filePath := fmt.Sprintf("%s/.labconfig.yml", dir)
	return filePath
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	if pathError, ok := err.(*os.PathError); ok {
		if pathError.Err == syscall.ENOTDIR {
			return false
		}
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func createConfig(filePath string) error {
	config := Config{}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed create config file: %s", err.Error())
	}
	defer file.Close()

	out, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("Failed marshal config: %v", err.Error())
	}

	_, err = file.Write(out)
	if err != nil {
		return fmt.Errorf("Failed write config file: %s", err.Error())
	}

	return nil
}

func (c *Config) Write() error {
	filePath := getConfigPath()
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("Failed open config file: %s", err.Error())
	}
	defer file.Close()

	out, err := yaml.Marshal(&c)
	if err != nil {
		return fmt.Errorf("Failed marshal config: %v", err.Error())
	}

	_, err = file.Write(out)
	if err != nil {
		return fmt.Errorf("Failed write config file: %s", err.Error())
	}

	return nil
}

func (c *Config) AddToken(domain string, token string) {
	item := yaml.MapItem{
		Key:   domain,
		Value: token,
	}
	c.Tokens = append(c.Tokens, item)
}

func (c *Config) AddRepository(repository string) {
	c.PreferredDomains = append(c.PreferredDomains, repository)
}

func (c *Config) MustDomain() string {
	if len(c.PreferredDomains) > 0 {
		return c.PreferredDomains[0]
	}
	return ""
}
