package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:   "comment [ticket-key] [comment-text]",
	Short: "Add a comment to a Jira ticket",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := args[0]
		commentText := args[1]

		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		fmt.Printf("Adding comment to %s...\n", ticketKey)
		err := client.AddComment(ticketKey, commentText)
		ui.FatalIfError(err, "Error adding comment")

		fmt.Printf("✅ Comment added successfully to %s\n", ticketKey)
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
