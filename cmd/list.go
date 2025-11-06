package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

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
		cfg := config.LoadAndValidate()
		jql := buildJQLQuery(cmd, cfg)
		limit, _ := cmd.Flags().GetInt("limit")
		client := cfg.NewAPIClient()

		fmt.Println("Fetching tickets...")
		results, err := client.SearchIssues(url.QueryEscape(jql), limit)
		ui.FatalIfError(err, "Error fetching tickets")

		ui.RenderIssueList(results)
	},
}

func buildJQLQuery(cmd *cobra.Command, cfg *config.Config) string {
	project, _ := cmd.Flags().GetString("project")
	status, _ := cmd.Flags().GetString("status")
	assignee, _ := cmd.Flags().GetString("assignee")

	var jqlParts []string

	// Use default project if specified and no project flag provided
	if project != "" {
		jqlParts = append(jqlParts, fmt.Sprintf("project = %s", project))
	} else if cfg.DefaultProject != "" {
		jqlParts = append(jqlParts, fmt.Sprintf("project = %s", cfg.DefaultProject))
	}

	if status != "" {
		jqlParts = append(jqlParts, fmt.Sprintf("status = \"%s\"", status))
	}

	if assignee != "" {
		if assignee == "@me" {
			jqlParts = append(jqlParts, "assignee = currentUser()")
		} else {
			jqlParts = append(jqlParts, fmt.Sprintf("assignee = \"%s\"", assignee))
		}
	}

	// Default query if no filters provided
	jql := "ORDER BY updated DESC"
	if len(jqlParts) > 0 {
		jql = strings.Join(jqlParts, " AND ") + " " + jql
	}

	return jql
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("project", "p", "", "filter by project key")
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().StringP("assignee", "a", "", "filter by assignee (@me for yourself)")
	listCmd.Flags().IntP("limit", "l", 20, "maximum number of tickets to show")
}
