/**
* @Author: hl
* @Description:
* @Date: 2022/4/30 21:07
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"rabbitmq"
	"syscall"
)

func ACallback(key string, msg []byte) error {
	fmt.Printf("ACallback start, key:%s msg:%s\n", key, msg)
	return nil
}

func BCallback(key string, msg []byte) error {
	fmt.Printf("BCallback start, key:%s msg:%s\n", key, msg)
	return nil
}

func main() {
	// 1.创建consumer
	client1, err := rabbitmq.NewMQConsumer("127.0.0.1", 5672, "guest", "guest")
	if err != nil {
		fmt.Println("rabbitmq.NewMQConsumer failed, err:", err.Error())
		return
	}
	fmt.Println("rabbitmq.NewMQConsumer success")

	// 2.consume
	err = client1.Consume("first_queue", "first_exchange", "first.#", ACallback)
	if err != nil {
		fmt.Println("Consume failed, err:", err.Error())
		return
	}
	fmt.Println("Consume success")

	client2, err := rabbitmq.NewMQConsumer("127.0.0.1", 5672, "guest", "guest")
	if err != nil {
		fmt.Println("rabbitmq.NewMQConsumer failed, err:", err.Error())
		return
	}
	fmt.Println("rabbitmq.NewMQConsumer success")

	err = client2.Consume("first_queue", "first_exchange", "first.#", BCallback)
	if err != nil {
		fmt.Println("Consume failed, err:", err.Error())
		return
	}
	fmt.Println("Consume success")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}