package podcasts

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/amoilanen/gopodder/pkg/config"
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

type FeedEpisodesToFetch struct {
	episodes  []*feed.RssItem
	feedTitle string
}

// TODO: Move some of the more detailed code into a separate module for fetching feeds? FeedFetcher?
func prepareFeedEpisodesFetch(feedUrl string, fromTime time.Time, toTime time.Time) FeedEpisodesToFetch {
	feedReader := feed.FeedReader{}
	parsedFeed, err := feedReader.GetFeed(feedUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	episodesToFetch := []*feed.RssItem{}
	for _, episode := range parsedFeed.Items {
		episodeDate, err := time.Parse(time.RFC1123, episode.PubDate)
		if err != nil {
			panic(err)
		}
		if episodeDate.After(fromTime) && episodeDate.Before(toTime) {
			episodesToFetch = append(episodesToFetch, episode)
		}
	}
	return FeedEpisodesToFetch{episodesToFetch, parsedFeed.Title}
}

func prepareEpisodesFetchForFeeds(feedUrls []string, fromTime time.Time, toTime time.Time) []FeedEpisodesToFetch {
	allEpisodesToFetch := []FeedEpisodesToFetch{}
	for _, feedUrl := range feedUrls {
		allEpisodesToFetch = append(allEpisodesToFetch, prepareFeedEpisodesFetch(feedUrl, fromTime, toTime))
	}
	return allEpisodesToFetch
}

// TODO: Fix fetching the "On Point" podcast feed go run main.go podcasts fetch https://rss.wbur.org/onpoint/podcast.xml --start-time=2023-12-01
func runFetchCmd(ccmd *cobra.Command, args []string) {
	fromTime := parseTime(fromTime)
	toTime := parseTime(toTime)

	allEpisodesToFetch := []FeedEpisodesToFetch{}
	if len(args) > 0 {
		feedUrl := args[0]
		allEpisodesToFetch = prepareEpisodesFetchForFeeds([]string{feedUrl}, fromTime, toTime)
	} else {
		feedUrls := []string{}
		for _, feedConfig := range config.ReadFeedConfigs() {
			feedUrls = append(feedUrls, feedConfig.Feed)
		}
		allEpisodesToFetch = prepareEpisodesFetchForFeeds(feedUrls, fromTime, toTime)
	}

	httpClient := http.HttpClient{}
	for _, feedEpisodesToFetch := range allEpisodesToFetch {
		fmt.Printf("Fetching feed %s...\n", feedEpisodesToFetch.feedTitle)
		feedPath := filepath.Join(downloadDirectory, feedEpisodesToFetch.feedTitle)
		if _, err := os.Stat(feedPath); os.IsNotExist(err) {
			os.MkdirAll(feedPath, 0700)
		}
		for _, episode := range feedEpisodesToFetch.episodes {
			url, err := url.Parse(episode.Enclosure.Url)
			if err != nil {
				println(fmt.Errorf("Not parsable url %s", episode.Enclosure.Url))
			}
			name := episode.Title + path.Ext(url.Path)
			fmt.Printf("Downloading \"%s\": episode \"%s\" published on %s\n", feedEpisodesToFetch.feedTitle, name, episode.PubDate)
			httpClient.DownloadFile(url.String(), filepath.Join(feedPath, name))
		}
	}
	//TODO: Show the overall progress, which episode of how many episodes is being downloaded, etc.
}
