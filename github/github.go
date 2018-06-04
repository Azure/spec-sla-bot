package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//PullRequestURL URL for pull request api
const PullRequestURL = "https://api.github.com/repos/Azure/azure-rest-api-specs/pulls"

type PullRequestsResult struct {
	// TotalCount int
	Items []*Request
}

type Request struct {
	Number     int
	HTMLURL    string `json:"html_url"`
	ID         int
	Title      string
	State      string
	User       *User
	Assignee   *Assignee
	Assignees  Assignees
	CreatedAt  time.Time `json:"created_at"`
	UpadatedAt time.Time `json:"updated_at"`
	Labels     Labels
	ClosedAt   time.Time `json:"closed_at"`
	MergedAt   time.Time `json:"merged_at"`
	Body       string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Assignee struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Label struct {
	ID   int
	Name string
}

type Assignees []*Assignee

type Labels []*Label

func PullRequests() (*PullRequestsResult, error) {
	resp, err := http.Get(PullRequestURL)
	if err != nil {
		return nil, err
	}
	//Pushed onto the call stack
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result PullRequestsResult
	if err := json.NewDecoder(resp.Body).Decode(&result.Items); err != nil {
		//resp.Body.Close()
		return nil, err
	}
	//resp.Body.Close()
	return &result, nil
}
