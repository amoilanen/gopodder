package podcasts

import (
	"github.com/amoilanen/gopodder/pkg/config"
	"github.com/amoilanen/gopodder/pkg/feed"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add - add a podcast to the list of downloaded podcasts",
	Long:  ``,
	Run:   runAddCmd,
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

func runAddCmd(ccmd *cobra.Command, args []string) {
	feedUrl := args[0]
	newFeedConfig := createFeedConfig(feedUrl)
	currentFeedConfigs := config.ReadFeedConfigs()
	updatedFeedConfigs := config.AppendUniqueFeedConfig(currentFeedConfigs, newFeedConfig)
	config.WriteFeedConfigs(updatedFeedConfigs)
}
