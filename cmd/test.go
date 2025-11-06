package cmd

import (
	"fmt"
	"os"

	"github.com/danielyan21/JiraCLI/internal/api"
	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Jira API connection",
	Long: `Test your Jira API credentials and connection.

This command will attempt to connect to Jira and retrieve your user information
to verify that your credentials are working correctly.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
			os.Exit(1)
		}

		// Validate configuration
		if err := config.ValidateConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
			os.Exit(1)
		}

		fmt.Println("Testing Jira API connection...")
		fmt.Printf("URL: %s\n", cfg.JiraURL)
		fmt.Printf("Email: %s\n", cfg.Email)
		fmt.Println("API Token: [HIDDEN]")
		fmt.Println()

		// Create API client with auth type
		authType := cfg.AuthType
		if authType == "" {
			authType = "basic" // default to basic auth for backwards compatibility
		}
		fmt.Printf("Auth Type: %s\n", authType)
		client := api.NewClientWithAuthType(cfg.JiraURL, cfg.Email, cfg.APIToken, authType)

		// Test connection
		fmt.Println("Attempting to connect...")
		err = client.TestConnection()
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
