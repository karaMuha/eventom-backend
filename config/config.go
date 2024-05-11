package config

import (
	"crypto/rsa"
	"log"

	"github.com/spf13/viper"
)

func ReadEnvFile(filename string) *viper.Viper {
	config := viper.New()

	config.AddConfigPath(".")
	config.SetConfigType("env")
	config.SetConfigName(filename)

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while parsing configuration file: %v", err)
	}

	return config
}

var PrivateKey *rsa.PrivateKey
