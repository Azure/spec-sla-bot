package messages

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/spec-sla-bot/models"
	"github.com/gobuffalo/pop"
	"github.com/google/go-github/github"
)

//type AcknowledgmentStatus unint32
//var tx *pop.Connection
var currentExpireTime time.Time

func CheckAcknowledgement(event github.PullRequestEvent) {
	//Check if PR is in the database
	//Add if not (now or in another function)
	log.Print("CONNECT TO DEVELOPEMENT DB")
	tx, err := pop.Connect("developement")
	if err != nil {
		log.Printf("Could not conntect to the developement database")
		return
	}
	log.Print("MADE IT HERE")
	if checkClosed(event, tx) || checkUnassigned(event, tx) || (event.PullRequest.Assignee == nil && checkOpened(event, tx)) {
		//update event in DB to show the PR is no longer open and no more messages will be accepted for that PR ID
		//don't send a message
	} else if event.PullRequest.Assignee != nil && (checkAssigned(event, tx) || checkReviewed(event, tx) || checkEdited(event, tx) || checkLabeled(event, tx) || checkOpened(event, tx)) {
		//send a message with PR id
		//Format string with PR ID
		log.Printf("Close: %s, Unassigned: %s, Opened: %s, Assigned: %s, Reviewed: %s, Edited: %s, Labeled: %s", checkClosed(event, tx), checkUnassigned(event, tx), checkOpened(event, tx),
			checkAssigned(event, tx), checkReviewed(event, tx), checkEdited(event, tx), checkLabeled(event, tx))
		log.Print("MADE IT HERE")
		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.PullRequest.Number, *event.PullRequest.URL, *event.PullRequest.Assignee.Login)
		log.Print(message)
		err = SendToQueue(message, currentExpireTime)
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
		log.Printf("Could not conntect to the developement database")
		return
	}
	if event.Issue.IsPullRequest() && checkCommented(event, tx) {
		message := fmt.Sprintf("PR id, %d, URL, %s, Assignee, %s", *event.Issue.ID, *event.Issue.URL, *event.Issue.Assignee.Login)
		log.Print(message)
		err = SendToQueue(message, currentExpireTime)
		log.Print("SENT TO QUEUE")
		if err != nil {
			log.Printf("Message for event %d not delivered", *event.Issue.ID)
		}
		return
	}
	log.Printf("Comment event was not on a pull request issue")
}

func updateTime() time.Time {
	currentTime := time.Now().Local()
	if strings.Compare(currentTime.Weekday().String(), "Friday") == 0 {
		currentExpireTime := currentTime.Add(time.Hour * time.Duration(48))
		return currentExpireTime
	}
	currentExpireTime := currentTime.Add(time.Hour * time.Duration(24))
	return currentExpireTime
}

func checkCommented(event github.IssueCommentEvent, tx *pop.Connection) bool {
	//check that the issue is not nil and that the issue id is a pr id in the db
	if event.Issue != nil && event.Issue.Assignee != nil {
		//check in the DB that the repo is not
		//query := tx.Where("issue_url = ?", event.Comment.IssueURL)
		/*expireTime := updateTime()
		time := models.ValidTime{
			Time:  expireTime,
			Valid: true,
		}
		if event.Issue.PullRequestLinks!= nil && event.PullRequest.Assignee == nil && event.Number != nil {
			err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", time, *event.Number)
			if err != nil {
				log.Print("Unable to update event number %d", *event.Number)
				return false
			}
		}*/
		return true
	}
	return false
}

func checkAssigned(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "assigned") == 0 {
		expireTime := updateTime()
		err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", expireTime, event.Number)
		if err != nil {
			log.Print("Unable to update event number %d", *event.Number)
			return false
		}
		if event.PullRequest != nil && event.PullRequest.Assignees != nil {
			//if *event.PullRequest != nil && *event.PullRequest.Assignees != nil {
			assignees := event.PullRequest.Assignees
			for _, assignee := range assignees {
				err = tx.RawQuery("INSERT INTO assignees (login, type, html_url) SELECT ?,?,? WHERE NOT EXISTS (SELECT ? FROM assignees WHERE login = ?)", assignee.Login, assignee.Type, assignee.HTMLURL, assignee.Login)
				if err != nil {
					log.Print("Unable to add assignee %s", *event.PullRequest.Assignee)
					return false
				}
			}
			//}

		}
		return true
	}
	return false
}

