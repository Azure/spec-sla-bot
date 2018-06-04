package ack

import (
	"playground/listpr/github"
	"time"
)

func determineAck(request github.Request) bool {
	//return request.Assignee == nil
	const sla = 24 * time.Hour
	return request.UpadatedAt.Sub(request.CreatedAt) > sla
}

func getUnacknowledgedPR(pullRequests *github.PullRequestsResult) []*github.Request {
	var unacknowledged []*github.Request
	for _, item := range pullRequests.Items {
		if determineAck(*item) {
			unacknowledged = append(unacknowledged, item)
		}
	}
	return unacknowledged
}

func filterByAssignee() {

}
