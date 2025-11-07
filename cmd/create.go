package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/danielyan21/JiraCLI/internal/ui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Jira issue",
	Long: `Create a new Jira issue interactively.

The command will prompt you for all required information with
arrow key navigation and dropdown selections.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadAndValidate()
		client := cfg.NewAPIClient()

		var project string
		projectPrompt := &survey.Input{
			Message: "Project:",
			Default: cfg.DefaultProject,
		}
		survey.AskOne(projectPrompt, &project, survey.WithValidator(survey.Required))

		var summary string
		summaryPrompt := &survey.Input{
			Message: "Summary:",
		}
		survey.AskOne(summaryPrompt, &summary, survey.WithValidator(survey.Required))

		var issueType string
		issueTypePrompt := &survey.Select{
			Message: "Issue Type:",
			Options: []string{"Task", "Bug", "Story", "Epic", "Subtask"},
			Default: "Task",
		}
		survey.AskOne(issueTypePrompt, &issueType)

		var priority string
		priorityPrompt := &survey.Select{
			Message: "Priority:",
			Options: []string{"None", "Highest", "High", "Medium", "Low", "Lowest"},
			Default: "None",
		}
		survey.AskOne(priorityPrompt, &priority)
		if priority == "None" {
			priority = ""
		}

		var description string
		descriptionPrompt := &survey.Multiline{
			Message: "Description (optional):",
		}
		survey.AskOne(descriptionPrompt, &description)

		var assignToMe bool
		assignPrompt := &survey.Confirm{
			Message: "Assign to yourself?",
			Default: false,
		}
		survey.AskOne(assignPrompt, &assignToMe)

		fmt.Printf("\nCreating issue in %s...\n", project)
		result, err := client.CreateIssue(project, summary, description, issueType, priority, assignToMe)
		ui.FatalIfError(err, "Error creating issue")

		fmt.Printf("\n✅ Issue created successfully!\n")
		fmt.Printf("   Key: %s\n", result.Key)
		fmt.Printf("   URL: %s/browse/%s\n", cfg.JiraURL, result.Key)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
