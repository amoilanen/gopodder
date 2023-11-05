package podcasts

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultDownloadDirectory := filepath.Join(userHomeDir, "gopodder")

	fetchCmd.Flags().BoolVarP(&useTitlesAsFileNames, "titles-as-file-names", "t", true, "Specifies that episode titles should be used as file names when downloading episodes")
	fetchCmd.Flags().StringVarP(&downloadDirectory, "download-dir", "d", defaultDownloadDirectory, "Directory where the episodes should be saved")
}

// TODO: Extract downloadFile into a separate package
func downloadFile(url, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP request failed with status code %d", response.StatusCode)
	}
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	contentLength := response.ContentLength

	progressBar := NewProgressBar(contentLength)
	writer := io.MultiWriter(outFile, progressBar)

	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Extract ProgressBar into a separate package
func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{total: total}
}

type ProgressBar struct {
	io.Writer
	total       int64
	bytesCopied int64
}

func (p *ProgressBar) Write(b []byte) (n int, err error) {
	n = len(b)
	p.bytesCopied += int64(n)

	//TODO: Show the download speed
	if p.total > 0 {
		progress := float64(p.bytesCopied) / float64(p.total)
		fmt.Printf("\r[%-80s] %.2f%%", strings.Repeat("#", int(progress*80)), progress*100)
		if p.bytesCopied >= p.total {
			fmt.Printf("\r %-100s", strings.Repeat(" ", 100))
			fmt.Printf("\r")
		}
	} else {
		fmt.Printf("\rDownloaded %d bytes", p.bytesCopied)
	}

	return n, nil
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
		downloadFile(url.String(), filepath.Join(feedPath, name))
	}
	//firstItem := *feed.Items[0]
	//fmt.Println(firstItem)
}
