package cmd

import (
	"errors"
	"strconv"
	"strings"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details for specific id",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("this command requiers at least one task id")
		}

		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		var res session.GenericTask
		for _, idStr := range args {
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}
			res = s.SearchTask(id)
			if res == nil {
				continue
			}
			res.PrettyPrint(c.Colors)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
