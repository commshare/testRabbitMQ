package amqpClient

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event : is the simplest form of event, provided as a quick way to
type Event struct {
	ID   string
	Name string
	Date time.Time
}

//Serialize implements how an event of type Event is serialize
func (e *Event) Serialize() ([]byte, error) {
	serialized, err := json.Marshal(e)
	if err != nil {
		fmt.Println(err)
		return []byte(""), err
	}
	return serialized, nil
}

//Deserialize implements how an event of type Event is deserialize
func (e *Event) Deserialize(jsonMessage []byte) error {
	err := json.Unmarshal(jsonMessage, &e)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
