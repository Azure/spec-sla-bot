package messages

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Azure/azure-service-bus-go"
	"github.com/Azure/spec-sla-bot/models"
	"github.com/gobuffalo/uuid"
)

//ReceiveFromQueue Sets up the queue to receive from service bus and determines whether
//or not an email should be sent to the assignees and updates the datebase accordingly
func ReceiveFromQueue(ctx context.Context, connStr string) (*servicebus.ListenerHandle, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		return nil, err
	}

	queueName := "24hrgitevents"
	q, err := getQueueToReceive(ns, queueName)
	if err != nil {
		log.Printf("failed to build a new queue named %q\n", queueName)
		return nil, err
	}

	listenHandle, err := q.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
		messageStruct := &MessageContent{}
		err := json.Unmarshal(message.Data, messageStruct)
		log.Print("parsed message")
		if err != nil {
			log.Println(err)
			return message.DeadLetter(err)
		}
		if ShouldSendAssigneeEmail(messageStruct) {
			err = SendEmailToAssignee(ctx, messageStruct)
			if err != nil {
				log.Println(err)
				return message.DeadLetter(err)
			}
			err = AddEmailToDB(messageStruct)
			if err != nil {
				log.Println("Unable to add the emails to the database")
				return message.DeadLetter(err)
			}
		} else if messageStruct.ManagerEmailReminder {
			err = SendEmailToManager(ctx, messageStruct)
			if err != nil {
				log.Println(err)
				return message.DeadLetter(err)
			}
		}
		return message.Complete()
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return listenHandle, nil
}

//getQueueToReveive gets the queue to receive messages from service bus
func getQueueToReceive(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil {
		return nil, err
	}
	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			return nil, err
		}
	}
	q, err := ns.NewQueue(queueName)
	return q, err
}

//ShouldSendAssigneeEmail determines if an email should be sent to the assignee based on
//whether the time is valid and the current time is greater than the expire
//time
func ShouldSendAssigneeEmail(messageStruct *MessageContent) bool {
	prs := []models.Pullrequest{}
	err := models.DB.Where("git_prid=?", messageStruct.PRID).All(&prs)
	if err != nil || prs == nil {
		return false
	}
	for _, pr := range prs {
		if pr.ValidTime && time.Now().Sub(pr.ExpireTime) >= 0 {
			log.Print("returning true, should send message")
			return true
		}
	}
	return false
}

//AddEmailToDB adds an email entry to the database for use when sending
//the manager email
func AddEmailToDB(messageStruct *MessageContent) error {
	emailID, err := uuid.NewV1()
	if err != nil {
		return err
	}
	q := models.DB.RawQuery(`INSERT INTO emails (id, created_at, updated_at, pullrequest_id, time_sent)
			VALUES (?, ?, ?, ?, ?)`,
		emailID, time.Now(), time.Now(), int(messageStruct.PRID), time.Now())
	err = q.Exec()
	if err != nil {
		log.Print(err)
		return errors.New("Could not complete insert to add email to database")
	}

	assignees := []models.Assignee{}
	err = models.DB.RawQuery(`SELECT * FROM assignees WHERE login=?`, messageStruct.AssigneeLogin).All(&assignees)
	if err != nil {
		return err
	}
	if assignees == nil {
		log.Print("Assignee is not in the database")
		return nil
	}

	assigneeID := assignees[0].ID
	emailAssigneeID, err := uuid.NewV1()
	if err != nil {
		return err
	}
	q = models.DB.RawQuery(`INSERT INTO email_assignees (id, email_id, assignee_id, created_at, 
		updated_at) VALUES (?, ?, ?, ?, ?)`,
		emailAssigneeID, emailID, assigneeID, time.Now(), time.Now())
	err = q.Exec()
	if err != nil {
		log.Print(err)
		return errors.New("Could not complete insert to add email to database")
	}
	return nil
}
