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

func RenderIssueDetail(issue *api.Issue, jiraURL string, comments []api.Comment) {
	c := NewColorFuncs()

	printIssueHeader(issue, c)
	printIssueSummary(issue, c)
	printIssueDescription(issue, c)
	printBrowserLink(issue, jiraURL, c)

	if len(comments) > 0 {
		printComments(comments, c)
	}
}

func printTableHeader(c *ColorFuncs) {
	keyHeader := fmt.Sprintf("%-10s", "KEY")
	statusHeader := fmt.Sprintf("%-12s", "STATUS")
	assigneeHeader := fmt.Sprintf("%-12s", "ASSIGNEE")
	fmt.Printf("%s %s %s %s\n",
		keyHeader,
		statusHeader,
		assigneeHeader,
		"SUMMARY",
	)
	fmt.Println(strings.Repeat("-", 70))
}

func printIssueRow(issue api.Issue, c *ColorFuncs) {
	assigneeName := "Unassigned"
	if issue.Fields.Assignee != nil {
		assigneeName = issue.Fields.Assignee.DisplayName
	}

	assigneeName = Truncate(assigneeName, 20)
	summary := Truncate(issue.Fields.Summary, 40)
	statusColor := GetStatusColor(issue.Fields.Status.Name)

	key := fmt.Sprintf("%-10s", issue.Key)
	status := fmt.Sprintf("%-12s", issue.Fields.Status.Name)
	assignee := fmt.Sprintf("%-12s", assigneeName)

	fmt.Printf("%s %s %s %s\n",
		c.Cyan(key),
		statusColor(status),
		c.Yellow(assignee),
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
				wrappedText := wrapText(desc, 78)
				for _, line := range strings.Split(wrappedText, "\n") {
					fmt.Printf("  %s\n", line)
				}
			} else {
				fmt.Printf("  %s\n", c.Gray("(No description)"))
			}
		case map[string]interface{}:
			// Try to extract text from ADF format
			text := extractDescriptionFromADF(desc)
			if text != "" {
				wrappedText := wrapText(text, 78)
				for _, line := range strings.Split(wrappedText, "\n") {
					fmt.Printf("  %s\n", line)
				}
			} else {
				fmt.Printf("  %s\n", c.Gray("(No description)"))
			}
		default:
			fmt.Printf("  %s\n", c.Gray("(No description)"))
		}
	} else {
		fmt.Printf("  %s\n", c.Gray("(No description)"))
	}
}

func extractDescriptionFromADF(adf map[string]interface{}) string {
	var text string

	if adf["type"] == "text" {
		if textVal, ok := adf["text"].(string); ok {
			return textVal
		}
	}

	if content, ok := adf["content"].([]interface{}); ok {
		for i, item := range content {
			if itemMap, ok := item.(map[string]interface{}); ok {
				itemText := extractDescriptionFromADF(itemMap)
				if itemText != "" {
					// Add spacing between paragraphs
					if i > 0 && adf["type"] == "doc" {
						text += "\n\n"
					}
					text += itemText
				}
			}
		}
	}

	return text
}

func printBrowserLink(issue *api.Issue, jiraURL string, c *ColorFuncs) {
	fmt.Printf("\n%s %s\n", c.Green("View in browser:"), c.Cyan(fmt.Sprintf("%s/browse/%s", jiraURL, issue.Key)))
}

func printComments(comments []api.Comment, c *ColorFuncs) {
	if len(comments) == 0 {
		fmt.Printf("\n%s\n", c.Gray("No comments"))
		return
	}

	fmt.Printf("\n%s (%d)\n", c.Bold("Comments:"), len(comments))
	fmt.Println(strings.Repeat("-", 80))

	for i, comment := range comments {
		if i > 0 {
			fmt.Println(strings.Repeat("-", 80))
		}

		fmt.Printf("\n%s %s\n", c.Yellow(comment.Author.DisplayName), c.Gray(comment.Created.Format("2006-01-02 15:04")))

		bodyText := comment.GetBodyText()
		if bodyText != "" {
			wrappedText := wrapText(bodyText, 78)
			for _, line := range strings.Split(wrappedText, "\n") {
				fmt.Printf("  %s\n", line)
			}
		} else {
			fmt.Printf("  %s\n", c.Gray("(Empty comment)"))
		}
		fmt.Println()
	}
}

func wrapText(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
