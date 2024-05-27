package config

import (
	"log"

	"github.com/spf13/viper"
)

func ReadEnvFile(environment string) *viper.Viper {
	config := viper.New()

	config.AddConfigPath(".")
	config.SetConfigType("env")
	config.SetConfigName(environment)

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while parsing configuration file: %v", err)
	}

	return config
}
