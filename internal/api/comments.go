package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}

func (c *Client) GetComments(issueKey string) ([]Comment, error) {
	apiVersion := c.getAPIVersion()
	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/comment", apiVersion, issueKey)

	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var result CommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling comments: %w", err)
	}

	return result.Comments, nil
}

func (c *Client) AddComment(issueKey, comment string) error {
	apiVersion := c.getAPIVersion()
	var requestBody []byte
	var err error

	if c.AuthType == "pat" {
		requestBody, err = json.Marshal(map[string]string{
			"body": comment,
		})
	} else {
		// API v3 uses Atlassian Document Format (ADF)
		requestBody, err = json.Marshal(map[string]interface{}{
			"body": map[string]interface{}{
				"type":    "doc",
				"version": 1,
				"content": []map[string]interface{}{
					{
						"type": "paragraph",
						"content": []map[string]interface{}{
							{
								"type": "text",
								"text": comment,
							},
						},
					},
				},
			},
		})
	}
	if err != nil {
		return fmt.Errorf("marshaling comment: %w", err)
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/comment", apiVersion, issueKey)
	resp, err := c.doRequest("POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}
