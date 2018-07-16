package grifts

import (
	"context"
	"time"

	"github.com/Azure/spec-sla-bot/models"
	"github.com/google/go-github/github"
	"github.com/markbates/grift/grift"
)

var ctx = context.Background()
var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seed the database with PRs in the specs repo")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		var client *github.Client
		var opt *github.PullRequestListOptions
		pullRequestList, _, err := client.PullRequests.List(ctx, "t-jaelli", "azure-service-bus-go", opt)
		if err != nil {
			return err
		}
		//var pr models.Pullrequest
		//var a models.Assignee
		if pullRequestList != nil {
			for _, pullRequest := range pullRequestList {
				expireTime := models.ValidTime{
					Time:  time.Now(),
					Valid: false,
				}
				pr := &models.Pullrequest{
					URL:              *pullRequest.URL,
					HtmlUrl:          *pullRequest.HTMLURL,
					IssueUrl:         *pullRequest.IssueURL,
					Number:           *pullRequest.Number,
					State:            *pullRequest.State,
					Locked:           *pullRequest.Title, //not correct. May need a new column for ID or get rid of locked
					Title:            *pullRequest.Title,
					Body:             *pullRequest.Body,
					RequestCreatedAt: *pullRequest.CreatedAt,
					RequestUpdatedAt: *pullRequest.UpdatedAt,
					RequestMergedAt:  *pullRequest.MergedAt,
					RequestClosedAt:  *pullRequest.ClosedAt,
					CommitsUrl:       *pullRequest.Commits,     // may need a null check to get the CommitsURL
					StatusUrl:        *pullRequest.StatusesURL, // consider changing name of column to match statuses
					ExpireTime:       expireTime,
				}
				err := models.DB.Create(pr)
				if err != nil {
					return err
				}
				if pullRequest.Assignee != nil {
					//add assignee info to the database
					a := &models.Assignee{
						Login:   *pullRequest.Assignee.Login,
						Type:    *pullRequest.Assignee.Type,
						HtmlUrl: *pullRequest.Assignee.HTMLURL,
					}
					err := models.DB.Create(a)
					if err != nil {
						return err
					}

				}
			}
		}
		return nil
	})

})
