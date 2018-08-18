package common

import (
	"github.com/streadway/amqp"
	"theAmazingSmsSender/app/config"
)

var rabbitMqConnection *amqp.Connection

func ConnectToRabbitMQ() {
	connection, err := amqp.Dial("amqp://" + config.GetConfig().RABBITMQ_USER + ":" + config.GetConfig().RABBITMQ_PASSWORD + "@" + config.GetConfig().RABBITMQ_HOST + ":" + config.GetConfig().RABBITMQ_PORT + "/")
	if err != nil {
		panic(err)
	}

	rabbitMqConnection = connection
}

func GetRabbitMQChannel() *amqp.Channel {
	ch, err := rabbitMqConnection.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
