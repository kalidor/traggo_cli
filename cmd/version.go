package cmd

import (
	"fmt"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Traggo version",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		fmt.Println(s.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
