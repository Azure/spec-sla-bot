package actions

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo/render"

	"github.com/gobuffalo/buffalo"
	"github.com/google/go-github/github"
)

// EventListen default implementation.
func EventListen(c buffalo.Context) error {
	/*errEmail := email.SendEmailToAssignee()
	if errEmail != nil {
		c.Error(http.StatusInternalServerError, errEmail)
		return errEmail
	}*/ //handle secret
	request := c.Request()
	payload, err := github.ValidatePayload(request, []byte(os.Getenv("APPSETTING_X_HUB_SIGNATURE")))
	if err != nil {
		log.Printf("secret key is not correct: err=%s\n", err)
		return c.Error(http.StatusInternalServerError, err)
	}
	defer request.Body.Close()
	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return c.Error(http.StatusInternalServerError, err)
	}
	repoName := ""
	switch e := event.(type) {
	case *github.PullRequestEvent:
		if e.Action != nil {
			repoName = *e.Repo.FullName
			fmt.Printf("Repository Name: %s", *e.Repo.FullName)
		}
	case *github.LabelEvent:
		if e.Action != nil {
			repoName = *e.Repo.FullName
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(request))
		//return err
	}
	return c.Render(200, render.JSON(map[string]string{"message": "Welcome to Buffalo!", "repo name": repoName}))
}
