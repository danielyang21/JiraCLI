package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Jira API client
type Client struct {
	BaseURL    string
	Email      string
	APIToken   string
	HTTPClient *http.Client
}

// NewClient creates a new Jira API client
func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Email:    email,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Basic authentication
	req.SetBasicAuth(c.Email, c.APIToken)

	// Perform request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}

	return resp, nil
}

// TestConnection tests the connection to Jira
func (c *Client) TestConnection() error {
	resp, err := c.doRequest("GET", "/rest/api/3/myself", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to connect to Jira (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetIssue retrieves a single issue by key
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	endpoint := fmt.Sprintf("/rest/api/3/issue/%s", issueKey)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get issue (status %d): %s", resp.StatusCode, string(body))
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("error decoding issue: %w", err)
	}

	return &issue, nil
}

// SearchIssues searches for issues using JQL
func (c *Client) SearchIssues(jql string, maxResults int) (*SearchResults, error) {
	endpoint := fmt.Sprintf("/rest/api/3/search?jql=%s&maxResults=%d", jql, maxResults)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to search issues (status %d): %s", resp.StatusCode, string(body))
	}

	var results SearchResults
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding search results: %w", err)
	}

	return &results, nil
}

// AddComment adds a comment to an issue
func (c *Client) AddComment(issueKey, comment string) error {
	// TODO: Implement comment addition
	return fmt.Errorf("not yet implemented")
}

// UpdateIssueStatus updates the status of an issue
func (c *Client) UpdateIssueStatus(issueKey, status string) error {
	// TODO: Implement status transition
	return fmt.Errorf("not yet implemented")
}

// AssignIssue assigns an issue to a user
func (c *Client) AssignIssue(issueKey, accountID string) error {
	// TODO: Implement issue assignment
	return fmt.Errorf("not yet implemented")
}
