package ui

import (
	"fmt"
	"strings"

	"github.com/danielyan21/JiraCLI/internal/api"
)

func RenderIssueList(results *api.SearchResults) {
	if len(results.Issues) == 0 {
		fmt.Println("\nNo tickets found.")
		return
	}

	c := NewColorFuncs()

	actualCount := len(results.Issues)
	totalCount := results.Total
	if totalCount == 0 {
		totalCount = actualCount // Fallback if API doesn't return total
	}

	fmt.Printf("\n%s\n\n", c.Bold(fmt.Sprintf("Found %d ticket(s):", actualCount)))
	printTableHeader(c)

	for _, issue := range results.Issues {
		printIssueRow(issue, c)
	}

	fmt.Printf("\n%s\n", c.Green(fmt.Sprintf("Showing %d of %d total results", actualCount, totalCount)))
}

func RenderIssueDetail(issue *api.Issue, jiraURL string, showComments bool) {
	c := NewColorFuncs()

	printIssueHeader(issue, c)
	printIssueSummary(issue, c)
	printIssueDescription(issue, c)
	printBrowserLink(issue, jiraURL, c)

	if showComments {
		fmt.Println("\n[Comments feature coming soon]")
	}
}

func printTableHeader(c *ColorFuncs) {
	fmt.Printf("%-12s %-15s %-20s %s\n",
		c.Bold("KEY"),
		c.Bold("STATUS"),
		c.Bold("ASSIGNEE"),
		c.Bold("SUMMARY"))
	fmt.Println(strings.Repeat("-", 80))
}

func printIssueRow(issue api.Issue, c *ColorFuncs) {
	assigneeName := "Unassigned"
	if issue.Fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.DisplayName
	}

	assigneeName = Truncate(assigneeName, 20)
	summary := Truncate(issue.Fields.Summary, 40)
	statusColor := GetStatusColor(issue.Fields.Status.Name)

	fmt.Printf("%-12s %-15s %-20s %s\n",
		c.Cyan(issue.Key),
		statusColor(issue.Fields.Status.Name),
		c.Yellow(assigneeName),
		summary,
	)
}

func printIssueHeader(issue *api.Issue, c *ColorFuncs) {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%s %s\n", c.Bold("Key:"), c.Cyan(issue.Key))
	fmt.Printf("%s %s\n", c.Bold("Type:"), issue.Fields.IssueType.Name)

	statusColor := GetStatusColor(issue.Fields.Status.Name)
	fmt.Printf("%s %s\n", c.Bold("Status:"), statusColor(issue.Fields.Status.Name))

	priorityColor := GetPriorityColor(issue.Fields.Priority.Name)
	fmt.Printf("%s %s\n", c.Bold("Priority:"), priorityColor(issue.Fields.Priority.Name))

	if issue.Fields.Assignee != nil {
		fmt.Printf("%s %s\n", c.Bold("Assignee:"), c.Yellow(issue.Fields.Assignee.DisplayName))
	} else {
		fmt.Printf("%s %s\n", c.Bold("Assignee:"), c.Gray("Unassigned"))
	}

	if issue.Fields.Reporter != nil {
		fmt.Printf("%s %s\n", c.Bold("Reporter:"), issue.Fields.Reporter.DisplayName)
	}

	fmt.Printf("%s %s (%s)\n", c.Bold("Project:"), issue.Fields.Project.Name, issue.Fields.Project.Key)
	fmt.Printf("%s %s\n", c.Bold("Created:"), c.Gray(issue.Fields.Created.Format("2006-01-02 15:04")))
	fmt.Printf("%s %s\n", c.Bold("Updated:"), c.Gray(issue.Fields.Updated.Format("2006-01-02 15:04")))
	fmt.Println(strings.Repeat("=", 80))
}

func printIssueSummary(issue *api.Issue, c *ColorFuncs) {
	fmt.Printf("\n%s\n", c.Bold("Summary:"))
	fmt.Printf("  %s\n", issue.Fields.Summary)
}

func printIssueDescription(issue *api.Issue, c *ColorFuncs) {
	fmt.Printf("\n%s\n", c.Bold("Description:"))
	if issue.Fields.Description != nil {
		switch desc := issue.Fields.Description.(type) {
		case string:
			if desc != "" {
				fmt.Printf("  %s\n", desc)
			} else {
				fmt.Printf("  %s\n", c.Gray("(No description)"))
			}
		case map[string]interface{}:
			fmt.Printf("  %s\n", c.Gray("(Rich text description available in web UI)"))
		default:
			fmt.Printf("  %s\n", c.Gray("(No description)"))
		}
	} else {
		fmt.Printf("  %s\n", c.Gray("(No description)"))
	}
}

func printBrowserLink(issue *api.Issue, jiraURL string, c *ColorFuncs) {
	fmt.Printf("\n%s %s\n", c.Green("View in browser:"), c.Cyan(fmt.Sprintf("%s/browse/%s", jiraURL, issue.Key)))
}
