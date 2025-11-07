package api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Client) CreateIssue(projectKey, summary, description, issueType, priority string, assignToMe bool) (*CreateIssueResponse, error) {
	fields := CreateIssueFields{
		Project: ProjectRef{
			Key: projectKey,
		},
		Summary: summary,
		IssueType: IssueTypeRef{
			Name: issueType,
		},
	}

	if description != "" {
		if c.AuthType == "pat" {
			fields.Description = description
		} else {
			fields.Description = map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []map[string]interface{}{
					{
						"type": "paragraph",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": description,
							},
						},
					},
				},
			}
		}
	}

	if priority != "" {
		fields.Priority = &PriorityRef{Name: priority}
	}

	if assignToMe {
		currentUser, err := c.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("getting current user: %w", err)
		}

		if c.AuthType == "pat" {
			fields.Assignee = &AssigneeRef{Name: currentUser.DisplayName}
		} else {
			fields.Assignee = &AssigneeRef{AccountID: currentUser.AccountID}
		}
	}

	requestBody, err := json.Marshal(CreateIssueRequest{Fields: fields})
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue", c.getAPIVersion())
	resp, err := c.doRequest("POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CreateIssueResponse
	return &result, decodeJSON(resp, &result)
}
