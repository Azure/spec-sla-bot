package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type EmailAssignee struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	EmailID    uuid.UUID `json:"email_id" db:"email_id"`
	AssigneeID string    `json:"assignee_id" db:"assignee_id"`
	Email      Email     `belongs_to:"emails" db:"-"`
	Assignee   Assignee  `belongs_to:"assignees" db:"-"`
}

// String is not required by pop and may be deleted
func (e EmailAssignee) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// EmailAssignees is not required by pop and may be deleted
type EmailAssignees []EmailAssignee

// String is not required by pop and may be deleted
func (e EmailAssignees) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *EmailAssignee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
	//&validators.StringIsPresent{Field: e.EmailID, Name: "EmailID"},
	//&validators.StringIsPresent{Field: e.AssigneeID, Name: "AssigneeID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *EmailAssignee) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *EmailAssignee) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
