package podcasts

import (
	"fmt"
	"os"

	"github.com/amoilanen/gopodder/pkg/feed"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch - fetches the podcast feed provided as an argument",
	Long:  ``,
	Run:   runFetchCmd,
}

func runFetchCmd(ccmd *cobra.Command, args []string) {
	//TODO: Implement fetching the episodes and also checking which episodes have already been fetched
	println("Fetching feed...")
	feedUrl := args[0]
	feedReader := feed.FeedReader{}
	feed, err := feedReader.GetFeed(feedUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	firstItem := *feed.Items[0]
	fmt.Println(firstItem)
}
