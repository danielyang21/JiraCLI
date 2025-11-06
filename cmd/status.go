package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [ticket-key] [new-status]",
	Short: "Update the status of a ticket",
	Long: `Update the status/transition of a Jira ticket.

Supports flexible status names and aliases:
  - Full names: "To Do", "In Progress", "Done"
  - Short codes: td, ip, d, p, r
  - Without spaces: todo, inprogress, done

Examples:
  jira status PROJ-123 done          # Update to Done
  jira status PROJ-123 ip            # Update to In Progress
  jira status PROJ-123 "in progress" # With spaces (needs quotes)
  jira status PROJ-123 td            # Update to To Do`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		newStatus := args[1]

		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Updating %s to '%s'...\n", ticketKey, newStatus)
		err := client.UpdateIssueStatus(ticketKey, newStatus)
		ui.FatalIfError(err, "Error updating status")

		fmt.Printf("Successfully updated %s to '%s'\n", ticketKey, newStatus)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
