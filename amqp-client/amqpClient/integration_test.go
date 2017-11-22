package amqpClient

import (
	"testing"
	"time"
)

func TestPkgIntegration(t *testing.T) {
	// <setup code>
	file := newAppConfigSetup()
	defer newAppConfigTearDown(file)

	err := Initialize(file.Name(), testQueues)
	if err != nil {
		t.Fail()
	}

	// creating a connection
	cl, err := New()
	defer cl.Disconnect()

	// we will publish this event 20 times and see f we can consume 20 messages in return
	for i := 1; i <= 20; i++ {
		testRandomEvenet := Event{
			ID:   "1",
			Name: "new_order.placed", // consumer should send notification to verify the new email address
			Date: time.Now(),
		}
		cl.Publish(&testRandomEvenet, testQueues.QueueUserAccount)
	}

	// consuming the events with 5 concurrent workers
	cl.Consume(testQueues.QueueUserAccount, 5, func(msg []byte, consumer string) {
		consumed++
	})

	time.Sleep(time.Second * 2)

	if consumed != 20 {
		t.Fail()
	}

}

var consumed = 0

type TestQueues struct {
	QueueUserAccount string
}

var testQueues = TestQueues{
	QueueUserAccount: "test_queue", // the actual name of the queue to be created by amqp
}
