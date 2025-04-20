package rabbit

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func CreateQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return queue, fmt.Errorf("failed to declare queue: %w", err)
	}
	log.Printf("fila criada: %s\n", queueName)
	return queue, nil
}

func BindQueue(ch *amqp.Channel, queueName, routingKey string) error {
	exchangeName := os.Getenv("EXCHANGE_NAME")

	err := ch.QueueBind(
		queueName,    // nome da fila
		routingKey,   // routing key
		exchangeName, // exchange
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}
	return nil
}
