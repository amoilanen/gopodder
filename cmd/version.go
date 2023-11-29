package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gopodder",
	Long:  "Print the version number of gopodder",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gopodder podcast fetching client v0.0.1")
	},
}
