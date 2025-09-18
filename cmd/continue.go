package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
)

// continueCmd represents the continue command
var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "continue a previous task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("this command requiers one task id/ticket id")
		}
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		var task session.GenericTask
		re := regexp.MustCompile(`(?P<TagName>[[:word:]]*):(?P<TagValue>[[:word:]]*)`)
		matches := re.FindStringSubmatch(args[0])
		if len(matches) > 0 {
			nIndex := re.SubexpIndex("TagName")
			tagName := matches[nIndex]
			tIndex := re.SubexpIndex("TagValue")
			tagValue := matches[tIndex]

			task = s.SearchTaskByTag(tagName, tagValue)
			// Let's look for this ticket in all tasks in tag 'ticket'
		} else {
			// Let's look for this id in all tasks
			argInt, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			task = s.SearchTask(argInt)
		}
		if task == nil {
			fmt.Println("Unable to retrieve the requested id / tag")
			return nil
		}
		s.Continue(task)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(continueCmd)
}
