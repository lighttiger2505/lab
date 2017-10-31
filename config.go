package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

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
