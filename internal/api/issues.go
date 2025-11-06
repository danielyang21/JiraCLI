package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}
	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s", apiVersion, issueKey)
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

func (c *Client) SearchIssues(jql string, maxResults int) (*SearchResults, error) {
	var endpoint string

	// Jira Cloud (basic auth) uses API v3 with /search/jql endpoint
	// Jira Server/DC (PAT) uses API v2 with /search endpoint
	if c.AuthType == "pat" {
		endpoint = fmt.Sprintf("/rest/api/2/search?jql=%s&maxResults=%d", jql, maxResults)
	} else {
		// Jira Cloud requires /search/jql endpoint (new as of 2024)
		// Must explicitly request fields (default is only "id")
		endpoint = fmt.Sprintf("/rest/api/3/search/jql?jql=%s&maxResults=%d&fields=*navigable", jql, maxResults)
	}
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
