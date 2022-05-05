/**
* @Author: hl
* @Description:
* @Date: 2022/4/30 17:56
 */

package rabbitmq

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
)

// Producer对象
type MQProducer struct {
	URI     	amqp.URI
	conn    	*amqp.Connection
	channel  	*amqp.Channel
	exchanges	[]string
}

// 创建MQProducer对象
func NewMQProducer(ip string, port int, username, password string) (*MQProducer, error) {
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

	producer := &MQProducer{
		URI:     uri,
		conn:    conn,
		channel: channel,
		exchanges: []string{},
	}

	return producer, nil
}

func (m *MQProducer) AddExchanges(exchanges []string) error {
	if m.channel == nil {
		return errors.New("channel is not open")
	}

	for _, e := range exchanges {
		err := m.channel.ExchangeDeclare(e, "topic", true, false, false, false, nil)
		if err != nil {
			return err
		}

		m.exchanges =  append(m.exchanges, e)
	}

	return nil
}

// publish
func (m *MQProducer) Publish(exchange, key string, msg []byte) error {
	if m.channel == nil {
		return errors.New("channel is not open")
	}

	err := m.channel.Publish(exchange, key, false, false,
		amqp.Publishing{
		ContentType: "text/plain",
		Body: msg,
	})
	if err != nil {
		return err
	}

	return nil
}