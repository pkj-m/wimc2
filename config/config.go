package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type AppConfig struct {
	Youtube YoutubeConfig `yaml:"youtube"`
	Mongo   MongoConfig   `yaml:"mongo"`
}

type MongoConfig struct {
	ApiKey           string `yaml:"api_key"`
	HostAddress      string `yaml:"host_address"`
	Port             int    `yaml:"port"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	Database         string `yaml:"database"`          // mongo database for storing/retrieving documents
	CollectionPrefix string `yaml:"collection_prefix"` // prefix for every collection name, appended with date of search
}

type YoutubeConfig struct {
	BaseURL          string   `yaml:"base_url"`
	ApiKey           string   `yaml:"api_key"`
	ChannelID        string   `yaml:"channel_id"`
	ChannelListParts []string `yaml:"channel_list_parts"`
	MaxSearchResults int      `yaml:"max_search_results"`
}

func LoadConfig() (*AppConfig, error) {
	f, err := os.Open("./config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config AppConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
