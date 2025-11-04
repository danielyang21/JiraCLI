package api

import "time"

// Issue represents a Jira issue
type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields represents the fields of a Jira issue
type IssueFields struct {
	Summary     string      `json:"summary"`
	Description interface{} `json:"description"` // Can be string or complex object
	IssueType   IssueType   `json:"issuetype"`
	Status      Status      `json:"status"`
	Priority    Priority    `json:"priority"`
	Assignee    *User       `json:"assignee"`
	Reporter    *User       `json:"reporter"`
	Created     time.Time   `json:"created"`
	Updated     time.Time   `json:"updated"`
	Project     Project     `json:"project"`
}

// IssueType represents the type of an issue
type IssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Status represents the status of an issue
type Status struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Priority represents the priority of an issue
type Priority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// User represents a Jira user
type User struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	Active      bool   `json:"active"`
}

// Project represents a Jira project
type Project struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// SearchResults represents the results of a JQL search
type SearchResults struct {
	Expand     string  `json:"expand"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// Comment represents a comment on an issue
type Comment struct {
	ID      string    `json:"id"`
	Body    string    `json:"body"`
	Author  User      `json:"author"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// Transition represents a status transition
type Transition struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	To   Status `json:"to"`
}

// TransitionResponse represents available transitions for an issue
type TransitionResponse struct {
	Expand      string       `json:"expand"`
	Transitions []Transition `json:"transitions"`
}
