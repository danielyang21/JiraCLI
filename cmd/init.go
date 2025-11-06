package cmd

import (
	"fmt"
	"os"

	"github.com/danielyan21/JiraCLI/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize JiraCLI configuration",
	Long: `Initialize JiraCLI by setting up your Jira credentials and preferences.

This command will guide you through setting up:
- Jira instance URL
- Authentication (email + API token)
- Default project
- Other preferences`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing JiraCLI configuration...")

		if err := config.InitializeConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✓ Configuration initialized successfully!")
		fmt.Println("\nNext steps:")
		fmt.Println("  • Run 'jira list' to view your tickets")
		fmt.Println("  • Run 'jira mine' to see tickets assigned to you")
		fmt.Println("  • Run 'jira --help' for more commands")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
