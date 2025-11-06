package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) AssignIssue(issueKey, assignee string) error {
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}

	// Handle @me or me - get current user
	normalizedAssignee := strings.ToLower(strings.TrimSpace(assignee))
	if normalizedAssignee == "@me" || normalizedAssignee == "me" {
		currentUser, err := c.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("error getting current user: %w", err)
		}

		if c.AuthType == "pat" {
			assignee = currentUser.DisplayName
		} else {
			assignee = currentUser.AccountID
		}
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/assignee", apiVersion, issueKey)

	var requestBody []byte
	var err error

	if c.AuthType == "pat" {
		// For v2, use "name" field
		requestBody, err = json.Marshal(map[string]string{
			"name": assignee,
		})
	} else {
		// For v3, use "accountId" field
		requestBody, err = json.Marshal(map[string]string{
			"accountId": assignee,
		})
	}

	if err != nil {
		return fmt.Errorf("error creating request body: %w", err)
	}

	bodyReader := bytes.NewReader(requestBody)
	resp, err := c.doRequest("PUT", endpoint, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to assign issue (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetCurrentUser() (*User, error) {
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}

	endpoint := fmt.Sprintf("/rest/api/%s/myself", apiVersion)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get current user (status %d): %s", resp.StatusCode, string(body))
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("error decoding user: %w", err)
	}

	return &user, nil
}
