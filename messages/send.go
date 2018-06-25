package messages

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func SendToQueue(message string) error {
	//connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING")
	connStr := "Endpoint=sb://spec-sla-bus.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=YetSU0WSSf0Bnb+Y9wndzDdXP2DKGoH70GiNAJNl9tk="
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
