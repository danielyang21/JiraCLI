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
  jira list                      # Your tickets (default)
  jira list --mine               # Your tickets (explicit)
  jira list --all                # All tickets in default project
  jira list --recent             # Recently updated (last 7 days)
  jira list -p KAN               # All tickets in KAN project
  jira list -s "In Progress"     # Tickets with specific status
  jira list -a @me -s Done       # Your done tickets`,
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
	mine, _ := cmd.Flags().GetBool("mine")
	all, _ := cmd.Flags().GetBool("all")
	recent, _ := cmd.Flags().GetBool("recent")
	project, _ := cmd.Flags().GetString("project")
	status, _ := cmd.Flags().GetString("status")
	assignee, _ := cmd.Flags().GetString("assignee")

	var jqlParts []string

	if recent {
		jqlParts = append(jqlParts, "updated >= -7d")
	}

	if all {
		if cfg.DefaultProject != "" && project == "" {
			jqlParts = append(jqlParts, fmt.Sprintf("project = %s", cfg.DefaultProject))
		}
	} else if mine || (!all && assignee == "" && !recent) {
		jqlParts = append(jqlParts, "assignee = currentUser()")
	}

	if project != "" {
		jqlParts = append(jqlParts, fmt.Sprintf("project = %s", project))
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

	if len(jqlParts) == 0 {
		jqlParts = append(jqlParts, "assignee = currentUser()")
	}

	return strings.Join(jqlParts, " AND ") + " ORDER BY status ASC, updated DESC"
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().Bool("mine", false, "show only your tickets (default behavior)")
	listCmd.Flags().Bool("all", false, "show all tickets in default project")
	listCmd.Flags().Bool("recent", false, "show recently updated tickets (last 7 days)")
	listCmd.Flags().StringP("project", "p", "", "filter by project key")
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().StringP("assignee", "a", "", "filter by assignee (@me for yourself)")
	listCmd.Flags().IntP("limit", "l", 20, "maximum number of tickets to show")
}
