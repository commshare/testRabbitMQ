package main

import (
	"fmt"
	"log"
	"../../amqpClient"
)

//import (
//	"fmt"
//	"log"
//	"../../amqpClient"
//
//
//	/*
//	"github.com/almasry/amqp-client"
//	*/
//)
//
type Queues struct {
/*	QueueUserAccount    string
	QueueOrderWorkflow  string*/
	QueueServiceChanges string
}
//
var AmqpQueues Queues

/*TODO 这个init 会让打印不出到控制台啊 ,而且这个函数貌似会被自动调用啊*/
func init() {
	fmt.Println("call init ...")
	AmqpQueues = Queues{
/*		QueueUserAccount:    "user_account", // the actual name of the queue to be created by amqp
		QueueOrderWorkflow:  "order_workflow",*/
		QueueServiceChanges: "TEST_vodimg",
	}
	// initialize the amqb client
	// inject configuration and a struct that has all AmqpQueues names
	err := amqpClient.Initialize("V:\\RECORD_TEST\\refref\\config.yml", AmqpQueues)
	failOnError(err)
}

func realmAIN(){
	client, err := amqpClient.New()
	failOnError(err)
	defer client.Disconnect()
	fmt.Println("------new client ok -----")
	// getting a list of all the AmqpQueues

	client.Consume(AmqpQueues.QueueServiceChanges, 2, func(msg []byte, consumer string) {
		eventHandler2(msg, consumer)
	})

	//client.Consume(AmqpQueues.QueueUserAccount, 2, func(msg []byte, consumer string) {
	//	eventHandler(msg, consumer)
	//	// Do something about this event ..
	//})
	//
	//// Do some other business logic here ...
	//// ....
	//// ....
	//
	//client.Consume(AmqpQueues.QueueOrderWorkflow, 2, func(msg []byte, consumer string) {
	//	eventHandler(msg, consumer)
	//})

	// Do some other business logic here ...
	// ....
	// ....

	/*TODO 这里不加,会导致main退出，加了就能让consumer不断的执行消费数据*/
		client.WatchWorkersStream()

}
func realmAIN2(){
	client, err := amqpClient.New()
	failOnError(err)
	defer client.Disconnect()
	fmt.Println("------new client ok -----")
	// getting a list of all the AmqpQueues
	/*TODO	这说明里头的msg chan 是阻塞的，一直有数据自产自销并阻塞住*/
	go client.Consume(AmqpQueues.QueueServiceChanges, 2, func(msg []byte, consumer string) {
		eventHandler2(msg, consumer)
	})


}
func realmAIN3(){
	client, err := amqpClient.New()
	failOnError(err)
	defer client.Disconnect()
	fmt.Println("------new client ok -----")
	// getting a list of all the AmqpQueues
	/*TODO	这说明里头的msg chan 是阻塞的，一直有数据自产自销并阻塞住*/
	go client.Consume2(AmqpQueues.QueueServiceChanges, 2, func(itemChan chan amqpClient.Element, consumer string) {
		it := <-itemChan
		eventHandler3(it, consumer)
	})


}
func main(){
	fmt.Println("---main---new client  -----")
	fmt.Println("---main---new client  -----")
	fmt.Println("---main---new client  -----")
	fmt.Println("---main---new client  -----")
	log.Println("----------------------")

	/*TODO 这里不加,会导致main退出，加了就能让consumer不断的执行消费数据*/
	if amqpClient.EnableLogChan() == true {
		log.Println("realmAIN1")
		realmAIN()
	}else{
		fmt.Println("----not use log chan-------")
		/* //work
		realmAIN2()
		*/
		if false {
			log.Println("realmAIN3")
			realmAIN3()
		}else{
			log.Println("realmAIN2")
			realmAIN2()
		}

	}
	select{}
}

func eventHandler(msg []byte, consumer string) {
	// assuming you already know the type of event you'e receiving ..
	var event amqpClient.Event
	event.Deserialize(msg)

	log.Println(fmt.Sprintf("Consumer %s just Received a message : %s", consumer, event.Name))
}
func eventHandler2(msg []byte, consumer string) {
	// assuming you already know the type of event you'e receiving ..
/*	var event amqpClient.Event
	event.Deserialize(msg)*/

	log.Println(fmt.Sprintf("Consumer %s eventHandler2 just Received a message ", consumer))
}
func eventHandler3(item amqpClient.Element, consumer string) {
	// assuming you already know the type of event you'e receiving ..
	/*	var event amqpClient.Event
		event.Deserialize(msg)*/

	log.Println(fmt.Sprintf("Consumer %s eventHandler3 Received %s", consumer,item.Body))
}
func failOnError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}
