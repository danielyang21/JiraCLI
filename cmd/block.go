package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var blockCmd = &cobra.Command{
	Use:   "block [ticket-key]",
	Short: "Mark a ticket as blocked",
	Long: `Mark a ticket as "Blocked" and optionally add a comment explaining why.

This is a quick action that updates the ticket status to "Blocked".
Use the --reason flag to add a comment explaining the blocker.

Examples:
  jira block PROJ-123
  jira block PROJ-123 --reason "Waiting for API access"
  jira block PROJ-123 -r "Dependencies not ready"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		reason, _ := cmd.Flags().GetString("reason")
		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Marking %s as blocked...\n", ticketKey)

		err := client.UpdateIssueStatus(ticketKey, "Blocked")
		ui.FatalIfError(err, "Error updating status")

		if reason != "" {
			fmt.Printf("Adding comment...\n")
			if err := client.AddComment(ticketKey, reason); err != nil {
				fmt.Printf("Warning: Could not add comment: %v\n", err)
			}
		}

		fmt.Printf("✅ %s is now Blocked\n", ticketKey)
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)
	blockCmd.Flags().StringP("reason", "r", "", "reason for blocking (adds as comment)")
}
