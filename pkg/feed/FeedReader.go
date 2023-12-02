package feed

import (
	"encoding/xml"

	"github.com/amoilanen/gopodder/pkg/http"
)

type RssItemEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
}

type RssItem struct {
	XMLName   xml.Name         `xml:"item"`
	Title     string           `xml:"title"`
	PubDate   string           `xml:"pubDate"`
	Enclosure RssItemEnclosure `xml:"enclosure"`
}

type RssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Title   string     `xml:"channel>title"`
	Items   []*RssItem `xml:"channel>item"`
}

type FeedReader struct {
}

func (r *FeedReader) GetFeed(feedUrl string) (*RssFeed, error) {
	httpClient := http.HttpClient{}
	respBody, err := httpClient.GetBytesFromUrl(feedUrl)
	if err != nil {
		return nil, err
	}
	var feed RssFeed
	if err := xml.Unmarshal(respBody, &feed); err != nil {
		return nil, err
	}
	return &feed, nil
}
