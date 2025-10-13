package cmd

import (
	"fmt"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

var (
	// testCmd represents the test command
	testCmd = &cobra.Command{
		Use:   "check",
		Short: "Check API connectivity with current token",
		Run: func(cmd *cobra.Command, args []string) {
			c := config.LoadConfig(configPath)
			s := session.NewTraggoSession(c)
			err := s.Ping()
			if err != nil {
				fmt.Println("Unable to request the API", err)
				return
			}
			fmt.Println("Ping success!")
		},
	}
)

func init() {
	rootCmd.AddCommand(testCmd)
}
