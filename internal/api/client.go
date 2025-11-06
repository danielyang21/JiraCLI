package api

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	Email      string
	APIToken   string
	AuthType   string // "basic" or "pat"
	HTTPClient *http.Client
}

func NewClient(baseURL, email, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Email:    email,
		APIToken: apiToken,
		AuthType: "basic",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func NewClientWithAuthType(baseURL, email, apiToken, authType string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Email:    email,
		APIToken: apiToken,
		AuthType: authType,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "JiraCLI/0.1.0")

	if c.AuthType == "pat" {
		req.Header.Set("Authorization", "Bearer "+c.APIToken)
	} else {
		req.SetBasicAuth(c.Email, c.APIToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}

	return resp, nil
}

func (c *Client) TestConnection() error {
	// Use API v3 for Jira Cloud (basic auth), v2 for Jira Server/DC (PAT)
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}
	resp, err := c.doRequest("GET", "/rest/api/"+apiVersion+"/myself", nil)
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
