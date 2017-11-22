# amqp-client

[![GoDoc](https://godoc.org/github.com/almasry/amqp-client?status.svg)](https://godoc.org/github.com/almasry/amqp-client) [![Go Report Card](https://goreportcard.com/badge/github.com/almasry/amqp-client)](https://goreportcard.com/report/github.com/almasry/amqp-client) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


A rabbitmq client-wrapper library that uses  [streadway/amqp](http://github.com/streadway/amqp).


## Installation

    go get github.com/almasry/amqp-client

## Documentation

Available on [godoc.org](http://www.godoc.org/github.com/almasry/amqp-client), also reading this 
document would help you quickly get started using the package, further examples are available at [amqp-client/examples](https://github.com/almasry/amqp-client/blob/master/examples) .

### Pkg design    

The pkg uses a good level of abstraction to give the user a great deal of flexibility through decoupling the api from
 the actual implementation. 
 
* The [streadway/amqp](http://github.com/streadway/amqp) can be replaced with another library if 
 needed, as the Client api is depending on the general type [MessageBrokerClientInterface](https://github.com/almasry/amqp-client/blob/master/messagebroker.go#L5)
 and any client can implement this interface to replace the current [RabbitMqClient](https://github.com/almasry/amqp-client/blob/master/rabbitmq_client.go#L16) implementation.


* The Event system is also very flexible, you can implement any type of event. You can ditch the amqpClient.Event and 
create a much more sophisticated one. The new event must be of type [SerializableEventInterface](https://github.com/almasry/amqp-client/blob/master/serializable_event.go#L5) and
must implement [Serialize](https://github.com/almasry/amqp-client/blob/master/event.go#L17) and [Deserialize](https://github.com/almasry/amqp-client/blob/master/event.go#L27) 
methods to tell the broker how to marshale / unmarshale the event while transporting it. You may want to have a look at the [SerializableEventInterface documentation](https://godoc.org/github.com/almasry/amqp-client#SerializableEventInterface) .

--------------
### Getting started   
Create a Struct that holds your queues as fields (probably in a separate file), the struct could have any name, 
but all its fields must be of type string, the fields might also start with 'Queue' prefix - but it's not mandatory.

```go
type Queues struct {
	QueueUserAccount    string
	QueueOrderWorkflow  string
	QueueServiceChanges string
	// ... 
	// ... all other queues should go here ..
}

```

* In your application init function (or the function you use to bootstrap), create an instance
 of the queues struct you just created, the values of the fields are the names of the queues as hey will appear in rabbitmq 

* Copy the [examples/config.yml](https://github.com/almasry/amqp-client/blob/master/examples/config.yml) file to your application
 directory or just copy the [rabbitmq ] configuration block to your current .yml file (if you use one). After editing the parameters, 
 make a call to the  amqpClient.Initialize( .. , ..)  with both file location and the new instance of queues struct as parameters to the function.
 
 
* The Initialize method will create all the queues provided, so you actually don't have to worry any more about creating 
queues in the run time. You can consume/ publish to any of these queues without having to worry about their creation process 
or which exchange they belong to.


You still can create queues dynamically during the run time by calling client method CreateQueue( name ) anywhere in your application.


```go
var AmqpQueues Queues

func init() {
    // setting values of queues
	AmqpQueues = Queues{
		QueueUserAccount:    "user_account",
		QueueOrderWorkflow:  "order_workflow",
		QueueServiceChanges: "service_changes",
	}
	// initializing the amqp client library
	amqpClient.Initialize("./config.yml", AmqpQueues)
}

```

Horraaay ! Now that you'e done with configuration, let's get to the fun part. 


--------------
### Example Publisher 

Publishing to any of the previously created queues is as simple as the following :

```go

	// creating a new client
	cl, err := amqpClient.New()
	failOnError(err, "Failed to create a new amqp client")
	defer cl.Disconnect()
	
	// creating an instance of the amqpClient.Event and publishing 
	userEmailChangedEvent := amqpClient.Event{
		ID:   "1",
		Name: "new_order.placed", // sending a notification to stat shipping the newly placed order ..
		Date: time.Now(),
	}
	cl.Publish(&userEmailChangedEvent, AmqpQueues.QueueOrderWorkflow)
```

* The default simple event  [amqpClient.Event](https://github.com/almasry/amqp-client/blob/master/event.go#L10)  is using json marshaling and json unmarshaling 
to handle the serialize /  deserialize process. why ? so in case if you have other consumers
(written in other programming languages ) they don' need to know about other types of serialization / deserialization, all they need is to decode a simple json message.


That' all about publishing events, very simple ! The full example is available at [examples/publisher](https://github.com/almasry/amqp-client/blob/master/examples/producer/main.go) .

-------------------
### Example Consumer 

Creating workers and consuming events is also very simple. Actually they are just one step, all you need to know is :
* the name of the queue you want to consume messages from 
* the number of concurrent workers you need (1 ... )
* the event handler or the callback function that will handle the event once it' received

```go

    // create a new client 
	client, err := amqpClient.New()
	failOnError(err)
	defer client.Disconnect()

	// QueuesList ==> is the global queues configuration we declared in the previous steps 

    // QueuesList.QueueServiceChanges  is th queue you wanna publish to 
    // 2  is th number of concurrent consumers (workers) that will be reading messages from  QueueServiceChanges 
	client.Consume(AmqpQueues.QueueServiceChanges, 2, func(msg []byte, consumer string) {
		eventHandler(msg, consumer)
	})
	
	// Do some other business logic here ...
	// ....
	// ....

	client.WatchWorkersStream()
```
* Workers are blocking by nature, but in this case Consume method runs asynchronously in a non-blocking way, at the end of your function 
you need to call client.WatchWorkersStream() to wait for the results of the consumers otherwise the consumers will exit.


* Yo can create as many consumers with any number of workers and run them simultanelously wihtout having to worry about how their 
concurrency works or fan out their results to on channel.


```go

    // create a new client 
	client, err := amqpClient.New()
	failOnError(err)
	defer client.Disconnect()

	client.Consume(AmqpQueues.QueueServiceChanges, 5, func(msg []byte, consumer string) {
		eventHandler(msg, consumer)
	})

	client.Consume(AmqpQueues.QueueUserAccount, 20, func(msg []byte, consumer string) {
		// 
		eventHandler(msg, consumer)
	})

	client.Consume(AmqpQueues.QueueOrderWorkflow, 9, func(msg []byte, consumer string) {
		eventHandler(msg, consumer)
	})

	// Do some other business logic here ...
	// ....

	client.WatchWorkersStream()
```


* In the call back function of the Consume method, you receive two parameters, the 'msg' in its row format, yo need to deserialize it, 
and the name of the consumer (the hash),  you may ignore the latest or use it for future debugging if needed.


```go

func eventHandler(msg []byte, consumer string) {
	// assuming you already know the type of event you'e receiving from the channel ..
	var event amqpClient.Event
	event.Deserialize(msg)

	log.Println(fmt.Sprintf("Consumer %s just processed event : %s", consumer, event.Name))
}
```



That' all about consuming events and creating workers. The full example is available at [examples/consumer](https://github.com/almasry/amqp-client/blob/master/examples/consumer/main.go) .

-------------------
### FAQ
 Why do you create queues while initializing the application ?

>  It's a good practice to create all the queues while bootstrapping the application, usually rabbitmq tries to allocate 
 disk space to redundantly save the queues and their messages to the disk, and f it fails for any reason - you'd like to 
 find out and fix it as early as possible, not during the run time.
 
> PS. the default configuration of the library sets all queues as 'durable', Durable queues will survive server restarts 
and remain even when there are no remaining consumers or bindings.  



 Why do I need to creating a struct for Queues names, why bother while you can just make them part of the yml configuration  ?

>  Using the queues names as fields of a struct makes it impossible for developers to miss-type the namesof the queues. While making the queues 
declaration part of the weakly typed yaml configuration will give more chance for error since
you have to go every time and look at the yml file and copy some string and paste it somewhere else, which is a very error-prone process. 

> A small bug of not publishing to the right queue can lead to losing thousands of events which is vry risky in an event-driven architecture. 


