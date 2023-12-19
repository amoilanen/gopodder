package feed

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/amoilanen/gopodder/pkg/http"
	timeutils "github.com/amoilanen/gopodder/pkg/time"
)

type FeedFetcher struct {
}

type FeedEpisodesToFetch struct {
	Episodes  []*RssItem
	FeedTitle string
}

func (f *FeedFetcher) prepareFeedEpisodesFetch(feedUrl string, fromTime time.Time, toTime time.Time) FeedEpisodesToFetch {
	feedReader := FeedReader{}
	parsedFeed, err := feedReader.GetFeed(feedUrl)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	episodesToFetch := []*RssItem{}
	for _, episode := range parsedFeed.Items {
		episodeDate := timeutils.ParseTime(episode.PubDate)
		if episodeDate.After(fromTime) && episodeDate.Before(toTime) {
			episodesToFetch = append(episodesToFetch, episode)
		}
	}
	return FeedEpisodesToFetch{episodesToFetch, parsedFeed.Title}
}

func (f *FeedFetcher) prepareEpisodesFetchForFeeds(feedUrls []string, fromTime time.Time, toTime time.Time) []FeedEpisodesToFetch {
	allEpisodesToFetch := []FeedEpisodesToFetch{}
	for _, feedUrl := range feedUrls {
		allEpisodesToFetch = append(allEpisodesToFetch, f.prepareFeedEpisodesFetch(feedUrl, fromTime, toTime))
	}
	return allEpisodesToFetch
}

func (f *FeedFetcher) fetch(downloadDirectory string, feedEpisodes []FeedEpisodesToFetch) {
	httpClient := http.HttpClient{}
	for _, feedEpisodesToFetch := range feedEpisodes {
		fmt.Printf("Fetching feed %s...\n", feedEpisodesToFetch.FeedTitle)
		feedPath := filepath.Join(downloadDirectory, feedEpisodesToFetch.FeedTitle)
		if _, err := os.Stat(feedPath); os.IsNotExist(err) {
			os.MkdirAll(feedPath, 0700)
		}
		for _, episode := range feedEpisodesToFetch.Episodes {
			url, err := url.Parse(episode.Enclosure.Url)
			if err != nil {
				println(fmt.Errorf("Not parsable url %s", episode.Enclosure.Url))
			}
			name := strings.TrimSpace(episode.Title + path.Ext(url.Path))
			fmt.Printf("Downloading \"%s\": episode \"%s\" published on %s\n", feedEpisodesToFetch.FeedTitle, name, episode.PubDate)
			httpClient.DownloadFile(url.String(), filepath.Join(feedPath, name))
		}
	}
}

func (f *FeedFetcher) FetchEpisodes(downloadDirectory string, fromTime time.Time, toTime time.Time, feedUrls []string) {
	allEpisodesToFetch := f.prepareEpisodesFetchForFeeds(feedUrls, fromTime, toTime)
	f.fetch(downloadDirectory, allEpisodesToFetch)
}
