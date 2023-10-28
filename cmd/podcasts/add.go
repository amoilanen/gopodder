package podcasts

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add - add a podcast to the list of downloaded podcasts",
	Long:  ``,
	Run:   runAddCmd,
}

func uniqueOf(values []string) []string {
	uniqueMap := map[string]bool{}
	for _, value := range values {
		uniqueMap[value] = true
	}
	unique := []string{}
	for value := range uniqueMap {
		unique = append(unique, value)
	}
	return unique
}

func runAddCmd(ccmd *cobra.Command, args []string) {
	feedUrl := args[0]
	savedPodcasts := viper.GetStringSlice("feeds")
	updated := append(savedPodcasts, feedUrl)
	viper.Set("feeds", uniqueOf(updated))
	viper.WriteConfig()
}
