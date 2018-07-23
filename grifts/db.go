package grifts

import (
	"context"
	"log"
	"time"

	"github.com/Azure/spec-sla-bot/messages"
	"github.com/Azure/spec-sla-bot/models"
	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var ctx = context.Background()

var _ = grift.Add("seed:without:connection", func(c *grift.Context) error {
	var client *github.Client
	var opt *github.PullRequestListOptions
	pullRequestList, _, err := client.PullRequests.List(ctx, "t-jaelli", "azure-service-bus-go", opt)
	if err != nil {
		return err
	}
	if pullRequestList != nil {
		for _, pullRequest := range pullRequestList {
			pr := &models.Pullrequest{
				GitPRID:          *pullRequest.ID,
				URL:              *pullRequest.URL,
				HtmlUrl:          *pullRequest.HTMLURL,
				IssueUrl:         *pullRequest.IssueURL,
				Number:           *pullRequest.Number,
				State:            *pullRequest.State,
				ValidTime:        false,
				Title:            *pullRequest.Title,
				Body:             *pullRequest.Body,
				RequestCreatedAt: *pullRequest.CreatedAt,
				RequestUpdatedAt: *pullRequest.UpdatedAt,
				RequestMergedAt:  messages.NullCheckTime(pullRequest.MergedAt),
				RequestClosedAt:  messages.NullCheckTime(pullRequest.ClosedAt),
				CommitsUrl:       messages.NullCheckInt(pullRequest.Commits),
				StatusUrl:        *pullRequest.StatusesURL, // consider changing name of column to match statuses
				ExpireTime:       time.Time{},
			}
			err := models.DB.Create(pr)
			if err != nil {
				return err
			}
			if pullRequest.Assignee != nil {
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

var _ = grift.Add("db:seed:truncateAll", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		err := tx.TruncateAll()
		if err != nil {
			return errors.WithStack(err)
		}
		c.Set("tx", tx)
		err = grift.Run("db:seed:without:connection", c)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
})

var _ = grift.Add("db:seed:with:connection", func(c *grift.Context) error {
	var client *github.Client = github.NewClient(nil)
	var opt *github.PullRequestListOptions
	db := models.DB
	if tx := c.Value("tx"); tx != nil {
		log.Printf("Made connection")
		db = tx.(*pop.Connection)
	}
	pullRequestList, _, err := client.PullRequests.List(ctx, "t-jaelli", "azure-rest-api-specs", opt)
	if err != nil {
		log.Printf("Failed to get prlist from github")
		return err
	}
	if pullRequestList != nil {
		log.Printf("list is not nil")
		log.Printf("length: %d", len(pullRequestList))
		for _, pullRequest := range pullRequestList {
			//log.Printf(pullRequest.String())
			pr := &models.Pullrequest{
				GitPRID:          *pullRequest.ID,
				URL:              *pullRequest.URL,
				HtmlUrl:          *pullRequest.HTMLURL,
				IssueUrl:         *pullRequest.IssueURL,
				Number:           *pullRequest.Number,
				State:            *pullRequest.State,
				ValidTime:        false,
				Title:            *pullRequest.Title,
				Body:             *pullRequest.Title,
				RequestCreatedAt: *pullRequest.CreatedAt,
				RequestUpdatedAt: *pullRequest.UpdatedAt,
				RequestMergedAt:  messages.NullCheckTime(pullRequest.MergedAt),
				RequestClosedAt:  messages.NullCheckTime(pullRequest.ClosedAt),
				CommitsUrl:       messages.NullCheckInt(pullRequest.Commits), // may need a null check to get the CommitsURL
				StatusUrl:        *pullRequest.StatusesURL,                   // consider changing name of column to match statuses
				ExpireTime:       time.Time{},
			}
			log.Printf("Made it here")
			log.Printf(pr.String())
			err = db.Create(pr)
			if err != nil {
				log.Printf("Could not create the pr in the database")
				return err
			}
			if pullRequest.Assignee != nil {
				a := &models.Assignee{
					Login:   *pullRequest.Assignee.Login,
					Type:    *pullRequest.Assignee.Type,
					HtmlUrl: *pullRequest.Assignee.HTMLURL,
				}
				err := db.Create(a)
				if err != nil {
					log.Printf("Could not create the assignee in the database")
					return err
				}
			}
		}
	}
	return nil
})

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		c.Set("tx", tx)
		err := grift.Run("db:seed:with:connection", c)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
})
