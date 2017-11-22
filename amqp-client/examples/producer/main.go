package main

import (
	"github.com/almasry/amqp-client"
	"log"
	"time"
)

type Queues struct {
	QueueUserAccount    string
	QueueOrderWorkflow  string
	QueueServiceChanges string
}

var AmqpQueues Queues

func init() {

	AmqpQueues = Queues{
		QueueUserAccount:    "user_account",
		QueueOrderWorkflow:  "order_workflow",
		QueueServiceChanges: "service_changes",
	}
	err := amqpClient.Initialize("./examples/config.yml", AmqpQueues)
	failOnError(err)
}

func main() {
	// creating a connection
	cl, err := amqpClient.New()
	failOnError(err)
	defer cl.Disconnect()

	userEmailChangedEvent := amqpClient.Event{
		ID:   "1",
		Name: "new_order.placed", // consumer should send notification to verify the new email address
		Date: time.Now(),
	}
	cl.Publish(&userEmailChangedEvent, AmqpQueues.QueueOrderWorkflow)

	// publishing another event ..
	userPasswordChangedEvent := amqpClient.Event{
		ID:   "2",
		Name: "usr_email.changed",
		Date: time.Now(),
	}
	cl.Publish(&userPasswordChangedEvent, AmqpQueues.QueueUserAccount)
}

func failOnError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
