package cmd

import (
	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

var (
	tags []string
	// note string // Already defined

	// startCmd represents the start command
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a task",
		Long: `Start a task with tags:
	
- traggo_cli start [-t | --tags key1:value1] [-t | --tags key2:value2]
- traggo_cli start -t tag:key -n "Test if this is possible to do"`,
		Run: func(cmd *cobra.Command, args []string) {
			c := config.LoadConfig(configPath)
			s := session.NewTraggoSession(c)
			s.Start(tags, note)
		},
	}
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "List of tags")
	startCmd.Flags().StringVarP(&note, "note", "n", "", "Note associated to this task")
	startCmd.MarkFlagRequired("tags")

}
