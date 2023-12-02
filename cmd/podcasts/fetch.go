package podcasts

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/amoilanen/gopodder/pkg/feed"
	"github.com/amoilanen/gopodder/pkg/http"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch - fetches the podcast feed provided as an argument",
	Long:  ``,
	Run:   runFetchCmd,
}

// TODO: Respect this flag and use different file names based on it
var useTitlesAsFileNames bool
var downloadDirectory string

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDownloadDirectory := filepath.Join(userHomeDir, "gopodder")

	fetchCmd.Flags().BoolVarP(&useTitlesAsFileNames, "titles-as-file-names", "t", true, "Specifies that episode titles should be used as file names when downloading episodes")
	fetchCmd.Flags().StringVarP(&downloadDirectory, "download-dir", "d", defaultDownloadDirectory, "Directory where the episodes should be saved")
}

// TODO: Extract downloadFile into a separate package, should it be a part of the

func runFetchCmd(ccmd *cobra.Command, args []string) {
	//TODO: Implement fetching the episodes and also checking which episodes have already been fetched
	println("Fetching feed...")
	feedUrl := args[0]
	feedReader := feed.FeedReader{}
	httpClient := http.HttpClient{}
	feed, err := feedReader.GetFeed(feedUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	feedPath := filepath.Join(downloadDirectory, feed.Title)
	if _, err := os.Stat(feedPath); os.IsNotExist(err) {
		os.MkdirAll(feedPath, 0700)
	}
	//TODO: Show the overall progress, which episode of how many episodes is being downloaded, etc.
	for _, episode := range feed.Items {
		url, err := url.Parse(episode.Enclosure.Url)
		if err != nil {
			println(fmt.Errorf("Not parsable url %s", episode.Enclosure.Url))
		}
		name := episode.Title + path.Ext(url.Path)
		fmt.Printf("Downloading \"%s\": episode \"%s\"\n", feed.Title, name)
		httpClient.DownloadFile(url.String(), filepath.Join(feedPath, name))
	}
	//firstItem := *feed.Items[0]
	//fmt.Println(firstItem)
}
