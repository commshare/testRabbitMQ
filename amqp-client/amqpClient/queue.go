package amqpClient

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var queuesNamesList = make([]string, 0)

var (
	errConfigVariableIsNotStruct = "The Queues type provided is not of a struct type"
	errFieldsTypeMustBeString    = "all fields of queue struct must be of type string "
	errNoQueueNameProvided       = "no name provided for queue : %s"
	errQueueNameNotUnique        = "two or more queues have the same name: %s "
)

//extractQueuesNames extracts and validates the provided queues struct and its values
func extractQueuesNames(queues interface{}) ([]string, error) {

	// checking type of the queues (it' supposed to be of type struct)
	if reflect.TypeOf(queues).Kind() != reflect.Struct {
		return nil, errors.New(errConfigVariableIsNotStruct)
	}

	for i := 0; i < reflect.TypeOf(queues).NumField(); i++ {

		// checking type of the queues names (it' supposed to be of type string)
		if reflect.TypeOf(queues).Field(i).Type.String() != reflect.String.String() {
			return nil, errors.New(errFieldsTypeMustBeString)
		}

		field := reflect.TypeOf(queues).Field(i).Name
		queueName := strings.TrimSpace(reflect.ValueOf(queues).Field(i).String())

		// making sure no queue name is left blank
		if queueName == "" {
			return nil, fmt.Errorf(errNoQueueNameProvided, field)
		}

		// checking if the queue has been registered before ?
		if true != isUniqueName(queueName, queuesNamesList) {
			return nil, fmt.Errorf(errQueueNameNotUnique, queueName)
		}
		queuesNamesList = append(queuesNamesList, queueName)
	}
	return queuesNamesList, nil
}

//isUniqueName checks if the queue name is unique and has not been used twice as an alias for other queues (to avoid misuse)
func isUniqueName(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return false
		}
	}
	return true
}
