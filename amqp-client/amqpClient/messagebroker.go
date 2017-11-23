package amqpClient
// Item struct is used to send down a channel
// the map received from a delivery and the callbacks
// for acknowledge(Ack) or negatively acknowledge(Nack)
type Element struct {
	Body  []byte
	ack  func(bool) error
	nack func(bool, bool) error
}

//MessageBrokerClientInterface : to create a message broker client, the client must implement this interface
//so that it can be used by other system components
type MessageBrokerClientInterface interface {
	Connect() error
	Disconnect() error
	CreateQueue(queueName string) error
	Publish(event SerializableEventInterface, queueName string) error
	Consume(queueName string, workers uint, callback func(msg []byte, consumer string)) error
	WatchWorkersStream()
	Consume2(queueName string, workers uint, callback func(items chan Element, consumer string)) error
}
