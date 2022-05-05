/**
* @Author: hl
* @Description:
* @Date: 2022/4/30 21:55
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"rabbitmq"
	"syscall"
)

type Student struct {
	Name	string
	Id		int
}

func main() {
	// 1.创建producer
	producer, err := rabbitmq.NewMQProducer("127.0.0.1", 5672, "guest", "guest")
	if err != nil {
		fmt.Println("rabbitmq.NewMQProducer failed, err:", err.Error())
		return
	}
	fmt.Println("rabbitmq.NewMQProducer success")

	// 2.添加exchanges
	exchanges := []string{"first_exchange", "second_exchange"}
	err = producer.AddExchanges(exchanges)
	if err != nil {
		fmt.Println("AddExchanges failed, err:", err.Error())
		return
	}
	fmt.Println("AddExchanges success")

	for i := 0; i < 1000000; i++ {
		student := &Student{
			Name: "hello",
			Id: i,
		}

		msg, err := json.Marshal(student)
		if err != nil {
			fmt.Println("json.Marshal failed, err:", err.Error())
			return
		}

		// 3.publish
		err = producer.Publish("first_exchange", "first.second.three", msg)
		if err != nil {
			fmt.Println("Publish failed, err:", err.Error())
			break
		}
		fmt.Printf("Publish success, msg:%s\n", msg)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}