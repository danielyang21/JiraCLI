package cmd

import (
	"fmt"
	"os"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Jira API connection",
	Long: `Test your Jira API credentials and connection.

This command will attempt to connect to Jira and retrieve your user information
to verify that your credentials are working correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadAndValidate()

		fmt.Println("Testing Jira API connection...")
		fmt.Printf("URL: %s\n", cfg.JiraURL)
		fmt.Printf("Email: %s\n", cfg.Email)
		fmt.Println("API Token: [HIDDEN]")
		fmt.Println()

		authType := cfg.AuthType
		if authType == "" {
			authType = "basic"
		}
		fmt.Printf("Auth Type: %s\n", authType)

		client := cfg.NewAPIClient()

		fmt.Println("Attempting to connect...")
		err := client.TestConnection()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Connection failed: %v\n\n", err)
			fmt.Fprintln(os.Stderr, "Troubleshooting tips:")
			fmt.Fprintln(os.Stderr, "1. Verify your Jira URL is correct (should start with https://)")
			fmt.Fprintln(os.Stderr, "2. Make sure you're using an API token, not your password")
			fmt.Fprintln(os.Stderr, "3. Check that your email matches your Jira account")
			fmt.Fprintln(os.Stderr, "4. Verify your API token hasn't expired")
			fmt.Fprintln(os.Stderr, "\nTo generate a new API token:")
			fmt.Fprintln(os.Stderr, "  - Jira Cloud: https://id.atlassian.com/manage-profile/security/api-tokens")
			fmt.Fprintln(os.Stderr, "  - Jira Server/DC: Use your username and password")
			fmt.Fprintln(os.Stderr, "\nThen run 'jira init' to update your credentials")
			os.Exit(1)
		}

		fmt.Println("✅ Connection successful!")
		fmt.Println("\nYour Jira API credentials are working correctly.")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
