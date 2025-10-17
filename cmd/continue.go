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
			return errors.New("this command requiers one task_id or TagName:TagValue")
		}
		c := config.LoadConfig(configPath)
		s := session.NewTraggoSession(c)
		var task session.GenericTask
		re := regexp.MustCompile(`(?P<TagName>[[:word:]]*):(?P<TagValue>[a-zA-Z_\-0-9]+)`)
		matches := re.FindStringSubmatch(args[0])
		if len(matches) > 0 {
			// Let's look by 'TagName & TagValue'
			nIndex := re.SubexpIndex("TagName")
			tagName := matches[nIndex]
			tIndex := re.SubexpIndex("TagValue")
			tagValue := matches[tIndex]

			task = s.SearchTaskByTag(tagName, tagValue)
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
		// TODO: show freshly created continued task
		return nil
	},
}

func init() {
	rootCmd.AddCommand(continueCmd)
}
