package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/danielyan21/JiraCLI/internal/api"
	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
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

		// Create API client with auth type
		authType := cfg.AuthType
		if authType == "" {
			authType = "basic" // default to basic auth for backwards compatibility
		}
		client := api.NewClientWithAuthType(cfg.JiraURL, cfg.Email, cfg.APIToken, authType)

		// Fetch issue
		fmt.Printf("Fetching details for %s...\n\n", ticketKey)
		issue, err := client.GetIssue(ticketKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching ticket: %v\n", err)
			os.Exit(1)
		}

		// Color definitions
		cyan := color.New(color.FgCyan).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		bold := color.New(color.Bold).SprintFunc()
		gray := color.New(color.FgHiBlack).SprintFunc()

		// Display issue details
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("%s %s\n", bold("Key:"), cyan(issue.Key))
		fmt.Printf("%s %s\n", bold("Type:"), issue.Fields.IssueType.Name)

		// Color-coded status
		statusColor := getStatusColorForView(issue.Fields.Status.Name)
		fmt.Printf("%s %s\n", bold("Status:"), statusColor(issue.Fields.Status.Name))

		// Color-coded priority
		priorityColor := getPriorityColor(issue.Fields.Priority.Name)
		fmt.Printf("%s %s\n", bold("Priority:"), priorityColor(issue.Fields.Priority.Name))

		if issue.Fields.Assignee != nil {
			fmt.Printf("%s %s\n", bold("Assignee:"), yellow(issue.Fields.Assignee.DisplayName))
		} else {
			fmt.Printf("%s %s\n", bold("Assignee:"), gray("Unassigned"))
		}

		if issue.Fields.Reporter != nil {
			fmt.Printf("%s %s\n", bold("Reporter:"), issue.Fields.Reporter.DisplayName)
		}

		fmt.Printf("%s %s (%s)\n", bold("Project:"), issue.Fields.Project.Name, issue.Fields.Project.Key)
		fmt.Printf("%s %s\n", bold("Created:"), gray(issue.Fields.Created.Format("2006-01-02 15:04")))
		fmt.Printf("%s %s\n", bold("Updated:"), gray(issue.Fields.Updated.Format("2006-01-02 15:04")))
		fmt.Println(strings.Repeat("=", 80))

		fmt.Printf("\n%s\n", bold("Summary:"))
		fmt.Printf("  %s\n", issue.Fields.Summary)

		// Handle description (can be string or complex object)
		fmt.Printf("\n%s\n", bold("Description:"))
		if issue.Fields.Description != nil {
			switch desc := issue.Fields.Description.(type) {
			case string:
				if desc != "" {
					fmt.Printf("  %s\n", desc)
				} else {
					fmt.Printf("  %s\n", gray("(No description)"))
				}
			case map[string]interface{}:
				// Atlassian Document Format (ADF) - simplified display
				fmt.Printf("  %s\n", gray("(Rich text description available in web UI)"))
			default:
				fmt.Printf("  %s\n", gray("(No description)"))
			}
		} else {
			fmt.Printf("  %s\n", gray("(No description)"))
		}

		// Show link to web UI
		fmt.Printf("\n%s %s\n", green("View in browser:"), cyan(fmt.Sprintf("%s/browse/%s", cfg.JiraURL, issue.Key)))

		// TODO: Implement comments when showComments is true
		if showComments {
			fmt.Println("\n[Comments feature coming soon]")
		}
	},
}

// getStatusColorForView returns a color function based on status
func getStatusColorForView(status string) func(a ...interface{}) string {
	statusLower := strings.ToLower(status)

	// Done/Closed statuses - green
	if strings.Contains(statusLower, "done") || strings.Contains(statusLower, "closed") || strings.Contains(statusLower, "resolved") {
		return color.New(color.FgGreen, color.Bold).SprintFunc()
	}

	// In Progress statuses - yellow
	if strings.Contains(statusLower, "progress") || strings.Contains(statusLower, "review") {
		return color.New(color.FgYellow, color.Bold).SprintFunc()
	}

	// To Do/Open statuses - blue
	if strings.Contains(statusLower, "to do") || strings.Contains(statusLower, "open") || strings.Contains(statusLower, "backlog") {
		return color.New(color.FgBlue, color.Bold).SprintFunc()
	}

	// Default - white
	return color.New(color.FgWhite, color.Bold).SprintFunc()
}

// getPriorityColor returns a color function based on priority
func getPriorityColor(priority string) func(a ...interface{}) string {
	priorityLower := strings.ToLower(priority)

	// High/Critical priorities - red
	if strings.Contains(priorityLower, "highest") || strings.Contains(priorityLower, "critical") {
		return color.New(color.FgRed, color.Bold).SprintFunc()
	}

	// High priority - red
	if strings.Contains(priorityLower, "high") {
		return color.New(color.FgRed).SprintFunc()
	}

	// Medium priority - yellow
	if strings.Contains(priorityLower, "medium") {
		return color.New(color.FgYellow).SprintFunc()
	}

	// Low priority - green
	if strings.Contains(priorityLower, "low") || strings.Contains(priorityLower, "lowest") {
		return color.New(color.FgGreen).SprintFunc()
	}

	// Default - white
	return color.New(color.FgWhite).SprintFunc()
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Flags
	viewCmd.Flags().BoolP("comments", "c", false, "show comments")
	viewCmd.Flags().BoolP("full", "f", false, "show full details including custom fields")
}
