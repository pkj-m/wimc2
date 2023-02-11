package config

import (
	"fmt"
	"os"
)

var Config AppConfig

type AppConfig struct {
	Youtube YoutubeConfig `yaml:"youtube"`
	Mongo   MongoConfig   `yaml:"mongo"`
}

type MongoConfig struct {
	ApiKey      string `yaml:"api_key"`
	HostAddress string `yaml:"host_address"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Document    string `yaml:"document"` // mongo document for storing/retrieving data
}

type YoutubeConfig struct {
	BaseURL          string   `yaml:"base_url"`
	ApiKey           string   `yaml:"api_key"`
	ChannelID        string   `yaml:"channel_id"`
	ChannelListParts []string `yaml:"channel_list_parts"`
}

func LoadConfig() error {
	f, err := os.Open("config.yml")
	if err != nil {
		fmt.Println("failed to load config file: %v", err.Error())
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}
}
