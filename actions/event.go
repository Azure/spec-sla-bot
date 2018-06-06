package actions

import (
	"net/http"

	"github.com/Azure/spec-sla-bot/email"
	"github.com/gobuffalo/buffalo"
)

// EventListen default implementation.
func EventListen(c buffalo.Context) error {
	errEmail := email.SendEmailToAssignee()
	if errEmail != nil {
		c.Error(http.StatusInternalServerError, errEmail)
		return errEmail
	}
	return c.Render(200, r.JSON(map[string]string{"message": "Welcome to Buffalo!"}))
}
