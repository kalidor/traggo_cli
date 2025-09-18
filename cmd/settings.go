package cmd

import (
	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Retrieve userSettings",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		s.GetSettings()
	},
}

func init() {
	rootCmd.AddCommand(settingsCmd)
}
