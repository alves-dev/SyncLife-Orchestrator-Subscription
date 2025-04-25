package main

import (
	"log"
	"orchestrator/internal/handler"
	"orchestrator/internal/rabbit"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	initLogger()
	loadAndValidEnvs()

	channel, connection, err := rabbit.GetChannel()

	if err != nil {
		log.Fatalf("RabbitMQ connection error: %v", err)
	}
	defer connection.Close()
	defer channel.Close()

	queueSubscription := os.Getenv("QUEUE_SUBSCRIPTION_NAME")
	queueDeprecatedEvents := os.Getenv("QUEUE_DEPRECATED_EVENTS_NAME")

	// Crete queues
	_, err = rabbit.CreateQueue(channel, queueSubscription)
	if err != nil {
		log.Fatalf("fila não criada: %v", err)
	}

	_, err = rabbit.CreateQueue(channel, queueDeprecatedEvents)
	if err != nil {
		log.Fatalf("fila não criada: %v", err)
	}

	// Create binds
	exchangeEvents := os.Getenv("EXCHANGE_EVENTS_NAME")
	exchangeDeprecatedNames := os.Getenv("EXCHANGE_DEPRECATED_NAMES")

	rabbit.BindQueue(channel, exchangeEvents, queueSubscription, "#")
	for _, name := range strings.Split(exchangeDeprecatedNames, ",") {
		rabbit.BindQueue(channel, name, queueDeprecatedEvents, "#")
	}

	// Start consuming
	msgsSubscription, err := channel.Consume(
		queueSubscription,
		"consumer-tag-subscription",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	msgsDeprecatedEvents, err := channel.Consume(
		queueDeprecatedEvents,
		"consumer-tag-deprecated-events",
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	// Consume loop
	go func() {
		for d := range msgsSubscription {
			handler.HandleCountEvent(d.Body, channel)
			handler.HandleSubscriptionEvent(d.Body, channel)
		}
	}()

	go func() {
		for d := range msgsDeprecatedEvents {
			handler.HandleCountEvent(d.Body, channel)
			handler.HandleDeprecatedEvent(d.Body, channel)
		}
	}()

	select {} // bloqueia pra sempre
}

func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.SetPrefix("[orchestrator] ")
}

func loadAndValidEnvs() {
	// Load .env file
	_ = godotenv.Load()

	requiredEnv := []string{"RABBITMQ_URL", "QUEUE_SUBSCRIPTION_NAME", "QUEUE_DEPRECATED_EVENTS_NAME", "EXCHANGE_EVENTS_NAME", "EXCHANGE_DEPRECATED_NAMES"}

	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required environment variable: %s", key)
		}
	}
}
