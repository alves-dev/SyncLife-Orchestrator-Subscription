package rabbit

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func GetChannel() (*amqp.Channel, *amqp.Connection, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open channel: %s", err)
	}

	return ch, conn, nil
}
