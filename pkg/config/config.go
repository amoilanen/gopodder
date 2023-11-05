package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type FeedConfig struct {
	Name string `yaml:"name"`
	Feed string `yaml:"feed"`
}

func uniqueFeedConfigs(values []FeedConfig) []FeedConfig {
	uniqueMap := map[string]FeedConfig{}
	for _, value := range values {
		uniqueMap[value.Feed] = value
	}
	unique := []FeedConfig{}
	for value := range uniqueMap {
		unique = append(unique, uniqueMap[value])
	}
	return unique
}

func AppendUniqueFeedConfig(feedConfigs []FeedConfig, feedConfig FeedConfig) []FeedConfig {
	return uniqueFeedConfigs(append(feedConfigs, feedConfig))
}

func OmitFeedConfigsMatching(feedConfigs []FeedConfig, matchBy string) []FeedConfig {
	matchByRegex, err := regexp.Compile(matchBy)
	if err != nil {
		matchByRegex = nil
	}
	result := []FeedConfig{}
	for _, config := range feedConfigs {
		matches := config.Name == matchBy ||
			config.Feed == matchBy ||
			(matchByRegex != nil && len(matchByRegex.FindStringIndex(config.Name)) > 0)
		if !matches {
			result = append(result, config)
		}
	}
	return result
}

func ReadFeedConfigs() []FeedConfig {
	savedFeedConfigs := []FeedConfig{}
	feedData := viper.GetString("feeds")
	err := yaml.Unmarshal([]byte(feedData), &savedFeedConfigs)
	if err != nil {
		panic(err)
	}
	return savedFeedConfigs
}

func WriteFeedConfigs(feedConfigs []FeedConfig) {
	var yamlData []byte
	yamlData, err := yaml.Marshal(&feedConfigs)
	if err != nil {
		panic(err)
	}
	viper.Set("feeds", string(yamlData))
	err = viper.WriteConfig()
	if err != nil {
		panic(err)
	}
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
