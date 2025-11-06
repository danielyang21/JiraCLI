package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/danielyan21/JiraCLI/internal/api"
	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/fatih/color"
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
		// Load configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
			os.Exit(1)
		}

		// Validate configuration
		if err := config.ValidateConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
			fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
			os.Exit(1)
		}

		// Get flags
		project, _ := cmd.Flags().GetString("project")
		status, _ := cmd.Flags().GetString("status")
		assignee, _ := cmd.Flags().GetString("assignee")
		limit, _ := cmd.Flags().GetInt("limit")

		// Build JQL query
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

		// Create API client with auth type
		authType := cfg.AuthType
		if authType == "" {
			authType = "basic" // default to basic auth for backwards compatibility
		}
		client := api.NewClientWithAuthType(cfg.JiraURL, cfg.Email, cfg.APIToken, authType)

		// Search issues
		fmt.Println("Fetching tickets...")
		results, err := client.SearchIssues(url.QueryEscape(jql), limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching tickets: %v\n", err)
			os.Exit(1)
		}

		// Display results
		if len(results.Issues) == 0 {
			fmt.Println("\nNo tickets found.")
			return
		}

		// Color definitions
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()

		fmt.Printf("\n%s\n\n", bold(fmt.Sprintf("Found %d ticket(s):", results.Total)))
		fmt.Printf("%-12s %-15s %-20s %s\n",
			bold("KEY"),
			bold("STATUS"),
			bold("ASSIGNEE"),
			bold("SUMMARY"))
		fmt.Println(strings.Repeat("-", 80))

		for _, issue := range results.Issues {
			assigneeName := "Unassigned"
			if issue.Fields.Assignee != nil {
				assigneeName = issue.Fields.Assignee.DisplayName
			}

			// Truncate long names and summaries
			if len(assigneeName) > 20 {
				assigneeName = assigneeName[:17] + "..."
			}
			summary := issue.Fields.Summary
			if len(summary) > 40 {
				summary = summary[:37] + "..."
			}

			// Color code status
			statusColor := getStatusColor(issue.Fields.Status.Name)

			fmt.Printf("%-12s %-15s %-20s %s\n",
				cyan(issue.Key),
				statusColor(issue.Fields.Status.Name),
				yellow(assigneeName),
				summary,
			)
		}

		fmt.Printf("\n%s\n", green(fmt.Sprintf("Showing %d of %d total results", len(results.Issues), results.Total)))
	},
}

// getStatusColor returns a color function based on status
func getStatusColor(status string) func(a ...interface{}) string {
	status = strings.ToLower(status)

	// Done/Closed statuses - green
	if strings.Contains(status, "done") || strings.Contains(status, "closed") || strings.Contains(status, "resolved") {
		return color.New(color.FgGreen).SprintFunc()
	}

	// In Progress statuses - yellow
	if strings.Contains(status, "progress") || strings.Contains(status, "review") {
		return color.New(color.FgYellow).SprintFunc()
	}

	// To Do/Open statuses - blue
	if strings.Contains(status, "to do") || strings.Contains(status, "open") || strings.Contains(status, "backlog") {
		return color.New(color.FgBlue).SprintFunc()
	}

	// Default - white
	return color.New(color.FgWhite).SprintFunc()
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Flags for filtering
	listCmd.Flags().StringP("project", "p", "", "filter by project key")
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().StringP("assignee", "a", "", "filter by assignee (@me for yourself)")
	listCmd.Flags().IntP("limit", "l", 20, "maximum number of tickets to show")
}
