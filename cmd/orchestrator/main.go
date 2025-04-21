package main

import (
	"log"
	"orchestrator/internal/handler"
	"orchestrator/internal/rabbit"
	"os"

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

	queueName := os.Getenv("QUEUE_NAME")

	_, err = rabbit.CreateQueue(channel, queueName)
	if err != nil {
		log.Fatalf("fila não criada: %v", err)
	}

	// Start consuming
	msgs, err := channel.Consume(
		queueName,
		"",
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
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handler.HandleMessage(d.Body, channel)
		}
	}()

	log.Println("Listening for events...")
	<-forever
}

func initLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[orchestrator] ")
}

func loadAndValidEnvs() {
	// Load .env file
	_ = godotenv.Load()

	requiredEnv := []string{"RABBITMQ_URL", "QUEUE_NAME", "EXCHANGE_NAME"}

	for _, key := range requiredEnv {
		if os.Getenv(key) == "" {
			log.Fatalf("Missing required environment variable: %s", key)
		}
	}
}
