package podcasts

import "github.com/spf13/cobra"

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch - fetches the podcast feed provided as an argument",
	Long:  ``,
	Run:   runFetchCmd,
}

func runFetchCmd(ccmd *cobra.Command, args []string) {
	//TODO: Implement
	println("Fetching feed...")
	feedUrl := args[0]
	println(feedUrl)
}
