package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Assignee struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Login         string    `json:"login" db:"login"`
	Type          string    `json:"type" db:"type"`
	HtmlUrl       string    `json:"html_url" db:"html_url"`
	AssigneeEmail Email     `many_to_many:"email_assignees"`
}

// String is not required by pop and may be deleted
func (a Assignee) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Assignees is not required by pop and may be deleted
type Assignees []Assignee

// String is not required by pop and may be deleted
func (a Assignees) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Assignee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Login, Name: "Login"},
		&validators.StringIsPresent{Field: a.Type, Name: "Type"},
		&validators.StringIsPresent{Field: a.HtmlUrl, Name: "HtmlUrl"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Assignee) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Assignee) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
