package cmd

import (
	"fmt"
	"os"

	"github.com/amoilanen/gopodder/cmd/podcasts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(podcasts.PodcastsCmd)
}

var rootCmd = &cobra.Command{
	Use:   "gopodder",
	Short: "gopodder - tool to download podcasts feeds",
	Long:  "gopodder - tool to download podcasts feeds",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
