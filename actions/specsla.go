package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/buffalo-azure/sdk/eventgrid"
	"github.com/Azure/spec-sla-bot/messages"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/google/go-github/github"
)

// SpecslaSubscriber gathers responses to all Requests sent to a particular endpoint.
type SpecslaSubscriber struct {
	eventgrid.Subscriber
}

// NewSpecslaSubscriber instantiates SpecslaSubscriber for use in a `buffalo.App`.
func NewSpecslaSubscriber(parent eventgrid.Subscriber) (created *SpecslaSubscriber) {
	dispatcher := eventgrid.NewTypeDispatchSubscriber(parent)

	created = &SpecslaSubscriber{
		Subscriber: dispatcher,
	}

	dispatcher.Bind("Github.PullRequestEvent", created.ReceivePullRequestEvent)
	//dispatcher.Bind("Github.LabelEvent", created.ReceiveLabelEvent)

	return
}

// ReceivePullRequestEvent will respond to an `eventgrid.Event` carrying a serialized `PullRequestEvent` as its payload.
func (s *SpecslaSubscriber) ReceivePullRequestEvent(c buffalo.Context, e eventgrid.Event) error {
	var payload github.PullRequestEvent

	if err := json.Unmarshal(e.Data, &payload); err != nil {
		return c.Error(http.StatusBadRequest, errors.New("unable to unmarshal request data"))
	}
	c.Logger().Debug("HERE")
	messages.CheckAcknowledgement(payload)

	// Replace the code below with your logic
	return c.Render(200, render.JSON(map[string]string{"message": "Hopefully this works"}))
}

// ReceiveLabelEvent will respond to an `eventgrid.Event` carrying a serialized `LabelEvent` as its payload.
/*func (s *SpecslaSubscriber) ReceiveLabelEvent(c buffalo.Context, e eventgrid.Event) error {
	var payload github.LabelEvent

	if err := json.Unmarshal(e.Data, &payload); err != nil {
		return c.Error(http.StatusBadRequest, errors.New("unable to unmarshal request data"))
	}

	// Replace the code below with your logic
	return c.Render(200, render.JSON(map[string]string{"message": "Hopefully this works", "label name for last event": *payload.Label.Name}))
}*/
