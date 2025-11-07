package api

import "fmt"

func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s", c.getAPIVersion(), issueKey)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issue Issue
	return &issue, decodeJSON(resp, &issue)
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

	var results SearchResults
	return &results, decodeJSON(resp, &results)
}
