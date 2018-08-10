package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/spec-sla-bot/models"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/google/go-github/github"
)

//SLAQueue name
const SLAQueue = "24hrgitevents"

//MessageContent sent to service bus
type MessageContent struct {
	PRID    int64
	HTMLURL string
	//This should be an array due to the potential of multiple assignees
	AssigneeLogin        string
	ManagerEmailReminder bool
}

//CheckAcknowledgement determines if a PullRequestEvent is a valid acknowledment of a PR
func CheckAcknowledgement(ctx context.Context, event github.PullRequestEvent) error {
	if checkClosed(ctx, event, models.DB) || checkUnassigned(ctx, event, models.DB) || (event.PullRequest.Assignee == nil && checkOpened(ctx, event, models.DB)) {
		err := UpsertPullRequestEntry(ctx, event, models.DB, false, time.Time{})
		if err != nil {
			return fmt.Errorf("unable to update event number %d in CheckAcknowledgement", *event.Number)
		}
	} else if event.PullRequest.Assignee != nil && (checkAssigned(ctx, event, models.DB) ||
		checkReviewed(ctx, event, models.DB) ||
		checkEdited(ctx, event, models.DB) ||
		checkLabeled(ctx, event, models.DB) ||
		checkOpened(ctx, event, models.DB)) {
		messageStruct := MessageContent{
			PRID:                 *event.PullRequest.ID,
			HTMLURL:              *event.PullRequest.HTMLURL,
			AssigneeLogin:        *event.PullRequest.Assignee.Login,
			ManagerEmailReminder: false,
		}
		message, err := json.Marshal(messageStruct)
		if err != nil {
			fmt.Errorf("Unable to Marshal message struct for PR %d", *event.PullRequest.ID)
			return err
		}
		err = SendToQueue(ctx, message, expireTime(time.Now()), SLAQueue)
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.PullRequest.ID)
			return err
		}
	}
	return nil
}

//CheckAcknowledgementComment determines if an IssuesCommentEvent on a PR is a valid acknowledgement
func CheckAcknowledgementComment(ctx context.Context, event github.IssueCommentEvent) error {
	if event.Issue != nil && event.Issue.IsPullRequest() && checkCommented(ctx, event, models.DB) {
		assignees := event.Issue.Assignees
		for _, assignee := range assignees {
			id, err := uuid.NewV1()
			if err != nil {
				return err
			}
			q := models.DB.RawQuery(`INSERT INTO assignees (id, created_at, updated_at, login, type, html_url)
			VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (login) DO UPDATE SET html_url = ?`,
				id, time.Now(), time.Now(), assignee.Login, assignee.Type, assignee.HTMLURL, assignee.HTMLURL)
			exErr := q.Exec()
			if exErr != nil {
				log.Println("Unable to update event in checkAssigned")
				return exErr
			}
		}
		messageStruct := MessageContent{
			PRID:          *event.Issue.ID,
			HTMLURL:       *event.Issue.HTMLURL,
			AssigneeLogin: *event.Issue.Assignee.Login,
		}
		message, err := json.Marshal(messageStruct)
		if err != nil {
			fmt.Errorf("Unable to Marshal message struct for PR %d", *event.Issue.ID)
			return err
		}
		err = SendToQueue(ctx, message, expireTime(time.Now()), SLAQueue)
		if err != nil {
			return fmt.Errorf("Message for event %d not delivered: %v", *event.Issue.ID, err)
		}
	}
	return nil
}

//CheckAcknowledgementLabel determines if a LabelEvent on a PR is a valid acknowledgement
func CheckAcknowledgementLabel(ctx context.Context, event github.LabelEvent) error {
	panic("not implemented")
}

//ExpireTime calculates the time a PR could violate the SLA depending on the current time
func expireTime(currentTime time.Time) time.Time {
	//return currentTime.Add(slaDuration(currentTime.Weekday()))
	return currentTime.Add(time.Minute * 2)
}

//slaDuration returns the amount of time an assignee has to respond to a PR given the day
func slaDuration(day time.Weekday) time.Duration {
	delay := 24 * time.Hour

	//Adjusted time for now for testing purposes
	//if the current weekday is Friday
	if day > time.Thursday {
		delay = time.Duration(8-day) * time.Hour * 24
	}
	return delay
}

