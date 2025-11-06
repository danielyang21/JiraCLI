package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var assignCmd = &cobra.Command{
	Use:   "assign [ticket-key] [new-assignee]",
	Short: "Update the assignee of a ticket",
	Long: `Update the assignee of a Jira ticket.

Examples:
  jira assign PROJ-123 @me          # Assign ticket to self`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		newAssignee := args[1]

		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Assigning %s to '%s'...\n", ticketKey, newAssignee)
		err := client.AssignIssue(ticketKey, newAssignee)
		ui.FatalIfError(err, "Error updating assignee")

		fmt.Printf("Successfully assigned %s to '%s'\n", ticketKey, newAssignee)
	},
}

func init() {
	rootCmd.AddCommand(assignCmd)
}
