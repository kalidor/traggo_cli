package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/kalidor/traggo_cli/utils"
	"github.com/spf13/cobra"
)

var (
	// tags []string // already declared
	// startDateStr  // already declared
	// endDateStr    // already declared
	delNote bool
	note    string
	add     bool

	// updateCmd represents the update command
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update field(s) from specific task ID",
		Long: `Update one or more fields from specific task ID:

- traggo_cli update taskId [-n | --note "This is a note"]
- traggo_cli update taskId [-a -n| --append -n "This is a note append to current note"]
- traggo_cli update taskId [-e | --end-date YYYY/MM/DD]
- traggo_cli update taskId [-s | --start-date YYYY/MM/DD]
- traggo_cli update taskId [-t | --tags ""]
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("this command requiers at least one task id")
			}
			c := config.LoadConfig(configPath)
			s := session.NewTraggoSession(c)
			if note != "" && delNote {
				return errors.New("cannot have --note and --delete-note in same command")
			}
			if len(tags) == 0 && startDateStr == "" && endDateStr == "" && !delNote && note == "" {
				rootCmd.Help()
				return nil
			}

			idStr := strings.TrimSpace(args[0])
			taskId, err := strconv.Atoi(idStr)
			if err != nil {
				return err
			}

			// TODO: avoid code duplication...
			// Update current task
			currentTimerTask, okTimer := s.SearchTask(taskId).(session.TimerTask)
			if okTimer {
				if startDateStr != "" {
					currentTimerTask.OldStart = currentTimerTask.Start
					currentTimerTask.Start, err = utils.StrToTime(startDateStr, time.DateTime)
					if err != nil {
						return err
					}
					fmt.Println(currentTimerTask.Start)
				}
				if note != "" {
					if add {
						currentTimerTask.Note = fmt.Sprintf("%s. %s", currentTimerTask.Note, note)
					} else {
						currentTimerTask.Note = note
					}
				}
				if delNote {
					currentTimerTask.Note = ""
				}
				if len(tags) > 0 {
					if add {
						for _, tag := range tags {
							if !strings.Contains(tag, ":") {
								continue
							}
							s := strings.SplitN(tag, ":", 2)
							currentTimerTask.Tags = append(currentTimerTask.Tags, session.Tag{
								Key:   s[0],
								Value: s[1],
							})

						}
					} else {
						var tagsStruct []session.Tag
						for _, tag := range tags {
							if !strings.Contains(tag, ":") {
								continue
							}
							s := strings.SplitN(tag, ":", 2)
							tagsStruct = append(tagsStruct, session.Tag{
								Key:   s[0],
								Value: s[1],
							})

						}
						currentTimerTask.Tags = tagsStruct
					}
				}
				fmt.Println(currentTimerTask.PreparePretty(c.Colors))
				s.UpdateTimerTask(currentTimerTask)
				return nil
			}

			// Update already done task
			currentTask, okSpan := s.SearchTask(taskId).(session.TimeSpanTask)
			if !okSpan {
				panic("cannot convert task to session.TimeSpanTask")
			}
			if startDateStr != "" {
				currentTask.OldStart = currentTask.Start
				currentTask.Start, err = utils.StrToTime(startDateStr, time.DateTime)
				if err != nil {
					return err
				}
				fmt.Println(currentTask.Start)
			}
			if endDateStr != "" {
				currentTask.End, err = utils.StrToTime(endDateStr, time.DateTime)
				if err != nil {
					return err
				}
			}
			if note != "" {
				if add {
					currentTask.Note = fmt.Sprintf("%s. %s", currentTask.Note, note)
				} else {
					currentTask.Note = note
				}
			}
			if delNote {
				currentTask.Note = ""
			}
			if len(tags) > 0 {
				if add {
					for _, tag := range tags {
						if !strings.Contains(tag, ":") {
							continue
						}
						s := strings.SplitN(tag, ":", 2)
						currentTask.Tags = append(currentTask.Tags, session.Tag{
							Key:   s[0],
							Value: s[1],
						})

					}
				} else {
					var tagsStruct []session.Tag
					for _, tag := range tags {
						if !strings.Contains(tag, ":") {
							continue
						}
						s := strings.SplitN(tag, ":", 2)
						tagsStruct = append(tagsStruct, session.Tag{
							Key:   s[0],
							Value: s[1],
						})

					}
					currentTask.Tags = tagsStruct
				}
			}
			fmt.Println(currentTask.PreparePretty(c.Colors))
			s.UpdateTimeSpanTask(currentTask)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringArrayVarP(&tags, "tags", "t", []string{}, "List of tags")
	updateCmd.Flags().BoolVarP(&add, "add", "a", false, "Append note to current one. Only useful when using '-n|--note' or '-t|--tags' flags.")
	updateCmd.Flags().BoolVarP(&delNote, "delete-note", "d", false, "Delete note")
	updateCmd.Flags().StringVarP(&note, "note", "n", "", "Note to add to task ID")
	updateCmd.Flags().StringVarP(&startDateStr, "start-date", "s", "", "Task new start date")
	updateCmd.Flags().StringVarP(&endDateStr, "end-date", "e", "", "Task new end date")

}
