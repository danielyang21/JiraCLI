package api

import (
	"strings"
	"time"
)

type JiraTime struct {
	time.Time
}

func (jt *JiraTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		jt.Time = time.Time{}
		return nil
	}

	// Try standard RFC3339 format first
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		jt.Time = t
		return nil
	}

	// Try Jira's format: 2025-11-05T22:21:49.975-0500
	// Convert -0500 to -05:00
	if len(s) >= 5 && (s[len(s)-5] == '+' || s[len(s)-5] == '-') {
		s = s[:len(s)-2] + ":" + s[len(s)-2:]
	}

	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	jt.Time = t
	return nil
}

type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string      `json:"summary"`
	Description interface{} `json:"description"`
	IssueType   IssueType   `json:"issuetype"`
	Status      Status      `json:"status"`
	Priority    Priority    `json:"priority"`
	Assignee    *User       `json:"assignee"`
	Reporter    *User       `json:"reporter"`
	Created     JiraTime    `json:"created"`
	Updated     JiraTime    `json:"updated"`
	Project     Project     `json:"project"`
}

type IssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Status struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Priority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
}

type Project struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type SearchResults struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

type Comment struct {
	ID      string    `json:"id"`
	Body    string    `json:"body"`
	Author  User      `json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   Status `json:"to"`
}

type TransitionResponse struct {
	Expand      string       `json:"expand"`
	Transitions []Transition `json:"transitions"`
}

type CreateIssueRequest struct {
	Fields CreateIssueFields `json:"fields"`
}

type CreateIssueFields struct {
	Project     ProjectRef     `json:"project"`
	Summary     string         `json:"summary"`
	Description interface{}    `json:"description,omitempty"`
	IssueType   IssueTypeRef   `json:"issuetype"`
	Priority    *PriorityRef   `json:"priority,omitempty"`
	Assignee    *AssigneeRef   `json:"assignee,omitempty"`
}

type ProjectRef struct {
	Key string `json:"key"`
}

type IssueTypeRef struct {
	Name string `json:"name"`
}

type PriorityRef struct {
	Name string `json:"name"`
}

type AssigneeRef struct {
	AccountID string `json:"accountId,omitempty"`
	Name      string `json:"name,omitempty"`
}

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}
