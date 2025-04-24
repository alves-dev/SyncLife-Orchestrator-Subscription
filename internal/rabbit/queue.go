package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func CreateQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	_ = deleteQueue(queueName)
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

func BindQueue(ch *amqp.Channel, exchangeName, queueName, routingKey string) error {
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

func deleteQueue(queueName string) bool {
	channel, connection, _ := GetChannel()
	defer connection.Close()
	defer channel.Close()

	count, err := channel.QueueDelete(
		queueName,
		false, // ifUnused: só deleta se não tiver consumidores
		true,  // ifEmpty: só deleta se não tiver mensagens
		false, // noWait
	)
	if err != nil {
		fmt.Printf("failed to delete queue: %s\n", queueName)
		return false
	}

	if count > 0 {
		fmt.Printf("Events deleted in queue %s: %d\n", queueName, count)
	}

	return true
}
