/**
* @Author: hl
* @Description:
* @Date: 2022/4/30 17:57
 */

package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

// 消息回调
type MsgCallback func(key string, msg []byte) error

// Consumer对象
type MQConsumer struct {
	URI     	amqp.URI
	conn    	*amqp.Connection
	channel  	*amqp.Channel
	queue    	amqp.Queue
	delivery 	<-chan amqp.Delivery
	MsgProc  	MsgCallback
}

// 创建Consumer
func NewMQConsumer(ip string, port int, username, password string) (*MQConsumer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d", username, password, ip, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	uri := amqp.URI{
		Scheme: "amqp",
		Host: ip,
		Port: port,
		Username: username,
		Password: password,
		Vhost: "/",
	}

	consumer := &MQConsumer{
		URI:     uri,
		conn:    conn,
		channel: channel,
	}

	return consumer, nil
}

// 消费
func (m *MQConsumer) Consume(queue, exchange, key string, msgCallback MsgCallback) error {
	var err error

	if m.channel == nil {
		return errors.New("channel is not open")
	}

	err = m.channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	m.queue, err =  m.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = m.channel.QueueBind(m.queue.Name, key, exchange, false, nil)
	if err != nil {
		return err
	}

	m.delivery, err = m.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	m.MsgProc = msgCallback

	go m.handleMsg()

	return nil
}

// 处理消息
func (m *MQConsumer) handleMsg() {
	defer func() {
		if r := recover(); r != nil {
			m.handleMsg()
		}
	}()

	for d := range m.delivery {
		if m.MsgProc != nil {
			m.MsgProc(d.RoutingKey, d.Body)
		}
		d.Ack(false)
	}
}

// Close
func (m *MQConsumer) Close() error {
	var err error

	if m.channel != nil {
		_, err = m.channel.QueueDelete(m.queue.Name, true, true, false)
		if err != nil {
			return nil
		}

		err = m.channel.Close()
		if err != nil {
			return err
		}
	}

	if m.conn != nil {
		err = m.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}