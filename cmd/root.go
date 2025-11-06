package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jira",
	Short: "A high-performance CLI tool for interacting with Jira",
	Long: `JiraCLI is a fast, intuitive command-line interface for Jira that eliminates
context-switching and streamlines developer workflows.

Features:
- Quick ticket updates (status, comments, assignments)
- Advanced search with JQL support
- Git integration for automatic ticket linking
- Local caching for instant responses
- Terminal UI for interactive workflows`,
	Version: "0.1.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jira-cli.yaml)")
	rootCmd.PersistentFlags().String("profile", "default", "configuration profile to use")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("json", "j", false, "output in JSON format")

	// Bind flags to viper
	viper.BindPFlag("profile", rootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jira-cli")
	}

	viper.SetEnvPrefix("JIRA")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
