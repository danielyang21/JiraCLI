package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) UpdateIssueStatus(issueKey, status string) error {
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}

	transitions, err := c.GetTransitions(issueKey)
	if err != nil {
		return fmt.Errorf("error fetching transitions: %w", err)
	}

	statusAliases := map[string][]string{
		"To Do":       {"todo", "td", "t", "to do"},
		"In Progress": {"inprogress", "progress", "ip", "p", "in progress"},
		"Done":        {"done", "d", "complete", "completed"},
		"In Review":   {"review", "r", "inreview", "in review"},
		"Blocked":     {"blocked", "b", "block"},
		"Backlog":     {"backlog", "bl"},
	}

	normalizedInput := strings.ToLower(strings.TrimSpace(status))
	normalizedInput = strings.ReplaceAll(normalizedInput, " ", "")

	var matchedTransition *Transition
	for i := range transitions.Transitions {
		t := &transitions.Transitions[i]
		normalizedTransition := strings.ToLower(t.To.Name)

		if normalizedTransition == strings.ToLower(status) {
			matchedTransition = t
			break
		}

		if aliases, ok := statusAliases[t.To.Name]; ok {
			for _, alias := range aliases {
				if alias == normalizedInput {
					matchedTransition = t
					break
				}
			}
			if matchedTransition != nil {
				break
			}
		}

		if strings.Contains(normalizedTransition, normalizedInput) {
			matchedTransition = t
		}
	}

	if matchedTransition == nil {
		availableStatuses := make([]string, len(transitions.Transitions))
		for i, t := range transitions.Transitions {
			availableStatuses[i] = t.To.Name
		}
		return fmt.Errorf("no matching transition found for '%s'. Available transitions: %s",
			status, strings.Join(availableStatuses, ", "))
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/transitions", apiVersion, issueKey)
	requestBody, err := json.Marshal(map[string]interface{}{
		"transition": map[string]string{
			"id": matchedTransition.ID,
		},
	})
	if err != nil {
		return fmt.Errorf("error creating request body: %w", err)
	}

	bodyReader := bytes.NewReader(requestBody)
	resp, err := c.doRequest("POST", endpoint, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update status (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) GetTransitions(issueKey string) (*TransitionResponse, error) {
	apiVersion := "3"
	if c.AuthType == "pat" {
		apiVersion = "2"
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/transitions", apiVersion, issueKey)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get transitions (status %d): %s", resp.StatusCode, string(body))
	}

	var transitions TransitionResponse
	if err := json.NewDecoder(resp.Body).Decode(&transitions); err != nil {
		return nil, fmt.Errorf("error decoding transitions: %w", err)
	}

	return &transitions, nil
}
