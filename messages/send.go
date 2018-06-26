package messages

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func SendToQueue(message string) error {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Print(message)
	q.Send(ctx, servicebus.NewMessageFromString(message))
	cancel()
	return nil
}

func getQueueToSend(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	q, err := ns.NewQueue(ctx, queueName)
	return q, err
}
