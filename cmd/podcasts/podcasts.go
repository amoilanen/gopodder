package podcasts

import (
	"github.com/spf13/cobra"
)

func init() {
	PodcastsCmd.AddCommand(fetchCmd)
}

var PodcastsCmd = &cobra.Command{
	Use:   "podcasts",
	Short: "podcasts - podcasts resource",
	Long:  "podcasts - podcasts resource",
}
