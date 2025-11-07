package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) AssignIssue(issueKey, assignee string) error {
	normalizedAssignee := strings.ToLower(strings.TrimSpace(assignee))
	if normalizedAssignee == "@me" || normalizedAssignee == "me" {
		currentUser, err := c.GetCurrentUser()
		if err != nil {
			return fmt.Errorf("getting current user: %w", err)
		}

		if c.AuthType == "pat" {
			assignee = currentUser.DisplayName
		} else {
			assignee = currentUser.AccountID
		}
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/assignee", c.getAPIVersion(), issueKey)

	var requestBody []byte
	var err error

	if c.AuthType == "pat" {
		requestBody, err = json.Marshal(map[string]string{
			"name": assignee,
		})
	} else {
		requestBody, err = json.Marshal(map[string]string{
			"accountId": assignee,
		})
	}
	if err != nil {
		return fmt.Errorf("marshaling assignee: %w", err)
	}

	resp, err := c.doRequest("PUT", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

func (c *Client) GetCurrentUser() (*User, error) {
	endpoint := fmt.Sprintf("/rest/api/%s/myself", c.getAPIVersion())
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	return &user, decodeJSON(resp, &user)
}
