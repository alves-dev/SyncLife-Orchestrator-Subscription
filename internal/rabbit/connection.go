package rabbit

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func GetChannel() (*amqp.Channel, *amqp.Connection, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	conn, err := amqp.DialConfig(rabbitURL, amqp.Config{
		Properties: amqp.Table{
			"connection_name": "orchestrator-subscription",
		},
	})
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
