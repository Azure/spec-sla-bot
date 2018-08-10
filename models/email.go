package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Email struct {
	ID            uuid.UUID `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	PullrequestID int       `json:"pullrequest_id" db:"pullrequest_id"`
	TimeSent      string    `json:"time_sent" db:"time_sent"`
	Assignees     Assignees `many_to_many:"email_assignees"`
}

// String is not required by pop and may be deleted
func (e Email) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Emails is not required by pop and may be deleted
type Emails []Email

// String is not required by pop and may be deleted
func (e Emails) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (e *Email) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		//&validators.StringIsPresent{Field: e.PullrequestID, Name: "PullrequestID"},
		&validators.StringIsPresent{Field: e.TimeSent, Name: "TimeSent"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (e *Email) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (e *Email) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
