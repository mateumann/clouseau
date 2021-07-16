package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// collectCmd represents the collect command
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Listen a server collecting traffic flow information",
	Long: `The listener component listens for incoming flow information from your
switch, router, access point or another network interface.  Use it to
receive and store information on traffic flows in your network.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("collect called")
	},
}

func init() {
	rootCmd.AddCommand(collectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// collectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// collectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
