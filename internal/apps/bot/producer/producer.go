package producer

import (
	"chatroom/config"
	"github.com/streadway/amqp"
	"log"
	"log/slog"
	"os"
)

type Producer struct{}

func NewProducer(config *config.Config) *Producer {
	return &Producer{}
}

func Init() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		slog.Error("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"jobs",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("failed to declare a queue: %v", err)
	}

	body := "Test!"
	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		slog.Error("failed to publish a message: %v", err)
	}

	slog.Info("message sent: %v", body)
}
