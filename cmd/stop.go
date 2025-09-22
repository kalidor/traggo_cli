package cmd

import (
	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

var (
	ids []int
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop given IDs",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		s.Stop(c.Colors, ids)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().IntSliceVarP(&ids, "ids", "i", []int{}, "List of id to stop")
	stopCmd.MarkFlagRequired("ids")

}
