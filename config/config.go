package config

import (
	"os"
)

type Config struct {
	Port string
}

func GetConfig() *Config {
	var config Config
	port, ok := os.LookupEnv("PORT")
	if ok {
		config.Port = port
	} else {
		config.Port = "3000"
	}

	return &config
}
