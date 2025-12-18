package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [ticket-key]",
	Short: "Mark a ticket as done",
	Long: `Mark a ticket as "Done".

This is a quick action that updates the ticket status to "Done".

Examples:
  jira done PROJ-123`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Marking %s as done...\n", ticketKey)

		err := client.UpdateIssueStatus(ticketKey, "Done")
		ui.FatalIfError(err, "Error updating status")

		fmt.Printf("✅ %s is now Done\n", ticketKey)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
