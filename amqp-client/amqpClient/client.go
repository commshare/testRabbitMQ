package amqpClient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var config Config

var (
	errUnableToCreateLogsFile = "unable to create logs file %s"
)

// New function creates a new client with a new connection (channel) to the message broker
func New() (MessageBrokerClientInterface, error) {
	var cl MessageBrokerClientInterface
	cl = getMessageBrokerClient()
	err := cl.Connect()
	return cl, err
}

// Initialize : function takes care of two things : it parses the config of the amqb system,
// and creates the queues as binds them to the default exchange
func Initialize(configFile string, queues interface{}) error {

	// creating configuring of the message broker
	c, err := newAppConfiguration(configFile)
	if err != nil {
		log.Println(err)
		return err
	}
	config = *c

	/*TODO 关闭输出到文件*/
	/*输出到文件？*/
	// setting logs output stream file
	//setLoggerFile()

	// creating a new client
	cl, err := New()
	if err != nil {
		return err
	}

	// validating, extracting queues names and creating them
	queuesNames, err := extractQueuesNames(queues)
	log.Printf("queuesNames %s",queuesNames)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	for _, queue := range queuesNames {
		cl.CreateQueue(queue)
	}

	return nil
}

func logError(err error, msg string) {
	if err != nil {
		log.Println(fmt.Sprintf("%s  -- details : %s", msg, err))
	}
}

// This method defines the which implementation (of type MessageBrokerClientInterface) is used as a client
func getMessageBrokerClient() MessageBrokerClientInterface {
	cl := &RabbitMqClient{}
	fmt.Println("-----make itemChanSlice")
	/*TODO 在这里分配内存，初始化为2个元素的slice，https://stackoverflow.com/questions/37690666/create-a-slice-of-buffered-channel-in-golang*/
	//cl.itemChanSlice=make([]RBElement,2)
	return cl
}

// sets log file
func setLoggerFile() {
	filePath, err := filepath.Abs(config.RabbitMq.Logs.Logfile)
	if err != nil {
		logError(err, fmt.Sprintf(errUnableToCreateLogsFile, filePath))
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logFile, err := os.Create(filePath)
		if err != nil {
			logError(err, fmt.Sprintf(errUnableToCreateLogsFile, filePath))
		}
		log.SetOutput(logFile)
		return
	}
	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	/*输出到文件了*/
	log.SetOutput(logFile)
}
