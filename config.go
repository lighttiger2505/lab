package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Tokens      yaml.MapSlice
	Repositorys []string
	Line        int
	Scope       string
	State       string
	Orderby     string
	Sort        string
}

func NewConfig() (*Config, error) {
	configData, err := getConfigData()
	if err != nil {
		return nil, fmt.Errorf("Failed read config file: %s", err.Error())
	}

	c := Config{}
	err1 := yaml.Unmarshal(configData, &c)
	if err1 != nil {
		return nil, fmt.Errorf("Failed unmarshal yaml: %s", err1.Error())
	}
	return &c, nil
}

func getConfigData() ([]byte, error) {
	filePath := getConfigPath()
	if !fileExists(filePath) {
		err := createConfig(filePath)
		if err != nil {
			return nil, fmt.Errorf("Not exist config: %s", filePath)
		}
	}

	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed read config file: %s", err.Error())
	}

	return configData, nil
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
	c.Repositorys = append(c.Repositorys, repository)
}
