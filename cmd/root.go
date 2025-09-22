package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	configPath string
	verbose    bool
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:              "traggo_cli",
		Short:            "Traggo CLI to interact with Traggo using API only.",
		PersistentPreRun: preRunRoot,
	}
	printBody func(r http.Response)
)

func preRunRoot(cmd *cobra.Command, args []string) {
	if verbose {
		fmt.Printf("verbose=%t\n", verbose)
		printBody = func(r http.Response) {
			b, _ := io.ReadAll(r.Body)
			fmt.Println(string(b))
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	defaultConfigPath := filepath.Join(homeDir, ".config/traggo_cli/config.json")

	rootCmd.Flags().BoolP("help", "h", false, "Help message")
	rootCmd.Flags().StringVarP(&configPath, "config", "c", defaultConfigPath, "Full path of config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print body response")
}
