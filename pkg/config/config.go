package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	MongoDbUrl string `json:"MongoDbUrl"`
	JwtSecret  string `json:"JwtSecret"`
}

func LoadConfig(configFileName string) *Config {
	var config Config

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Fatal("Cannot open config file: ", err)
		panic(err)
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return &config
}
