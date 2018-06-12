package actions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/buffalo-azure/sdk/eventgrid"
	"github.com/gobuffalo/buffalo"
	"github.com/google/go-github/github"
)

// MySpecslaSubscriber gathers responds to all Requests sent to a particular endpoint.
type SpecslaSubscriber struct {
	eventgrid.Subscriber
}

// NewSpecslaSubscriber instantiates SpecslaSubscriber for use in a `buffalo.App`.
func NewSpecslaSubscriber(parent eventgrid.Subscriber) (created *SpecslaSubscriber) {
	dispatcher := eventgrid.NewTypeDispatchSubscriber(parent)

	created = &SpecslaSubscriber{
		Subscriber: dispatcher,
	}

	dispatcher.Bind("Github.PullRequest", created.ReceivePullRequest)

	return
}

// ReceivePullRequest will respond to an `eventgrid.Event` carrying a serialized `PullRequest` as its payload.
func (s *SpecslaSubscriber) ReceivePullRequest(c buffalo.Context, e eventgrid.Event) error {
	var payload github.PullRequest
	if err := json.Unmarshal(e.Data, &payload); err != nil {
		return c.Error(http.StatusBadRequest, errors.New("unable to unmarshal request data"))
	}

	// Replace the code below with your logic
	return nil
}
