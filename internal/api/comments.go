package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) AddComment(issueKey, comment string) error {
	apiVersion := "3"
	var requestBody []byte
	var err error

	if c.AuthType == "pat" {
		apiVersion = "2"
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
		return fmt.Errorf("error creating request body: %w", err)
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/comment", apiVersion, issueKey)
	bodyReader := bytes.NewReader(requestBody)

	resp, err := c.doRequest("POST", endpoint, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to add comment (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
