package cmd

import (
	"fmt"
	"os"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test [ticket-key]",
	Short: "Test Jira API connection or debug ticket data",
	Long: `Test your Jira API credentials and connection.
If a ticket key is provided, shows raw API response for debugging.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadAndValidate()

		if len(args) > 0 {
			debugTicket(cfg, args[0])
			return
		}

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

func debugTicket(cfg *config.Config, ticketKey string) {
	client := cfg.NewAPIClient()
	c := ui.NewColorFuncs()

	fmt.Printf("Fetching %s for debugging...\n\n", c.Cyan(ticketKey))

	apiVersion := "3"
	if cfg.AuthType == "pat" {
		apiVersion = "2"
	}
	fmt.Printf("%s %s\n", c.Bold("Auth type:"), cfg.AuthType)
	fmt.Printf("%s %s\n", c.Bold("API version:"), apiVersion)

	// Try without any field filter first
	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s", apiVersion, ticketKey)
	fmt.Printf("%s %s\n\n", c.Bold("Endpoint:"), endpoint)

	issue, err := client.GetIssue(ticketKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", c.Red("Error:"), err)
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", c.Bold("Issue key:"), c.Cyan(issue.Key))
	fmt.Printf("%s %s\n", c.Bold("Summary:"), issue.Fields.Summary)
	fmt.Printf("\n%s %T\n", c.Bold("Description field type:"), issue.Fields.Description)
	if issue.Fields.Description == nil {
		fmt.Printf("%s\n", c.Red("Description is nil (not returned by API)"))
		fmt.Printf("\n%s\n", c.Yellow("This might mean:"))
		fmt.Println("  1. The ticket has no description")
		fmt.Println("  2. The API version needs different field parameters")
		fmt.Println("  3. The field name might be different in your Jira instance")
	} else {
		switch desc := issue.Fields.Description.(type) {
		case string:
			fmt.Printf("%s %q\n", c.Bold("Description (string):"), desc)
		case map[string]interface{}:
			fmt.Printf("%s\n", c.Bold("Description (ADF format):"))
			fmt.Printf("  %s %v\n", c.Bold("Type:"), desc["type"])
			if content, ok := desc["content"].([]interface{}); ok {
				fmt.Printf("  %s %d\n", c.Bold("Content items:"), len(content))
				for i, item := range content {
					if itemMap, ok := item.(map[string]interface{}); ok {
						fmt.Printf("    %s %v\n", c.Cyan(fmt.Sprintf("[%d] Type:", i)), itemMap["type"])
					}
				}
			}
		default:
			fmt.Printf("%s %#v\n", c.Bold("Description (unknown type):"), desc)
		}
	}
}

func init() {
	rootCmd.AddCommand(testCmd)
}
