package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	RethinkDB `json:"rethinkdb"`
}

type RethinkDB struct {
	Address  string `json:"address"`
	Database string `json:"database"`
}

func LoadConfig() Config {
	config := Config{
		RethinkDB: RethinkDB{},
	}

	configPath, isExist := os.LookupEnv("CONFIG_PATH")
	if !isExist {
		panic("Please set the environment variable!: CONFIG_PATH")
	}

	jsonBytes, _ := os.ReadFile(configPath)
	err := json.Unmarshal(jsonBytes, &config)
	if err != nil {
		panic("Cannot load config file!")
	}

	return config
}
