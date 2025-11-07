package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) UpdateIssueStatus(issueKey, status string) error {
	transitions, err := c.GetTransitions(issueKey)
	if err != nil {
		return fmt.Errorf("fetching transitions: %w", err)
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
		return fmt.Errorf("no matching transition found for '%s'. Available: %s",
			status, strings.Join(availableStatuses, ", "))
	}

	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/transitions", c.getAPIVersion(), issueKey)
	requestBody, err := json.Marshal(map[string]interface{}{
		"transition": map[string]string{
			"id": matchedTransition.ID,
		},
	})
	if err != nil {
		return fmt.Errorf("marshaling transition: %w", err)
	}

	resp, err := c.doRequest("POST", endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return checkResponse(resp)
}

func (c *Client) GetTransitions(issueKey string) (*TransitionResponse, error) {
	endpoint := fmt.Sprintf("/rest/api/%s/issue/%s/transitions", c.getAPIVersion(), issueKey)
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var transitions TransitionResponse
	return &transitions, decodeJSON(resp, &transitions)
}
