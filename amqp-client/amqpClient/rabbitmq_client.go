package amqpClient

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var enableLogChan = false

func EnableLogChan() bool {
	return enableLogChan
}

/*有上限的，最大是1024个*/
var workersResultStream = make(chan string, 1024)

//RabbitMqClient is the default Message Broker client that wraps streadway/amqp implementation of rabbitmq api
//streadway/amqp can be easily replaced with other implementations, all is need is a struct of type MessageBrokerClientInterface
type RabbitMqClient struct {
	channel *amqp.Channel
	itemChanSlice[] RBElement
}
/*https://stackoverflow.com/questions/37690666/create-a-slice-of-buffered-channel-in-golang*/
type RBElement chan Element
var (
	errFailedToConnectToRabbitMq = "no name provided for queue : %s"
	errWorkersRegisteredForQueue = "no workers were registered for queue %s "
	errFailedToCreateConsumer    = "couldn' create consumer channel for queue %s "
	errFailedToPublishMessage    = "Failed to publish a message %s"
	errFailedToCreateQueue       = "Failed to create queue : %s"
)
/*每个包的init都是自动被调用的*/
func init() {
	// initialize the generation of pseudo-random sequence to be used with the hash
	rand.Seed(time.Now().UnixNano())
}

//Connect connects the client to the amqp server
func (cl *RabbitMqClient) Connect() error {
	// creating a new connection to the amqp
	conn, err := amqp.Dial(connectionURL())
	if err != nil {
		logError(err, errFailedToConnectToRabbitMq)
		return err
	}
	// creating a connection channel (aka socket)  /*连接用的channel又名socket*/
	ch, err := conn.Channel()
	if err != nil {
		logError(err, errFailedToConnectToRabbitMq)
		return err
	}

	cl.channel = ch
	return nil
}

/*断开到服务器的连接*/
//Disconnect closes the connection (channel) with the amqp server
func (cl *RabbitMqClient) Disconnect() error {
	return cl.channel.Close()
}

//CreateQueue is used to create a queue, by default all the queues created are durable, which means even if the message
//broker restarts, it will automatically recreate the queues upon recovery without data loss
func (cl *RabbitMqClient) CreateQueue(queueName string) error {
	_, err := cl.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	fmt.Println("creating queue ", queueName)
	if err != nil {
		err = fmt.Errorf(errFailedToCreateQueue, queueName)
		log.Println(err)
		return err
	}
	return nil
}

