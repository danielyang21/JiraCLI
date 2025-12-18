package config

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/danielyan21/JiraCLI/internal/api"
	"github.com/spf13/viper"
)

type Config struct {
	JiraURL        string `mapstructure:"jira_url"`
	Email          string `mapstructure:"email"`
	APIToken       string `mapstructure:"api_token"`
	AuthType       string `mapstructure:"auth_type"` // "basic", "pat", "bearer"
	DefaultProject string `mapstructure:"default_project"`
}

func InitializeConfig() error {
	var jiraURL string
	urlPrompt := &survey.Input{
		Message: "Jira URL:",
		Help:    "e.g., https://yourcompany.atlassian.net",
	}
	if err := survey.AskOne(urlPrompt, &jiraURL, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var authMethod string
	authPrompt := &survey.Select{
		Message: "Authentication method:",
		Options: []string{
			"Email + API Token (Jira Cloud)",
			"Personal Access Token (Jira Server/DC)",
			"Username + Password (Basic Auth)",
		},
		Default: "Email + API Token (Jira Cloud)",
	}
	if err := survey.AskOne(authPrompt, &authMethod); err != nil {
		return err
	}

	var authType, email, apiToken string

	switch authMethod {
	case "Personal Access Token (Jira Server/DC)":
		authType = "pat"
		fmt.Printf("\nTo create a PAT, go to: %s/secure/ViewProfile.jspa\n", jiraURL)
		fmt.Println("Then click 'Personal Access Tokens' in the sidebar\n")

		patPrompt := &survey.Password{
			Message: "Personal Access Token:",
		}
		if err := survey.AskOne(patPrompt, &apiToken, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

	case "Username + Password (Basic Auth)":
		authType = "basic"

		usernamePrompt := &survey.Input{
			Message: "Username:",
		}
		if err := survey.AskOne(usernamePrompt, &email, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		passwordPrompt := &survey.Password{
			Message: "Password:",
		}
		if err := survey.AskOne(passwordPrompt, &apiToken, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

	default:
		authType = "basic"

		emailPrompt := &survey.Input{
			Message: "Email:",
		}
		if err := survey.AskOne(emailPrompt, &email, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		tokenPrompt := &survey.Password{
			Message: "API Token:",
			Help:    "Create one at: " + jiraURL + "/secure/ViewProfile.jspa?selectedTab=com.atlassian.pats.pats-plugin:jira-user-personal-access-tokens",
		}
		if err := survey.AskOne(tokenPrompt, &apiToken, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	var defaultProject string
	projectPrompt := &survey.Input{
		Message: "Default project key (optional):",
		Help:    "e.g., PROJ, KAN",
	}
	survey.AskOne(projectPrompt, &defaultProject)

	viper.Set("jira_url", jiraURL)
	viper.Set("auth_type", authType)
	viper.Set("email", email)
	viper.Set("api_token", apiToken)
	viper.Set("default_project", defaultProject)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}

	configPath := home + "/.jira-cli.yaml"

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("setting config file permissions: %w", err)
	}

	fmt.Printf("\nConfiguration saved to: %s\n", configPath)
	return nil
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	if cfg.JiraURL == "" {
		return fmt.Errorf("jira_url is required")
	}

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

func (cfg *Config) NewAPIClient() *api.Client {
	authType := cfg.AuthType
	if authType == "" {
		authType = "basic"
	}
	return api.NewClientWithAuthType(cfg.JiraURL, cfg.Email, cfg.APIToken, authType)
}
