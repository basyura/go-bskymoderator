package config

import (
	"fmt"
	"os"
)

type Config struct {
	UserId    string
	UserDid   string
	Password  string
	ListId    string
	ListAtUri string
	Query     string
}

var _config *Config

func InitializeConfig() (*Config, error) {

	if len(os.Args) != 5 {
		return nil, fmt.Errorf("$ goskymoderator id password uri query")
	}

	conf := &Config{
		UserId:   os.Args[1],
		Password: os.Args[2],
		ListId:   os.Args[3],
		Query:    os.Args[4],
	}

	_config = conf

	return conf, nil
}

func Instance() *Config {
	return _config
}
