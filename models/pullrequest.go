package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Pullrequest struct {
	ID               uuid.UUID `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	URL              string    `json:"url" db:"url"`
	HtmlUrl          string    `json:"html_url" db:"html_url"`
	IssueUrl         string    `json:"issue_url" db:"issue_url"`
	Number           string    `json:"number" db:"number"`
	State            string    `json:"state" db:"state"`
	Locked           string    `json:"locked" db:"locked"`
	Title            string    `json:"title" db:"title"`
	Body             string    `json:"body" db:"body"`
	RequestCreatedAt string    `json:"request_created_at" db:"request_created_at"`
	RequestUpdatedAt string    `json:"request_updated_at" db:"request_updated_at"`
	RequestMergedAt  string    `json:"request_merged_at" db:"request_merged_at"`
	RequestClosedAt  string    `json:"request_closed_at" db:"request_closed_at"`
	CommitsUrl       string    `json:"commits_url" db:"commits_url"`
	StatusUrl        string    `json:"status_url" db:"status_url"`
	ExpireTime       string    `json:"expire_time" db:"expire_time"`
}

// String is not required by pop and may be deleted
func (p Pullrequest) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Pullrequests is not required by pop and may be deleted
type Pullrequests []Pullrequest

// String is not required by pop and may be deleted
func (p Pullrequests) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Pullrequest) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.URL, Name: "URL"},
		&validators.StringIsPresent{Field: p.HtmlUrl, Name: "HtmlUrl"},
		&validators.StringIsPresent{Field: p.IssueUrl, Name: "IssueUrl"},
		&validators.StringIsPresent{Field: p.Number, Name: "Number"},
		&validators.StringIsPresent{Field: p.State, Name: "State"},
		&validators.StringIsPresent{Field: p.Locked, Name: "Locked"},
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Body, Name: "Body"},
		&validators.StringIsPresent{Field: p.RequestCreatedAt, Name: "RequestCreatedAt"},
		&validators.StringIsPresent{Field: p.RequestUpdatedAt, Name: "RequestUpdatedAt"},
		&validators.StringIsPresent{Field: p.RequestMergedAt, Name: "RequestMergedAt"},
		&validators.StringIsPresent{Field: p.RequestClosedAt, Name: "RequestClosedAt"},
		&validators.StringIsPresent{Field: p.CommitsUrl, Name: "CommitsUrl"},
		&validators.StringIsPresent{Field: p.StatusUrl, Name: "StatusUrl"},
		&validators.StringIsPresent{Field: p.ExpireTime, Name: "ExpireTime"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Pullrequest) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Pullrequest) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
