package messages

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-service-bus-go"
)

func SendToQueue(ctx context.Context, message []byte, postTime time.Time, queueName string) error {
	connStr := os.Getenv("CUSTOMCONNSTR_SERVICEBUS_CONNECTION_STRING")
	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		return err
	}

	fmt.Println("Post Time: ", postTime, " Now: ", time.Now())

	q, err := getQueueToSend(ns, queueName)
	fmt.Println("Queue Name: ", queueName)
	if err != nil {
		log.Printf("failed to build a new queue named %q\n", queueName)
		return err
	}

	msg := servicebus.NewMessage(message)
	msg.SystemProperties = &servicebus.SystemProperties{
		ScheduledEnqueueTime: &postTime,
	}

	fmt.Println("Scheduled Enqueue Time: ", msg.SystemProperties.ScheduledEnqueueTime)

	err = q.Send(ctx, msg)
	if err != nil {
		return err
	}
	q.Close(ctx)

	return nil
}

func getQueueToSend(ns *servicebus.Namespace, queueName string) (*servicebus.Queue, error) {
	q, err := ns.NewQueue(queueName)
	return q, err
}
