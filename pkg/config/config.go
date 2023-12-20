package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	MongoDbUrl string `json:"MongoDbUrl"`
}

var configFileName = "./config/config.json"

func LoadConfig() *Config {
	var config Config

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Cannot open config file: ", err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return &config
}
