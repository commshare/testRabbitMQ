package amqpClient

//SerializableEventInterface : to create an event struct that can be published and consumed by the message broker
// the struct has to be of type SerializableEventInterface and should implement Serialize() and Deserialize() methods
type SerializableEventInterface interface {
	Serialize() ([]byte, error)
	Deserialize(jsonString []byte) error
}