//CheckCommented determines if the comment on the PR was valid
func checkCommented(ctx context.Context, event github.IssueCommentEvent, tx *pop.Connection) bool {
	//check that the issue is not nil and that the issue id is a pr id in the db
	if event.Issue != nil && event.Issue.Assignee != nil && event.Sender != nil && event.Sender.Login != nil {
		expireTime := expireTime(time.Now())
		validTime := true
		//Check if the right person assignee commented
		prs := []models.Pullrequest{}
		err := tx.RawQuery(`SELECT * FROM pullrequests WHERE issue_url=?`, event.Issue.URL).All(&prs)
		if err != nil || prs == nil {
			log.Print("Could not make query")
			return false
		}
		//Might not be the best
		issueURL := prs[0].IssueUrl
		if strings.EqualFold(*event.Issue.Assignee.Login, *event.Sender.Login) && strings.EqualFold(*event.Issue.URL, issueURL) {
			q := tx.RawQuery(`UPDATE pullrequests SET valid_time=?, expire_time=? WHERE issue_url=?`, validTime, expireTime, issueURL)
			err := q.Exec()
			if err != nil {
				log.Print(err)
				log.Printf("Unable to update event number %d in CheckCommented", *event.Issue.ID)
				return false
			}
			return true
		}

	}
	return false
}

//CheckAssigned Determines if the PullRequestEvent action was assigneed and if the assignement was valid
func checkAssigned(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "assigned") {
		expireTime := expireTime(time.Now())
		validTime := false
		if event.PullRequest != nil && event.PullRequest.Assignees != nil {
			validTime = true
		}
		err := UpsertPullRequestEntry(ctx, event, tx, validTime, expireTime)
		if err != nil {
			return false
		}
		//Add new assignee if it does not exist in the repo
		err = AddAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		err = AddPullrequestAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false
}

//CheckUnassigned Determines if the PullRequestEvent action was unassigned and updates the database accordingly
func checkUnassigned(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "unassigned") {
		//Update PR in DB to no longer accept messages until assigned
		//if the assignee is nil update the expire time to be invalid
		validTime := false
		expireTime := time.Time{}
		if event.PullRequest != nil && event.Number != nil {
			err := UpsertPullRequestEntry(ctx, event, tx, validTime, expireTime)
			if err != nil {
				log.Printf("Unable to update event number %d in checkUnassigned", *event.Number)
				return false
			}
		}
		return true
	}
	return false
}

//CheckAssigned Determines if the PullRequestEvent action was reviewed and if the review was valid
func checkReviewed(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "review_requested") {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			err := AddAssigneeToDB(ctx, event, tx)
			if err != nil {
				log.Println(err)
				return false
			}
			if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
				validTime := true
				expireTime := expireTime(time.Now())
				if event.PullRequest != nil && event.Number != nil {
					err := UpsertPullRequestEntry(ctx, event, tx, validTime, expireTime)
					if err != nil {
						log.Printf("Unable to update event number %d in checkReviewed", *event.Number)
						return false
					}
				}
				err = AddPullrequestAssigneeToDB(ctx, event, tx)
				if err != nil {
					log.Println(err)
					return false
				}
				return true
			}
		}
	}
	return false
}

//CheckLabeled Determines if the PullRequestEvent action was labeled and if the label action was valid
func checkLabeled(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "labeled") {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			err := AddAssigneeToDB(ctx, event, tx)
			if err != nil {
				log.Println(err)
				return false
			}
			//Check that the label action was done by the assignee
			if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
				validTime := true
				expireTime := expireTime(time.Now())
				if event.PullRequest != nil && event.Number != nil {
					err := UpsertPullRequestEntry(ctx, event, tx, validTime, expireTime)
					if err != nil {
						log.Printf("Unable to update event number %d in checkLabeled", *event.Number)
						return false
					}
				}
				err = AddPullrequestAssigneeToDB(ctx, event, tx)
				if err != nil {
					log.Println(err)
					return false
				}
				return true
			}
		}
	}
	return false
}

//CheckClosed Determines if the PullRequestEvent action was closed and marks the PR as invalid
func checkClosed(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "closed") {
		expireTime := time.Time{}
		validTime := false
		if event.PullRequest != nil && event.Number != nil {
			err := UpsertPullRequestEntry(ctx, event, tx, validTime, expireTime)
			if err != nil {
				log.Printf("Unable to update event number %d in checkClosed", *event.Number)
				return false
			}
		}
		err := AddAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		err = AddPullrequestAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false
}

//CheckOpened Determines if the PullRequestEvent action was open and updated the entry in the DB
func checkOpened(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "opened") {
		expiration := time.Time{}
		validTime := false
		if event.PullRequest.Assignee != nil {
			expiration = expireTime(time.Now())
			validTime = true
		}
		err := UpsertPullRequestEntry(ctx, event, tx, validTime, expiration)
		if err != nil {
			log.Printf("Unable to update event number %d in checkOpened", *event.Number)
			return false
		}
		err = AddAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		err = AddPullrequestAssigneeToDB(ctx, event, tx)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return false

}