func checkUnassigned(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "unassigned") == 0 {
		//Update PR in DB to no longer accept messages until assigned
		//if the assignee is nil update the expire time to be invalid
		expireTime := updateTime()
		if event.PullRequest != nil && event.PullRequest.Assignee == nil && event.Number != nil {
			err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", expireTime, *event.Number)
			if err != nil {
				log.Print("Unable to update event number %d", *event.Number)
				return false
			}
		}
		return true
	}
	return false
}

func checkReviewed(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "review_requested") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.Name) == 0 {
				expireTime := updateTime()
				if event.PullRequest != nil && event.PullRequest.Assignee == nil && event.Number != nil {
					err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", expireTime, *event.Number)
					if err != nil {
						log.Print("Unable to update event number %d", *event.Number)
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
	if strings.Compare(*event.Action, "labeled") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.Name) == 0 {
			//Update DB
			expireTime := updateTime()
			if event.PullRequest != nil && event.PullRequest.Assignee == nil && event.Number != nil {
				err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", expireTime, *event.Number)
				if err != nil {
					log.Print("Unable to update event number %d", *event.Number)
					return false
				}
			}
			return true
			//}
		}
	}
	return false
}

func checkClosed(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "closed") == 0 {
		//Update DB to not accept messages
		//expireTime := updateTime()
		if event.PullRequest != nil && event.PullRequest.Assignee == nil && event.Number != nil {
			err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", time.Time{}, *event.Number)
			if err != nil {
				log.Print("Unable to update event number %d", *event.Number)
				return false
			}
		}
	}
	return false
}

func checkOpened(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "opened") == 0 {
		//expireTime := updateTime()
		//check if PR ID is in the DB
		//Create new entry with assignee if not
		if event.PullRequest.Assignee == nil {
			//PR id cannot accept messages (not assigned yet)
			pr := &models.Pullrequest{
				GitPRID:          *event.PullRequest.ID,
				URL:              *event.PullRequest.URL,
				HtmlUrl:          *event.PullRequest.HTMLURL,
				IssueUrl:         *event.PullRequest.IssueURL,
				Number:           *event.PullRequest.Number,
				State:            *event.PullRequest.State,
				ValidTime:        false,
				Title:            *event.PullRequest.Title,
				Body:             *event.PullRequest.Body,
				RequestCreatedAt: *event.PullRequest.CreatedAt,
				RequestUpdatedAt: NullCheckTime(event.PullRequest.UpdatedAt),
				RequestMergedAt:  NullCheckTime(event.PullRequest.MergedAt),
				RequestClosedAt:  NullCheckTime(event.PullRequest.ClosedAt),
				CommitsUrl:       NullCheckInt(event.PullRequest.Commits), // may need a null check to get the CommitsURL
				StatusUrl:        *event.PullRequest.StatusesURL,          // consider changing name of column to match statuses
				ExpireTime:       time.Time{},
			}
			err := models.DB.Create(pr)
			if err != nil {
				log.Printf("Could not create pr entry for checkopened, pr number %d", *event.Number)
				return false
			}
		} else {
			//time.Valid = true

			//update DB for PR to accept messages
			//Case when PR closed with an assignee but then reopenned

			//Not sure what to do here

		}
		return true
	}
	return false
}

func checkEdited(event github.PullRequestEvent, tx *pop.Connection) bool {
	if strings.Compare(*event.Action, "edited") == 0 {
		if event.PullRequest.Assignee != nil && event.Sender != nil {
			//if strings.Compare(*event.PullRequest.Assignee.Name, *event.Sender.) == 0 {
			//Update DB to accept messages
			expireTime := updateTime()
			if event.PullRequest != nil && event.PullRequest.Assignee == nil && event.Number != nil {
				err := tx.RawQuery("UPDATE pullrequests SET expire_time=? WHERE number=?", expireTime, *event.Number)
				if err != nil {
					log.Print("Unable to update event number %d", *event.Number)
					return false
				}
			}
			return true
		}
	}
	return false
}

func NullCheckTime(x *time.Time) time.Time {
	if x == nil {
		return time.Time{}
	}
	return *x
}

func NullCheckInt(x *int) int {
	if x == nil {
		return 0
	}
	return *x
}
