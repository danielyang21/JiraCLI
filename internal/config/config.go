package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

// Config represents the application configuration
type Config struct {
	JiraURL      string `mapstructure:"jira_url"`
	Email        string `mapstructure:"email"`
	APIToken     string `mapstructure:"api_token"`
	DefaultProject string `mapstructure:"default_project"`
}

// InitializeConfig prompts the user for configuration values and saves them
func InitializeConfig() error {
	reader := bufio.NewReader(os.Stdin)

	// Jira URL
	fmt.Print("Jira URL (e.g., https://yourcompany.atlassian.net): ")
	jiraURL, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading Jira URL: %w", err)
	}
	jiraURL = strings.TrimSpace(jiraURL)

	// Email
	fmt.Print("Email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading email: %w", err)
	}
	email = strings.TrimSpace(email)

	// API Token (hidden input)
	fmt.Print("API Token (hidden): ")
	apiTokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("error reading API token: %w", err)
	}
	apiToken := strings.TrimSpace(string(apiTokenBytes))
	fmt.Println() // New line after hidden input

	// Default Project (optional)
	fmt.Print("Default project key (optional, e.g., PROJ): ")
	defaultProject, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading default project: %w", err)
	}
	defaultProject = strings.TrimSpace(defaultProject)

	// Set values in viper
	viper.Set("jira_url", jiraURL)
	viper.Set("email", email)
	viper.Set("api_token", apiToken)
	viper.Set("default_project", defaultProject)

	// Get config file path
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %w", err)
	}

	configPath := home + "/.jira-cli.yaml"

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	// Set restrictive permissions on config file (600)
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("error setting config file permissions: %w", err)
	}

	fmt.Printf("\nConfiguration saved to: %s\n", configPath)
	return nil
}

// LoadConfig loads and returns the configuration
func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &cfg, nil
}

// ValidateConfig checks if the configuration is valid
func ValidateConfig(cfg *Config) error {
	if cfg.JiraURL == "" {
		return fmt.Errorf("jira_url is required")
	}
	if cfg.Email == "" {
		return fmt.Errorf("email is required")
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("api_token is required")
	}
	return nil
}