//Publish : publishes an event of type SerializableEventInterface to a specific queue
func (cl *RabbitMqClient) Publish(event SerializableEventInterface, queueName string) error {
	body, _ := event.Serialize()
	var channel = cl.channel

	// Publishing to a channel
	err := channel.Publish(
		"",        // all queues bind to the default exchange
		queueName, // queueName is used as a routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		log.Println(fmt.Sprintf(errFailedToPublishMessage, queueName))
		return err
	}
	return nil
}
//Consume is used to consume events / messages from a specific queue.
//A callback function is required as a parameter, and it' called whn a new message has been received
func (cl *RabbitMqClient) Consume2(queueName string, workers uint, callback func(items chan Element, consumer string)) error {
	fmt.Println("------------begin of Consume2----------")

	if workers < 1 {
		err := fmt.Errorf(errWorkersRegisteredForQueue, queueName)
		log.Println(err, nil)
		return err
	}
	for i := 1; i <= int(workers); i++ {
		cl.createConsumer(queueName, callback, i)
	}
	/*TODO 这里也会执行啊，consume会退出*/
	fmt.Println("------------end of Consume2----------")
	return nil
}
//creates a new consumer
func (cl *RabbitMqClient) createConsumer(queueName string, callback func(itemChan chan Element, consumer string), rank int) error {
	fmt.Println("------------begin of createConsumer------rank ----",rank)
	consumerHash := consumerHashValue(queueName, rank)
	log.Println("--gen consumerHash ok ",consumerHash)
	msgs, err := cl.channel.Consume(
		queueName,    // queue
		consumerHash, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	log.Println("---after cl.channel.Consume ")
	log.Println("make chan for rank ",rank)
	if false {
		cl.itemChanSlice[rank]=make(chan Element,100)
	}

	//items = make(chan Msg,100)
	if err != nil {
		err = fmt.Errorf(errFailedToCreateConsumer, queueName)
		return err
	}

	log.Println(consumerHash, " .. was created")
	go func() {
		/*这个协程不退出都是因为logOnHandleSuccess里有chan阻塞？ TODO */
		for msg := range msgs {
			if false {
				item := Element{
					Body:msg.Body,
					ack:msg.Ack,
					nack:msg.Nack,
				}
				//	var items chan Element
				cl.itemChanSlice[rank] <- item
				callback(cl.itemChanSlice[rank], consumerHash)
			}else{
				log.Println("msg:",msg)
			}

			if enableLogChan == true {
				/*可以先不写入，免得阻塞？好像不行，不断的写入chan中*/
				logOnHandleSuccess(msg.Body)
			}

		}
	}()
	/*TODO 这里会执行啊newConsumer也会退出，只能靠阻塞的协程来搞起*/
	fmt.Println("------------end of createConsumer----------")

	return nil
}
//Consume is used to consume events / messages from a specific queue.
//A callback function is required as a parameter, and it' called whn a new message has been received
func (cl *RabbitMqClient) Consume(queueName string, workers uint, callback func(msg []byte, consumer string)) error {
	fmt.Println("begin of Consume")
	if workers < 1 {
		err := fmt.Errorf(errWorkersRegisteredForQueue, queueName)
		log.Println(err, nil)
		return err
	}
	fmt.Println("Consume call newConsumer ")
	for i := 1; i <= int(workers); i++ {
		cl.newConsumer(queueName, callback, i)
	}
	/*TODO 这里也会执行啊，consume会退出*/
	fmt.Println("------------end of Consume----------")
	return nil
}

//creates a new consumer
func (cl *RabbitMqClient) newConsumer(queueName string, callback func(msg []byte, consumer string), rank int) error {
	consumerHash := consumerHashValue(queueName, rank)
	fmt.Println("------------begin of newConsumer------rank ----",rank,"consumerHash :",consumerHash)
	msgs, err := cl.channel.Consume(
		queueName,    // queue
		consumerHash, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		err = fmt.Errorf(errFailedToCreateConsumer, queueName)
		return err
	}

	log.Println(consumerHash, " .. was created")
	go func() {
		/*这个协程不退出都是因为logOnHandleSuccess里有chan阻塞？ TODO */
		for msg := range msgs {
			callback(msg.Body, consumerHash)
			if enableLogChan == true {
				/*可以先不写入，免得阻塞？好像不行，不断的写入chan中*/
				logOnHandleSuccess(msg.Body)
			}

		}
	}()
	/*TODO 这里会执行啊newConsumer也会退出，只能靠阻塞的协程来搞起*/
	fmt.Println("------------end of newConsumer----------")

	return nil
}
/*写入chan*/
//logOnHandleSuccess writes to the workersResultStream (logs results) for monitoring purposes
func logOnHandleSuccess(message []byte) {
	go func() {
		workersResultStream <- fmt.Sprintf("%v", string(message))
	}()
}
/*读取chan TODO 非常重要,没有这个consumer只能做一次*/
//第一个作用：是所有被成功消费的日志的复用器，处理系统的events，
//WatchWorkersStream does a couple of things : it's a multiplexer of the logs of all the successfully consumed and handled
//events in the system, it also blocks the function calling it from exiting the consumers.  第二个作用，就是阻塞函数，以阻止退出consumer
//Not calling this function would make it necessary for the developer to create consumers with a go routine [ go client.Consume(..)]
//如果不调用这个函数，则开发者非常有必要使用一个go routine去创建consumers
//and use the wait method from the builtin asynch package to manage concurrency in order to prevent the program from exiting.
func (cl *RabbitMqClient) WatchWorkersStream() {
	for {
		/*这个是取数据出来*/
		//log.Println(fmt.Sprintf("successfully handled : %s", <-workersResultStream))
		/*不打印，只取也可以吧*/
		<-workersResultStream
	}
}

//consumerHashValue return a hash value "a unique name" of the consumer with a suffix that is a random value
//to make sure 100% that no consumers would be created with the same name
//P.S.creating a consumer with the same name more than once would fail tha application
func consumerHashValue(queueName string, rank int) string {
	n := 6 //random pseudo-suffix size
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	consumerHash := queueName + "_consumer_" + strconv.Itoa(rank) + "__" + string(b)
	return consumerHash
}
/*	uri:="amqp://crs_video_app:jafepu123485%40%40%21%21%25%25@183.36.121.50:5672/%2fCRS"
*/
//connectionURL return the message broker connection url
func connectionURL() string {
	c := config.RabbitMq.Connection
	//return "amqp://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/"
	uri:= c.Uri
	return uri
}
