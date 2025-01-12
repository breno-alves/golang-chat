package exchanger

import (
	"chatroom/config"
	"github.com/streadway/amqp"
)

func NewExchanger(config *config.Config) *amqp.Connection {
	conn, err := amqp.Dial(config.Exchanger.Host)
	if err != nil {
		panic(err)
	}
	return conn
}
