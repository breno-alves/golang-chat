package consumer

import (
	"github.com/streadway/amqp"
	"log/slog"
)

type Consumer struct {
	exchanger *amqp.Connection
}

func NewConsumer(exchanger *amqp.Connection) *Consumer {
	consumer := &Consumer{
		exchanger: exchanger,
	}
	consumer.Init()
	return consumer
}

func (c *Consumer) Init() {
	ch, err := c.exchanger.Channel()
	if err != nil {
		slog.Error("failed to open a channel: %v", err)
		return
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			slog.Error("failed to close channel: %v", err)
		}
	}()

	q, err := ch.QueueDeclare(
		"bot_response",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("failed to declare a queue: %v", err)
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
		slog.Error("failed to register a consumer: %v", err)
		return
	}

	go func() {
		for d := range msgs {
			slog.Info("message received: %v", string(d.Body))
		}
	}()
}
