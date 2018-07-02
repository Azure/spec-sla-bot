package messages

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

type Message struct {
	PRID           string
	PullRequestURL string
	Assignee       string
}

func ReceiveFromQueue(ctx context.Context) (*servicebus.ListenerHandle, error) {
	connStr := mustGetenv("CUSTOMCONNSTR_SERVICEBUS_CONNECTION_STRING")
	//connStr :=
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

	log.Print("before receive")
	listenHandle, err := q.Receive(ctx, func(ctx context.Context, message *servicebus.Message) servicebus.DispositionAction {
		text := string(message.Data)
		log.Print(text)

		//The message is not invalid so parse
		log.Print("MADE IT TO RECEIVE")
		messageStruct, err := parseMessage(message.Data)
		log.Print("parsed message")
		if err != nil {
			log.Println(err)
			//os.Exit(1)
			return message.DeadLetter(err)
		}
		err = SendEmailToAssignee(messageStruct)
		if err != nil {
			log.Println(err)
			return nil
		}
		log.Print("sent email")
		return message.Complete()
	})

	//Not sure if this should stay
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("I am listening...")
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

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

func parseMessage(data []byte) (*Message, error) {
	str := string(data[:])
	log.Print(str)
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
