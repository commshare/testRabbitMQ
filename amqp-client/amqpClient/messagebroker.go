package amqpClient

//MessageBrokerClientInterface : to create a message broker client, the client must implement this interface
//so that it can be used by other system components
type MessageBrokerClientInterface interface {
	Connect() error
	Disconnect() error
	CreateQueue(queueName string) error
	Publish(event SerializableEventInterface, queueName string) error
	Consume(queueName string, workers uint, callback func(msg []byte, consumer string)) error
	WatchWorkersStream()
}
