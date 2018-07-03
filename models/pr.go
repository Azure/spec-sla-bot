package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type PR struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	URL         string    `json:"url" db:"url"`
	HtmlUrl     string    `json:"html_url" db:"html_url"`
	IssueUrl    string    `json:"issue_url" db:"issue_url"`
	Number      string    `json:"number" db:"number"`
	State       string    `json:"state" db:"state"`
	Locked      string    `json:"locked" db:"locked"`
	Title       string    `json:"title" db:"title"`
	Body        string    `json:"body" db:"body"`
	PrCreatedAt string    `json:"pr_created_at" db:"pr_created_at"`
	PrUpdatedAt string    `json:"pr_updated_at" db:"pr_updated_at"`
	PrMergedAt  string    `json:"pr_merged_at" db:"pr_merged_at"`
	PrClosedAt  string    `json:"pr_closed_at" db:"pr_closed_at"`
	CommitsUrl  string    `json:"commits_url" db:"commits_url"`
	CommentsUrl string    `json:"comments_url" db:"comments_url"`
	StatusUrl   string    `json:"status_url" db:"status_url"`
	ExpireTime  string    `json:"expire_time" db:"expire_time"`
}

// String is not required by pop and may be deleted
func (p PR) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PRs is not required by pop and may be deleted
type PRs []PR

// String is not required by pop and may be deleted
func (p PRs) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PR) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.URL, Name: "URL"},
		&validators.StringIsPresent{Field: p.HtmlUrl, Name: "HtmlUrl"},
		&validators.StringIsPresent{Field: p.IssueUrl, Name: "IssueUrl"},
		&validators.StringIsPresent{Field: p.Number, Name: "Number"},
		&validators.StringIsPresent{Field: p.State, Name: "State"},
		&validators.StringIsPresent{Field: p.Locked, Name: "Locked"},
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Body, Name: "Body"},
		&validators.StringIsPresent{Field: p.PrCreatedAt, Name: "PrCreatedAt"},
		&validators.StringIsPresent{Field: p.PrUpdatedAt, Name: "PrUpdatedAt"},
		&validators.StringIsPresent{Field: p.PrMergedAt, Name: "PrMergedAt"},
		&validators.StringIsPresent{Field: p.PrClosedAt, Name: "PrClosedAt"},
		&validators.StringIsPresent{Field: p.CommitsUrl, Name: "CommitsUrl"},
		&validators.StringIsPresent{Field: p.CommentsUrl, Name: "CommentsUrl"},
		&validators.StringIsPresent{Field: p.StatusUrl, Name: "StatusUrl"},
		&validators.StringIsPresent{Field: p.ExpireTime, Name: "ExpireTime"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PR) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PR) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
