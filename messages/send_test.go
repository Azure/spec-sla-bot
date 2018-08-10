package messages

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

func Test_SendToQueue(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	connStr := os.Getenv("CUSTOMCONNSTR_SERVICEBUS_CONNECTION_STRING")
	if connStr == "" {
		t.Skip()
		return
	}

	want := RandomString(15)
	tempQueueName := RandomString(10)
	if testing.Verbose() {
		t.Logf("queue name: %q", tempQueueName)
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		t.Error(err)
		return
	}

	queueManager := ns.NewQueueManager()

	q, err := getQueueToSend(ns, tempQueueName)
	if err != nil {
		t.Error(err)
		return
	}

	const waitTime = time.Duration(2 * time.Minute)
	fmt.Println("Wait time: ", waitTime)
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second+waitTime)
	defer cancel()

	_, err = queueManager.Put(ctx, tempQueueName)
	if err != nil {
		t.Error(err)
		return
	}
	defer queueManager.Delete(ctx, tempQueueName)

	SendToQueue(ctx, want, time.Now().Add(waitTime), tempQueueName)

	start := time.Now()
	message, err := q.ReceiveOne(ctx)
	defer message.Complete()

	timeWaited := time.Now().Sub(start)

	if err != nil {
		t.Error(err)
		return
	}

	got := string(message.Data)
	if got != want {
		t.Logf("got:\n\t%q\nwant:\n\t%q", got, want)
		t.Fail()
	} else {
		t.Logf("message matched")
	}

	const buffer = time.Duration(30 * time.Second)
	if min := waitTime - buffer; timeWaited < min {
		t.Logf("received message after %v, expected to wait at least %v", timeWaited, min)
		t.Fail()
	} else if max := waitTime + buffer; timeWaited > max {
		t.Logf("received message after %v, expected to wait no longer than %v", timeWaited, max)
		t.Fail()
	}
}

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
