package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	Tokens      *yaml.MapSlice
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
	dir, err := homedir.Dir()
	if err != nil {
		return nil, fmt.Errorf("Failed get home dir: %s", err.Error())
	}

	filePath := fmt.Sprintf("%s/.labconfig.yml", dir)
	if !fileExists(filePath) {
		return nil, fmt.Errorf("Not exist config: %s", filePath)
	}

	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed read config file: %s", err.Error())
	}

	return configData, nil
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

func ReadConfig() error {
	// Read config file
	viper.SetConfigName(".labconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.lab")
	if err := viper.ReadInConfig(); err != nil {
		if err := CreateConfig(); err != nil {
			return errors.New(fmt.Sprintf("Failed create config file: %s", err.Error()))
		}

		if err := viper.ReadInConfig(); err != nil {
			return errors.New(fmt.Sprintf("Failed read config file: %s", err.Error()))
		}
	}
	return nil
}

func CreateConfig() error {
	dir, err := homedir.Dir()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed get home dir: %s", err.Error()))
	}

	file, err := os.Create(fmt.Sprintf("%s/.labconfig.yml", dir))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed create config file: %s", err.Error()))
	}
	defer file.Close()

	fmt.Print("Plase input GitLab private token :")
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	_, err = file.Write([]byte(fmt.Sprintf("private_token: %s", stdin.Text())))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed write config file: %s", err.Error()))
	}

	return nil
}

func GetPrivateToken() string {
	return getString("private_token")
}

func GetLine() int {
	return getInt("line")
}

func GetState() string {
	return getString("state")
}

func GetScope() string {
	return getString("scope")
}

func GetOrderBy() string {
	return getString("orderby")
}

func GetSort() string {
	return getString("sort")
}

func getInt(key string) int {
	if viper.InConfig(key) {
		return viper.GetInt(key)
	} else {
		return -1
	}
}

func getString(key string) string {
	if viper.InConfig(key) {
		return viper.GetString(key)
	} else {
		return ""
	}
}
