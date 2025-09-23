package cmd

import (
	"fmt"
	"strconv"
	"strings"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	utils "github.com/kalidor/traggo_cli/utils"
	"github.com/spf13/cobra"
)

var (
	rangeIds string
	rmAll    bool
	rmAllYes bool

	// rmCmd represents the rm command
	rmCmd = &cobra.Command{
		Use:   "rm",
		Short: "Delete task(s)",
		Long: `Delete task(s).
	
- ./traggo_cli rm -i 222,223,224
- ./traggo_cli rm 222-300
- ./traggo_cli rm --all # will ask confirmation
- ./traggo_cli rm --all --yes # will NOT ask confirmation
`,
		RunE: runRmE,
	}
)

func runRmE(cmd *cobra.Command, args []string) error {
	c := config.LoadConfig(configPath)
	s := session.NewTraggoSession(c)
	s.ListCurrentTasks()
	if rmAll {
		if rmAllYes {
			fmt.Println("TODO: will remove all without confirmation")
			// s.DeleteAll() // TODO: to implement
		} else {
			r, err := utils.AskAndCompare("Delete all. Confirm (\"Yes, I'm sure\"/N): ", "Yes, I'm sure")
			if err != nil {
				return err
			}
			if r {
				fmt.Println("TODO: will remove all after confirmation")
				// s.DeleteAll()
			} else {
				fmt.Println("Aborting...")
			}
		}
		return nil
	}
	if len(ids) > 0 {
		s.Delete(ids)
		return nil
	}

	if strings.Contains(rangeIds, "-") {
		ids, err := handleRangeIds(rangeIds)
		if err != nil {
			return err
		}
		s.Delete(ids)
		return nil

	}
	return nil
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().IntSliceVarP(&ids, "ids", "i", []int{}, "List of id to delete")
	rmCmd.Flags().StringVarP(&rangeIds, "range", "r", "", "IDs range to delete (1-12)")
	rmCmd.Flags().BoolVar(&rmAll, "all", false, "Remove all tasks. Will ask confirmation.")
	rmCmd.Flags().BoolVar(&rmAllYes, "yes", false, "Remove all tasks without confirmation /!\\")

	rmCmd.MarkFlagsOneRequired("ids", "range")

}

func handleRangeIds(arg string) ([]int, error) {

	splittedArgs := strings.SplitN(arg, "-", 2)
	startId, err := strconv.Atoi(strings.TrimSpace(splittedArgs[0]))
	if err != nil {
		return []int{}, err
	}
	endId, err := strconv.Atoi(strings.TrimSpace(splittedArgs[1]))
	if err != nil {
		return []int{}, err
	}
	// generate all ids between those two
	if startId > endId {
		endId, startId = startId, endId
	}
	ids := make([]int, endId-startId+1)
	index := 0
	for i := startId; i <= endId; i++ {
		ids[index] = i
		index++
	}
	return ids, nil
}
