package api

import (
	"encoding/json"
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

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}

func decodeJSON(resp *http.Response, v interface{}) error {
	if err := checkResponse(resp); err != nil {
		return err
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) getAPIVersion() string {
	if c.AuthType == "pat" {
		return "2"
	}
	return "3"
}

func (c *Client) TestConnection() error {
	apiVersion := c.getAPIVersion()
	resp, err := c.doRequest("GET", "/rest/api/"+apiVersion+"/myself", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return checkResponse(resp)
}
