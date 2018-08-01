package messages

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/spec-sla-bot/models"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/google/go-github/github"
)

var currentExpireTime time.Time

func CheckAcknowledgement(event github.PullRequestEvent) {
	if checkClosed(event, models.DB) || checkUnassigned(event, models.DB) || (event.PullRequest.Assignee == nil && checkOpened(event, models.DB)) {
		err := UpsertPullRequestEntry(event, models.DB, false, time.Time{})
		if err != nil {
			log.Printf("Unable to update event number %d", *event.Number)
		}
	} else if event.PullRequest.Assignee != nil && (checkAssigned(event, models.DB) || checkReviewed(event, models.DB) || checkEdited(event, models.DB) || checkLabeled(event, models.DB) || checkOpened(event, models.DB)) {
		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.PullRequest.ID, *event.PullRequest.HTMLURL, *event.PullRequest.Assignee.Login)
		log.Print(message)
		err := SendToQueue(message, currentExpireTime)
		log.Print("SENT TO QUEUE")
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.PullRequest.ID)
		}
	}
}

func CheckAcknowledgementComment(event github.IssueCommentEvent) {
	log.Print("CONNECT TO DEVELOPEMENT DB")
	tx, err := pop.Connect("developement")
	if err != nil {
		log.Printf("Could not connect to the developement database")
		return
	}
	if event.Issue != nil && event.Issue.IsPullRequest() && checkCommented(event, tx) {

		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.Issue.ID, *event.Issue.URL, *event.Issue.Assignee.Login)
		log.Print(message)
		err = SendToQueue(message, currentExpireTime)
		log.Print("SENT TO QUEUE")
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.Issue.ID)
		}
		return
	}
}

func CheckAcknowledgementLabel(event github.LabelEvent) {
	log.Printf("Logic to handle a label event not implemented")
}

func updateTime() time.Time {
	currentTime := time.Now().Local()
	//Adjusted time for now for testing purposes
	//if the current weekday is Friday
	if currentTime.Weekday() == 5 {
		currentExpireTime := currentTime.Add(time.Hour * time.Duration(72))
		return currentExpireTime
	}
	currentExpireTime := currentTime.Add(time.Minute * time.Duration(5))
	return currentExpireTime
}

func checkCommented(event github.IssueCommentEvent, tx *pop.Connection) bool {
	//check that the issue is not nil and that the issue id is a pr id in the db
	if event.Issue != nil && event.Issue.Assignee != nil && event.Sender != nil && event.Sender.Login != nil {
		expireTime := updateTime()
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
				log.Printf("Unable to update event number %d", *event.Issue.ID)
				return false
			}
			return true
		}

	}
	return false
}

func checkAssigned(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "assigned") {
		expireTime := updateTime()
		validTime := false
		if event.PullRequest != nil && event.PullRequest.Assignees != nil {
			validTime = true
		}
		err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
		if err != nil {
			return false
		}
		//Add new assignee if it does not exist in the repo
		if event.PullRequest != nil && event.PullRequest.Assignees != nil {
			assignees := event.PullRequest.Assignees
			for _, assignee := range assignees {
				err := tx.RawQuery(`INSERT INTO pullrequests (login, type, html_url)
				VALUES (?, ?, ?) ON CONFLICT CONSTRAINT (login) DO UPDATE SET html_url = ?`,
					assignee.Login, assignee.Type, assignee.HTMLURL, assignee.HTMLURL)
				if err != nil {
					log.Printf("Unable to update event number %d", *event.Number)
					return false
				}
			}
			return true
		}
	}
	return false
}

func checkUnassigned(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "unassigned") {
		//Update PR in DB to no longer accept messages until assigned
		//if the assignee is nil update the expire time to be invalid
		validTime := false
		expireTime := time.Time{}
		if event.PullRequest != nil && event.Number != nil {
			err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
			if err != nil {
				log.Printf("Unable to update event number %d", *event.Number)
				return false
			}
		}
		return true
	}
	return false
}

func checkReviewed(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "review_requested") {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
				validTime := true
				expireTime := updateTime()
				if event.PullRequest != nil && event.Number != nil {
					err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
					if err != nil {
						log.Printf("Unable to update event number %d", *event.Number)
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func checkLabeled(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "labeled") {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//Check that the label action was done by the assignee
			if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
				validTime := true
				expireTime := updateTime()
				if event.PullRequest != nil && event.Number != nil {
					err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
					if err != nil {
						log.Printf("Unable to update event number %d", *event.Number)
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

func checkClosed(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "closed") {
		expireTime := time.Time{}
		validTime := false
		if event.PullRequest != nil && event.Number != nil {
			err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
			if err != nil {
				log.Printf("Unable to update event number %d", *event.Number)
				return false
			}
		}
	}
	return false
}

func checkOpened(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "opened") {
		expireTime := time.Time{}
		validTime := false
		if event.PullRequest.Assignee != nil {
			expireTime = updateTime()
			validTime = true
		}
		err := UpsertPullRequestEntry(event, tx, validTime, expireTime)
		if err != nil {
			log.Printf("Unable to update event number %d", *event.Number)
			return false
		}
		return true
	}
	return false

}

func checkEdited(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.EqualFold(*event.Action, "edited") {
		log.Print("made it to editied event")
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//Make sure the sender is the same as the assignee
			//if strings.EqualFold(*event.PullRequest.Assignee.Name, *event.Sender.Name) {
			expireTime := updateTime()
			if event.PullRequest != nil && event.Number != nil {
				err := UpsertPullRequestEntry(event, tx, true, expireTime)
				log.Print("updated event")
				if err != nil {
					log.Printf("Unable to update event number %d", *event.Number)
					return false
				}
			}
			return true
			//}
		}
	}
	return false
}

func UpsertPullRequestEntry(event github.PullRequestEvent, tx *pop.Connection, valid_time bool, expire_time time.Time) error {
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
		log.Print(err)
		log.Printf("Unable to update event number %d", *event.Number)
		return errors.New("Could not complete upsert")
	}
	return nil
}

//NullCheckTime
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
