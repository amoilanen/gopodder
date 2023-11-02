package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type FeedConfig struct {
	Name string `yaml:"name"`
	Feed string `yaml:"feed"`
}

func InitViper() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("User home directory is not known! %w", err))
	}
	configPath := filepath.Join(userHomeDir, ".gopodder")
	configFile := filepath.Join(configPath, "config.yaml")

	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, 0700)
	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		viper.SetDefault("feeds", []FeedConfig{})
		err = viper.SafeWriteConfigAs(configFile)
		if err != nil {
			panic(fmt.Errorf("Could not write default config file: %w", err))
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Could not read default config file: %w", err))
	}
}
