amqp协议的封装，提供rabbitmq消费者和生产者接口。

# 生产者

```go
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
    
    sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
```

# 消费者

```go
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

func main() {
    // 1.创建consumer
	client, err := rabbitmq.NewMQConsumer("127.0.0.1", 5672, "guest", "guest")
	if err != nil {
		fmt.Println("rabbitmq.NewMQConsumer failed, err:", err.Error())
		return
	}
	fmt.Println("rabbitmq.NewMQConsumer success")

    // 2.consume
	err = client.Consume("first_queue", "first_exchange", "first.#", ACallback)
	if err != nil {
		fmt.Println("Consume failed, err:", err.Error())
		return
	}
	fmt.Println("Consume success")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}
```

