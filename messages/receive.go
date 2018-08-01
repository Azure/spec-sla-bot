package messages

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-service-bus-go"
	"github.com/Azure/spec-sla-bot/models"
)

type Message struct {
	PRID           string
	PullRequestURL string
	Assignee       string
}

func ReceiveFromQueue(ctx context.Context, connStr string) (*servicebus.ListenerHandle, error) {
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	log.Print("new namespace created")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	queueName := "24hrgitevents"
	q, err := getQueueToReceive(ns, queueName)
	if err != nil {
		log.Printf("failed to build a new queue named %q\n", queueName)
		return nil, err
	}
	log.Print("got queue to receive")
	listenHandle, err := q.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
		messageStruct, err := parseMessage(message.Data)
		log.Print("parsed message")
		if err != nil {
			log.Println(err)
			return message.DeadLetter(err)
		}
		if ShouldSend(messageStruct) {
			err = SendEmailToAssignee(ctx, messageStruct)
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
	q, err := ns.NewQueue(ctx, queueName)
	return q, err
}

//Redo this
func parseMessage(data []byte) (*Message, error) {
	str := string(data[:])
	if len(str) != 0 {
		strSplit := strings.FieldsFunc(str, Split)
		for i, v := range strSplit {
			strSplit[i] = strings.TrimSpace(v)
		}
		return &Message{PRID: strSplit[1], PullRequestURL: strSplit[3], Assignee: strSplit[5]}, nil
	}
	return nil, errors.New("could not parse messages returned by service bus")
}

func Split(r rune) bool {
	return r == ','
}

func ShouldSend(messageStruct *Message) bool {
	gitPRID, _ := strconv.Atoi(messageStruct.PRID)
	prs := []models.Pullrequest{}
	err := models.DB.Where("git_prid=?", int64(gitPRID)).All(&prs)
	if err != nil || prs == nil {
		log.Print("Could not make querey")
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
