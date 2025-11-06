package cmd

import (
	"fmt"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

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

		cfg := config.LoadAndValidate()

		client := cfg.NewAPIClient()

		fmt.Printf("Fetching details for %s...\n\n", ticketKey)
		issue, err := client.GetIssue(ticketKey)
		ui.FatalIfError(err, "Error fetching ticket")

		ui.RenderIssueDetail(issue, cfg.JiraURL, showComments)
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolP("comments", "c", false, "show comments")
	viewCmd.Flags().BoolP("full", "f", false, "show full details including custom fields")
}
