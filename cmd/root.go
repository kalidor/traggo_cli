package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	configPath string
	cfgFile    string
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "traggo_cli",
		Short: "Traggo CLI to interact with Traggo using API only.",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	}
)

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
}
