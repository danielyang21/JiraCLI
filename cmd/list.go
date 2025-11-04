package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Jira tickets",
	Long: `List Jira tickets with various filtering options.

Examples:
  jira list                    # List recent tickets
  jira list --project PROJ     # List tickets in a specific project
  jira list --status "To Do"   # List tickets with specific status
  jira list --assignee @me     # List tickets assigned to you`,
	Run: func(cmd *cobra.Command, args []string) {
		project, _ := cmd.Flags().GetString("project")
		status, _ := cmd.Flags().GetString("status")
		assignee, _ := cmd.Flags().GetString("assignee")

		fmt.Println("Listing Jira tickets...")

		// TODO: Implement actual API call
		fmt.Println("\nFilters:")
		if project != "" {
			fmt.Printf("  Project: %s\n", project)
		}
		if status != "" {
			fmt.Printf("  Status: %s\n", status)
		}
		if assignee != "" {
			fmt.Printf("  Assignee: %s\n", assignee)
		}

		fmt.Println("\n[TODO: API integration coming soon]")
		fmt.Println("\nExample output:")
		fmt.Println("PROJ-123  In Progress  Fix login bug")
		fmt.Println("PROJ-124  To Do        Add dark mode")
		fmt.Println("PROJ-125  Done         Update README")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Flags for filtering
	listCmd.Flags().StringP("project", "p", "", "filter by project key")
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().StringP("assignee", "a", "", "filter by assignee (@me for yourself)")
	listCmd.Flags().IntP("limit", "l", 20, "maximum number of tickets to show")
}
