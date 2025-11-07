package api

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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
