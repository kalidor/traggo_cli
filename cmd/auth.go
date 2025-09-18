package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	config "github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	DEFAULT_ENDPOINT = "http://localhost:3030/graphql"
	// authCmd represents the auth command
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Request token and save it for later use",
		Long: `Request token and save it in configuration file for later use. Example:
	- ./traggo_cli auth --save-token`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: check if current configuration exists
			// Load it if so in order to keep other configuration items than auth:
			// - colors
			//
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Full URL endpoint (default: %s): ", DEFAULT_ENDPOINT)
			url, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			url = strings.TrimSuffix(url, "\n")
			if url == "" {
				url = DEFAULT_ENDPOINT
			}
			fmt.Printf("Login: ")
			login, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			login = strings.TrimSuffix(login, "\n")
			fmt.Printf("Password (Hidden): ")
			bytePassword, err := term.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}

			token, err := session.RequestPermanentTokenAndTest(url, login, string(bytePassword))
			if err != nil {
				return err
			}
			fmt.Println("Ping success!")
			config.NewConfig(url, token).Save(configPath)
			fmt.Println("You can edit the configuration file to add some color by tagName:tagValue")
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(authCmd)
}
