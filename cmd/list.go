package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	utils "github.com/kalidor/traggo_cli/utils"
	"github.com/spf13/cobra"
)

var (
	endDateStr   string
	highlight    string
	period       string
	startDateStr string
	today        bool
	listCmd      = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List current running tasks",
		Long: `List current running tasks. Examples:
- ./traggo_cli list
- ./traggo_cli list [-s | --start-date 2025-08-12] [-e | --end-date 2025-08-20]
- './traggo_cli list --period -1m' is the same as './traggo_cli list -s 2025-07-01 -e 2025-08-22' # if today is 2025-08-22
- ./traggo_cli list --period 1w`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := config.LoadConfig(configPath)
			s := session.NewTraggoSession(c)

			var (
				endDate   time.Time
				startDate time.Time
				err       error
			)

			delta := func(sDate time.Time, eDate *time.Time) {}

			if period != "" {
				re := regexp.MustCompile(`(?P<Number>(?:-)?\d+)(?P<Type>[[:alpha:]]{1})`)
				matches := re.FindStringSubmatch(period)
				if len(matches) > 0 {
					nIndex := re.SubexpIndex("Number")
					nString := matches[nIndex]
					number, _ := strconv.Atoi(nString)
					tIndex := re.SubexpIndex("Type")

					c := matches[tIndex]
					switch c {
					case "d":
						delta = func(sDate time.Time, eDate *time.Time) {
							*eDate = sDate.AddDate(0, 0, number)
						}

					case "m":
						delta = func(sDate time.Time, eDate *time.Time) {
							*eDate = sDate.AddDate(0, number, 0)
						}

					case "w":
						delta = func(sDate time.Time, eDate *time.Time) {
							*eDate = sDate.AddDate(0, 0, number*7)
						}
					default:
						return fmt.Errorf("invalid period provided: '%s'", period)
					}
				}
			}

			if startDateStr != "" {
				startDate, err = utils.StrToTime(startDateStr, time.DateOnly)
				if err != nil {
					return err
				}

				if period != "" {
					delta(startDate, &endDate)
				}
				fmt.Println(startDate)
			}

			if endDateStr != "" {
				endDate, err = utils.StrToTime(endDateStr, time.DateOnly)
				if err != nil {
					return err
				}

				if period != "" {
					delta(startDate, &endDate)
				}
			}

			if startDate.IsZero() && endDate.IsZero() && !today {
				if period == "" {
					// if there is no parameter, display current tasks
					tasks := s.ListCurrentTasks()
					if tasks.IsEmpty() {
						return nil
					}
					tasks.PrettyPrint(c.Colors, highlight)
				} else {
					endDate = time.Now()
					delta(endDate, &startDate)
				}
			}
			if today {
				tmp := time.Now().Format(time.DateOnly)
				startDate, _ = utils.StrToTime(tmp, time.DateOnly)
				// Done tasks
				fmt.Printf("Date range: [%s -> %s]\n", startDate.Format(time.DateOnly), startDate.Format(time.DateOnly))
				doneTasks := s.ListBetweenDates(startDate, time.Now())
				if doneTasks.IsEmpty() {
					return nil
				}
				doneTasks.PrettyPrint(c.Colors, highlight)

				tasks := s.ListCurrentTasks()
				if tasks.IsEmpty() {
					return nil
				}
				tasks.PrettyPrint(c.Colors, highlight)
				return nil
			}

			if !startDate.IsZero() && !today && period == "" && endDate.IsZero() {
				fmt.Println("coin")
				endDate = time.Now()
			}

			if !startDate.IsZero() && !endDate.IsZero() {
				fmt.Printf("Date range: [%s -> %s]\n", startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))
				tasks := s.ListBetweenDates(startDate, endDate)
				if tasks.IsEmpty() {
					return nil
				}
				tasks.PrettyPrint(c.Colors, highlight)
			}

			return nil

		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(
		&today,
		"today",
		"t",
		false, // default value
		"List today tasks",
	)
	listCmd.Flags().StringVarP(
		&startDateStr,
		"start-date",
		"s",
		"", // default value
		"Start date to list tasks. To use with -end-date/-e",
	)
	listCmd.Flags().StringVarP(
		&endDateStr,
		"end-date",
		"e",
		"",
		"End date to list tasks. To use with -start-date/-s",
	)
	listCmd.Flags().StringVarP(
		&period,
		"period",
		"p",
		"",
		"Period of time to list task (1d= 1day, 1m=1 month). To be used with --month or --start-date´",
	)
	listCmd.Flags().StringVarP(
		&highlight,
		"Highlight",
		"H",
		"",
		"Highlight line with matching string (case sensitive) in Tags or Note")

}
