package main

import (
	"gopkg.in/yaml.v3"

	"log"
	"os"
	"sync"
)

type Config struct {
	DBDriver string `yaml:"db.driver"`
	DBSource string `yaml:"db.source"`
}

var config *Config
var configOnce sync.Once

func GetConfig() (c Config) {
	configOnce.Do(readConfig)

	return *config
}

func readConfig() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Error reading config.yml: %s", err)
	}

	defer f.Close()

	config = new(Config)

	dec := yaml.NewDecoder(f)
	err = dec.Decode(config)
	if err != nil {
		log.Fatalf("Error parsing config.yml: %s", err)
	}

	return

}
