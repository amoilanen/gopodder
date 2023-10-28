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

var useTitlesAsFileNames bool
var downloadDirectory string

func init() {
	fetchCmd.Flags().BoolVarP(&useTitlesAsFileNames, "titles-as-file-names", "t", false, "Specifies that episode titles should be used as file names when downloading episodes")
	fetchCmd.Flags().StringVarP(&downloadDirectory, "download-dir", "d", "~/.gopodder", "Directory where the episodes should be saved")
}

func runFetchCmd(ccmd *cobra.Command, args []string) {
	//TODO: Implement fetching the episodes and also checking which episodes have already been fetched
	println("Fetching feed...")
	println(fmt.Sprintf("use titles as file names = %t, download directory = %s", useTitlesAsFileNames, downloadDirectory))
	//TODO: Create downloadDirectory if it does not exist
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
