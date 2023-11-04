package podcasts

import (
	"github.com/amoilanen/gopodder/pkg/config"
	"github.com/amoilanen/gopodder/pkg/feed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add - add a podcast to the list of downloaded podcasts",
	Long:  ``,
	Run:   runAddCmd,
}

func uniqueFeedConfigs(values []config.FeedConfig) []config.FeedConfig {
	uniqueMap := map[string]config.FeedConfig{}
	for _, value := range values {
		uniqueMap[value.Feed] = value
	}
	unique := []config.FeedConfig{}
	for value := range uniqueMap {
		unique = append(unique, uniqueMap[value])
	}
	return unique
}

func createFeedConfig(feedUrl string) config.FeedConfig {
	feedReader := feed.FeedReader{}
	feed, err := feedReader.GetFeed(feedUrl)
	if err != nil {
		panic(err)
	}
	return config.FeedConfig{
		Name: feed.Title,
		Feed: feedUrl,
	}
}

func readFeedConfigs() []config.FeedConfig {
	savedFeedConfigs := []config.FeedConfig{}
	feedData := viper.GetString("feeds")
	err := yaml.Unmarshal([]byte(feedData), &savedFeedConfigs)
	if err != nil {
		panic(err)
	}
	return savedFeedConfigs
}

func writeFeedConfigs(feedConfigs []config.FeedConfig) {
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

func runAddCmd(ccmd *cobra.Command, args []string) {
	feedUrl := args[0]
	newFeedConfig := createFeedConfig(feedUrl)
	currentFeedConfigs := readFeedConfigs()
	updatedFeedConfigs := uniqueFeedConfigs(append(currentFeedConfigs, newFeedConfig))
	writeFeedConfigs(updatedFeedConfigs)
}
