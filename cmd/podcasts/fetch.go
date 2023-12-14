package podcasts

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

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
var fromTime string
var toTime string

func tryParseWithLayouts(timeInput string, formats ...string) (time.Time, error) {
	var parsed time.Time
	var err error
	for _, layout := range formats {
		parsed, err = time.Parse(layout, timeInput)
		if err == nil {
			return parsed, nil
		}
	}
	return parsed, err
}

func currentTime() string {
	currentTime := time.Now()
	return currentTime.Format(time.RFC3339)
}

func parseTime(timeInput string) time.Time {
	parsedDate, err := tryParseWithLayouts(timeInput, time.RFC3339, time.DateOnly, time.RFC1123)
	if err != nil {
		panic(err)
	}
	return parsedDate
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDownloadDirectory := filepath.Join(userHomeDir, "gopodder")
	defaultFromTime := "1970-01-01T00:00:00.000Z"
	defaultToTime := currentTime()

	fetchCmd.Flags().BoolVarP(&useTitlesAsFileNames, "titles-as-file-names", "t", true, "Specifies that episode titles should be used as file names when downloading episodes")
	fetchCmd.Flags().StringVarP(&downloadDirectory, "download-dir", "d", defaultDownloadDirectory, "Directory where the episodes should be saved")
	fetchCmd.Flags().StringVarP(&fromTime, "start-time", "s", defaultFromTime, "Start timestamp after which the episodes should be fetched")
	fetchCmd.Flags().StringVarP(&toTime, "end-time", "e", defaultToTime, "End timestamp before which the episodes should be fetched")
}

// TODO: Implement fetching all the configured feeds if no argument is provided
func runFetchCmd(ccmd *cobra.Command, args []string) {
	afterTime := parseTime(fromTime)
	fmt.Printf("Downloading after %s\n", afterTime.Format(time.RFC3339))
	untilTime := parseTime(toTime)
	fmt.Printf("Downloading before %s\n", untilTime.Format(time.RFC3339))
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
		episodeDate, err := time.Parse(time.RFC1123, episode.PubDate)
		if err != nil {
			panic(err)
		}
		if episodeDate.After(afterTime) && episodeDate.Before(untilTime) {
			url, err := url.Parse(episode.Enclosure.Url)
			if err != nil {
				println(fmt.Errorf("Not parsable url %s", episode.Enclosure.Url))
			}
			name := episode.Title + path.Ext(url.Path)
			fmt.Printf("Downloading \"%s\": episode \"%s\" published on %s\n", feed.Title, name, episode.PubDate)
			httpClient.DownloadFile(url.String(), filepath.Join(feedPath, name))
		}
	}
	//firstItem := *feed.Items[0]
	//fmt.Println(firstItem)
}
