package broker

import (
	"chatroom/config"
	"fmt"
	"github.com/streadway/amqp"
	"log/slog"
)

type Channel struct {
	queueName string
	ch        *chan []byte
}

type Broker struct {
	conn     *amqp.Connection
	Channels map[string]*Channel
}

func NewBroker(config *config.Config) *Broker {
	conn, err := amqp.Dial(config.Exchanger.Host)
	if err != nil {
		panic(err)
	}
	return &Broker{
		conn: conn,
	}
}

func (broker *Broker) Consume(queueName string, receiver *chan []byte) error {
	if broker.conn.IsClosed() {
		slog.Error("exchanger connection is closed")
		return nil
	}

	ch, err := broker.conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("Error closing channel: %s", err.Error()))
		}
	}()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to declare a queue: %v", err.Error()))
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to register a consumer: %v", err.Error()))
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			slog.Debug(fmt.Sprintf("Received a message: %s", string(d.Body)))
			msg := d.Body
			*receiver <- msg
		}
	}()
	<-forever
	return nil
}

func (broker *Broker) Produce(queueName string, sender *chan []byte) error {
	if broker.conn.IsClosed() {
		slog.Error("exchanger connection is closed")
		return nil
	}

	ch, err := broker.conn.Channel()
	if err != nil {
		return err
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("Error closing channel: %s", err.Error()))
		}
	}()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to declare a queue: %v", err.Error()))
		return err
	}

	for {
		select {
		case msg := <-*sender:
			fmt.Println(string(msg))
			err := ch.Publish("", "key", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			})
			if err != nil {
				slog.Error(fmt.Sprintf("failed to publish a message: %v", err.Error()))
				continue
			}
		}
	}
}

func (broker *Broker) NewConsumer(queueName string) (*Channel, error) {
	ch := &Channel{
		queueName: queueName,
		ch:        new(chan []byte),
	}
	err := broker.Consume(ch.queueName, ch.ch)
	if err != nil {
		return nil, err
	}
	broker.Channels[queueName] = ch
	return ch, nil
}

func (broker *Broker) NewProducer(queueName string) (*Channel, error) {
	ch := &Channel{
		queueName: queueName,
		ch:        new(chan []byte),
	}
	err := broker.Produce(ch.queueName, ch.ch)
	if err != nil {
		return nil, err
	}
	broker.Channels[queueName] = ch
	return ch, nil
}
