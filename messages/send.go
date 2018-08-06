package messages

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func SendToQueue(message string, postTime time.Time) error {
	connStr := os.Getenv("CUSTOMCONNSTR_SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		return err
	}

	queueName := "24hrgitevents"
	q, err := getQueueToSend(ns, queueName)

	if err != nil {
		log.Printf("failed to build a new queue named %q\n", queueName)
		return err
	}
	postTime = postTime.Add(time.Minute * 10)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	msg := servicebus.NewMessageFromString(message)
	msg.SystemProperties = &servicebus.SystemProperties{
		ScheduledEnqueueTime: &postTime,
	}
	log.Printf("ABOUT TO SEND MESSAGE")
	log.Print(message)
	q.Send(ctx, msg)
	cancel()
	return nil
}

func getQueueToSend(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q, err := ns.NewQueue(ctx, queueName)
	return q, err
}
