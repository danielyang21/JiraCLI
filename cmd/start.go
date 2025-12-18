package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [ticket-key]",
	Short: "Start working on a ticket",
	Long: `Mark a ticket as "In Progress" and assign it to yourself.

This is a quick action that combines:
  - Updating status to "In Progress"
  - Assigning the ticket to you

Examples:
  jira start PROJ-123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Starting work on %s...\n", ticketKey)

		if err := client.AssignIssue(ticketKey, "@me"); err != nil {
			fmt.Printf("Warning: Could not assign ticket: %v\n", err)
		}

		err := client.UpdateIssueStatus(ticketKey, "In Progress")
		ui.FatalIfError(err, "Error updating status")

		fmt.Printf("✅ %s is now In Progress and assigned to you\n", ticketKey)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
