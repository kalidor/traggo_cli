package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/kalidor/traggo_cli/tui"
	"github.com/spf13/cobra"
)

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Live dashboard useful to interact with traggo",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		err := s.CheckTagsInConfig()
		// TODO: add command to force tag creation
		if err != nil {
			panic(err)
		}

		var dump *os.File
		if _, ok := os.LookupEnv("DEBUG"); ok {
			var err error
			dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
			if err != nil {
				os.Exit(1)
			}
		}
		m, _ := tui.NewMainModel(dump, s, tui.TableView)

		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(liveCmd)

}
