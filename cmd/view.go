package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view [ticket-key]",
	Short: "View detailed information about a ticket",
	Long: `View comprehensive details about a specific Jira ticket.

Examples:
  jira view PROJ-123        # View full ticket details
  jira view PROJ-123 -c     # View ticket with comments`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		showComments, _ := cmd.Flags().GetBool("comments")

		fmt.Printf("Fetching details for %s...\n\n", ticketKey)

		// TODO: Implement actual API call
		fmt.Println("[TODO: API integration coming soon]")
		fmt.Println("\nExample output:")
		fmt.Printf("Key:         %s\n", ticketKey)
		fmt.Println("Type:        Bug")
		fmt.Println("Status:      In Progress")
		fmt.Println("Priority:    High")
		fmt.Println("Assignee:    John Doe")
		fmt.Println("Reporter:    Jane Smith")
		fmt.Println("Created:     2024-01-15")
		fmt.Println("Updated:     2024-01-16")
		fmt.Println("\nSummary:")
		fmt.Println("  Fix critical login bug")
		fmt.Println("\nDescription:")
		fmt.Println("  Users are unable to log in when using OAuth.")

		if showComments {
			fmt.Println("\nComments:")
			fmt.Println("  [2024-01-16] John Doe: Working on a fix")
			fmt.Println("  [2024-01-15] Jane Smith: This is blocking production")
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Flags
	viewCmd.Flags().BoolP("comments", "c", false, "show comments")
	viewCmd.Flags().BoolP("full", "f", false, "show full details including custom fields")
}
