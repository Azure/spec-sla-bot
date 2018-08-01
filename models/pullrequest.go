package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

/*type ValidTime struct {
	Time  time.Time
	Valid bool
}*/

type Pullrequest struct {
	GitPRID          int64     `json:"git_prid" db:"git_prid"`
	ID               uuid.UUID `json:"id" db:"id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	URL              string    `json:"url" db:"url"`
	HtmlUrl          string    `json:"html_url" db:"html_url"`
	IssueUrl         string    `json:"issue_url" db:"issue_url"`
	Number           int       `json:"number" db:"number"`
	State            string    `json:"state" db:"state"`
	ValidTime        bool      `json:"valid_time" db:"valid_time"`
	Title            string    `json:"title" db:"title"`
	Body             string    `json:"body" db:"body"`
	RequestCreatedAt time.Time `json:"request_created_at" db:"request_created_at"`
	RequestUpdatedAt time.Time `json:"request_updated_at" db:"request_updated_at"`
	RequestMergedAt  time.Time `json:"request_merged_at" db:"request_merged_at"`
	RequestClosedAt  time.Time `json:"request_closed_at" db:"request_closed_at"`
	CommitsUrl       int       `json:"commits_url" db:"commits_url"`
	StatusUrl        string    `json:"status_url" db:"status_url"`
	ExpireTime       time.Time `json:"expire_time" db:"expire_time"`
	Assignees        Assignees `many_to_many:"pullrequest_assignees" db:"-"`
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
		&validators.IntIsPresent{Field: p.Number, Name: "Number"},
		&validators.StringIsPresent{Field: p.State, Name: "State"},
		//&validators.StringIsPresent{Field: p.ValidTime, Name: "ValidTime"},
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Body, Name: "Body"},
		//&validators.StringIsPresent{Field: p.RequestCreatedAt, Name: "RequestCreatedAt"},
		//&validators.StringIsPresent{Field: p.RequestUpdatedAt, Name: "RequestUpdatedAt"},
		//&validators.StringIsPresent{Field: p.RequestMergedAt, Name: "RequestMergedAt"},
		//&validators.StringIsPresent{Field: p.RequestClosedAt, Name: "RequestClosedAt"},
		&validators.IntIsPresent{Field: p.CommitsUrl, Name: "CommitsUrl"},
		&validators.StringIsPresent{Field: p.StatusUrl, Name: "StatusUrl"},
		//&validators.StringIsPresent{Field: p.ExpireTime, Name: "ExpireTime"},
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
