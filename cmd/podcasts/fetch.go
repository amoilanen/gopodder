package podcasts

import (
	"os"
	"path/filepath"
	"time"

	"github.com/amoilanen/gopodder/pkg/config"
	"github.com/amoilanen/gopodder/pkg/feed"
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

//TODO: Add option last-week, last-month, last-year, today?

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

// TODO: Fix fetching the "On Point" podcast feed go run main.go podcasts fetch https://rss.wbur.org/onpoint/podcast.xml --start-time=2023-12-01
func runFetchCmd(ccmd *cobra.Command, args []string) {
	fromTime := parseTime(fromTime)
	toTime := parseTime(toTime)

	feedFetcher := feed.FeedFetcher{}
	feedUrls := []string{}
	if len(args) > 0 {
		feedUrls = append(feedUrls, args[0])
	} else {
		for _, feedConfig := range config.ReadFeedConfigs() {
			feedUrls = append(feedUrls, feedConfig.Feed)
		}
	}
	feedFetcher.FetchEpisodes(downloadDirectory, fromTime, toTime, feedUrls)
	//TODO: Show the overall progress, which episode of how many episodes is being downloaded, etc.
}
