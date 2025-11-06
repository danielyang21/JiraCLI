package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/danielyan21/JiraCLI/internal/api"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

type Config struct {
	JiraURL        string `mapstructure:"jira_url"`
	Email          string `mapstructure:"email"`
	APIToken       string `mapstructure:"api_token"`
	AuthType       string `mapstructure:"auth_type"` // "basic", "pat", "bearer"
	DefaultProject string `mapstructure:"default_project"`
}

// InitializeConfig prompts the user for configuration values and saves them
func InitializeConfig() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Jira URL (e.g., https://yourcompany.atlassian.net): ")
	jiraURL, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading Jira URL: %w", err)
	}
	jiraURL = strings.TrimSpace(jiraURL)

	// Ask for auth type
	fmt.Println("\nAuthentication method:")
	fmt.Println("  1. Email + API Token (Jira Cloud)")
	fmt.Println("  2. Personal Access Token / PAT (Jira Server/DC)")
	fmt.Println("  3. Username + Password (Basic Auth)")
	fmt.Print("Select (1, 2, or 3) [default: 1]: ")
	authChoice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading auth choice: %w", err)
	}
	authChoice = strings.TrimSpace(authChoice)
	if authChoice == "" {
		authChoice = "1"
	}

	var authType, email, apiToken string

	if authChoice == "2" {
		// PAT authentication
		authType = "pat"
		fmt.Println("\nUsing Personal Access Token authentication")
		fmt.Println("To create a PAT, go to: " + jiraURL + "/secure/ViewProfile.jspa")
		fmt.Println("Then click 'Personal Access Tokens' in the sidebar")
		fmt.Print("\nPersonal Access Token (hidden): ")
		apiTokenBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("error reading PAT: %w", err)
		}
		apiToken = strings.TrimSpace(string(apiTokenBytes))
		fmt.Println()
	} else if authChoice == "3" {
		// Username + Password authentication
		authType = "basic"
		fmt.Println("\nUsing Username + Password authentication")
		fmt.Print("Username: ")
		email, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading username: %w", err)
		}
		email = strings.TrimSpace(email)

		fmt.Print("Password (hidden): ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("error reading password: %w", err)
		}
		apiToken = strings.TrimSpace(string(passwordBytes))
		fmt.Println()
	} else {
		// Email + API Token authentication (default)
		authType = "basic"
		fmt.Println("\nUsing Email + API Token authentication")
		fmt.Print("Email: ")
		email, err = reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading email: %w", err)
		}
		email = strings.TrimSpace(email)

		fmt.Print("API Token (hidden): ")
		apiTokenBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("error reading API token: %w", err)
		}
		apiToken = strings.TrimSpace(string(apiTokenBytes))
		fmt.Println()
	}


	fmt.Print("Default project key (optional, e.g., PROJ): ")
	defaultProject, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading default project: %w", err)
	}
	defaultProject = strings.TrimSpace(defaultProject)

	// Set values in viper
	viper.Set("jira_url", jiraURL)
	viper.Set("auth_type", authType)
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

	// For PAT auth, we don't need email
	if cfg.AuthType != "pat" {
		if cfg.Email == "" {
			return fmt.Errorf("email is required for basic authentication")
		}
	}

	if cfg.APIToken == "" {
		return fmt.Errorf("api_token is required")
	}
	return nil
}

// LoadAndValidate loads and validates the configuration, exiting on error
func LoadAndValidate() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
		os.Exit(1)
	}

	if err := ValidateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run 'jira init' to set up your configuration")
		os.Exit(1)
	}

	return cfg
}

// NewAPIClient creates a new API client from config
func (cfg *Config) NewAPIClient() *api.Client {
	authType := cfg.AuthType
	if authType == "" {
		authType = "basic" // default to basic auth for backwards compatibility
	}
	return api.NewClientWithAuthType(cfg.JiraURL, cfg.Email, cfg.APIToken, authType)
}
