package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
)

type PullrequestAssignee struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	PullrequestID uuid.UUID   `json:"pullrequest_id" db:"pullrequest_id"`
	AssigneeID    uuid.UUID   `json:"assignee_id" db:"assignee_id"`
	Assignee      Assignee    `belongs_to:"assignees" db:"-"`
	Pullrequest   Pullrequest `belongs_to:"pullrequests" db:"-"`
}

// String is not required by pop and may be deleted
func (p PullrequestAssignee) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// PullrequestAssignees is not required by pop and may be deleted
type PullrequestAssignees []PullrequestAssignee

// String is not required by pop and may be deleted
func (p PullrequestAssignees) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *PullrequestAssignee) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
	//&validators.StringIsPresent{Field: p.PullrequestID, Name: "PullrequestID"},
	//&validators.StringIsPresent{Field: p.AssigneeID, Name: "AssigneeID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *PullrequestAssignee) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *PullrequestAssignee) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
