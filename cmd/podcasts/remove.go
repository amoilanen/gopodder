package podcasts

import (
	"github.com/amoilanen/gopodder/pkg/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "add - remove a podcast from the list of configured podcasts",
	Long:  "add - remove a podcast from the list of configured podcasts, can remove by providing the name or the feed id",
	Run:   runRemoveCmd,
}

var removeAll bool

func init() {
	removeCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Whether all the saved podcast configurations should be removed")
}

func runRemoveCmd(ccmd *cobra.Command, args []string) {
	var updatedFeedConfigs []config.FeedConfig
	if removeAll {
		updatedFeedConfigs = []config.FeedConfig{}
	} else {
		omitBy := args[0]
		currentFeedConfigs := config.ReadFeedConfigs()
		updatedFeedConfigs = config.OmitFeedConfigsMatching(currentFeedConfigs, omitBy)
	}
	config.WriteFeedConfigs(updatedFeedConfigs)
}