//CheckAssigned Determines if the PullRequestEvent action was reviewed and if the review was valid
func checkEdited(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "edited") {
		log.Print("made it to editied event")
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//Make sure the sender is the same as the assignee
			//if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
			err := AddAssigneeToDB(ctx, event, tx)
			if err != nil {
				log.Println(err)
				return false
			}
			expireTime := expireTime(time.Now())
			if event.PullRequest != nil && event.Number != nil {
				err := UpsertPullRequestEntry(ctx, event, tx, true, expireTime)
				log.Print("updated event")
				if err != nil {
					log.Printf("Unable to update event number %d in checkEdited", *event.Number)
					return false
				}
			}
			err = AddPullrequestAssigneeToDB(ctx, event, tx)
			if err != nil {
				log.Println(err)
				return false
			}
			return true
		}
	}
	return false
}

//UpsertPullRequestEntry performs an upsert on the pullrequests table to update or add a new pullrequest entry
func UpsertPullRequestEntry(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection, valid_time bool, expire_time time.Time) error {
	id, err := uuid.NewV1()
	if err != nil {
		return err
	}
	q := tx.RawQuery(`INSERT INTO pullrequests (id, created_at, updated_at, git_prid, url, html_url, issue_url, number, state,
			valid_time, title, body, request_created_at, request_updated_at, request_merged_at,
			request_closed_at, commits_url, status_url, expire_time)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
			ON CONFLICT (git_prid) DO UPDATE SET valid_time=?, expire_time=?`,
		id, time.Now(), time.Now(), *event.PullRequest.ID, *event.PullRequest.URL,
		*event.PullRequest.HTMLURL, *event.PullRequest.IssueURL, *event.PullRequest.Number,
		*event.PullRequest.State, valid_time, *event.PullRequest.Title, "",
		NullCheckTime(event.PullRequest.CreatedAt), NullCheckTime(event.PullRequest.UpdatedAt),
		NullCheckTime(event.PullRequest.MergedAt), NullCheckTime(event.PullRequest.ClosedAt),
		NullCheckInt(event.PullRequest.Commits), *event.PullRequest.StatusesURL, expire_time, valid_time, expire_time)
	err = q.Exec()
	if err != nil {
		return fmt.Errorf("Could not complete upsert: %v", err)
	}
	return nil
}

//AddAssigneeToDB adds the assigneeo of the pull request to the DB
func AddAssigneeToDB(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) error {
	if event.PullRequest != nil && event.PullRequest.Assignees != nil {
		assignees := event.PullRequest.Assignees
		for _, assignee := range assignees {
			assigneeID, err := uuid.NewV1()
			if err != nil {
				return err
			}
			q := tx.RawQuery(`INSERT INTO assignees (id, created_at, updated_at, login, type, html_url)
			VALUES (?, ?, ?, ?, ?, ?) ON CONFLICT (login) DO NOTHING`,
				assigneeID, time.Now(), time.Now(), assignee.Login, assignee.Type, assignee.HTMLURL)
			exErr := q.Exec()
			if exErr != nil {
				log.Printf("Unable to update event number %d in checkAssigned", *event.Number)
				return exErr
			}
		}
		return nil
	}
	return fmt.Errorf("Assignee is nil. Cannot add nil assignee to DB")
}

//AddPullrequestAssigneeToDB creates an entry in the pullrequest_assignee table
func AddPullrequestAssigneeToDB(ctx context.Context, event github.PullRequestEvent, tx *pop.Connection) error {
	if event.PullRequest != nil && event.PullRequest.Assignee != nil {
		//Get assignee ID
		assignees := []models.Assignee{}
		err := models.DB.RawQuery(`SELECT * FROM assignees WHERE login=?`, event.PullRequest.Assignee.Login).All(&assignees)
		if err != nil {
			return err
		}
		if assignees == nil {
			log.Print("Assignee is not in the database")
			return nil
		}
		assigneeID := assignees[0].ID

		//Get pullrequest ID
		prs := []models.Pullrequest{}
		err = models.DB.RawQuery(`SELECT * FROM pullrequests WHERE git_prid=?`, event.PullRequest.ID).All(&prs)
		if err != nil {
			return err
		}
		if prs == nil {
			log.Print("The pull request is not in the database")
			return nil
		}
		pullrequestID := prs[0].ID

		pullrequestAssigneeID, err := uuid.NewV1()
		if err != nil {
			return err
		}
		q := tx.RawQuery(`INSERT INTO pullrequest_assignees (id, created_at, updated_at, pullrequest_id, assignee_id)
			VALUES (?, ?, ?, ?, ?)`,
			pullrequestAssigneeID, time.Now(), time.Now(), pullrequestID, assigneeID)
		exErr := q.Exec()
		if exErr != nil {
			log.Printf("Unable to update event number %d in checkAssigned", *event.Number)
			return exErr
		}
	}
	return nil
}

//NullCheckTime determines if the time.Time value is nil
func NullCheckTime(x *time.Time) time.Time {
	if x == nil {
		return time.Time{}
	}
	return *x
}

//NullCheckInt checks if the int point is nil and return 0 if it is nil
func NullCheckInt(x *int) int {
	if x == nil {
		return 0
	}
	return *x
}
