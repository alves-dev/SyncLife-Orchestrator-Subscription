package handler

import (
	"encoding/json"
	"log"
	"orchestrator/internal/rabbit"
	"orchestrator/pkg/events"
	"os"

	"github.com/streadway/amqp"
)

func HandleSubscriptionEvent(msg []byte, ch *amqp.Channel) {
	var rawEvent events.Base
	err := json.Unmarshal(msg, &rawEvent)
	if err != nil {
		log.Printf("Failed to parse event envelope: %v", err)
		return
	}

	var subscription events.SubscriptionRequestedV1
	switch rawEvent.Type {
	case events.EventTypeSubscriptionRequested:
		err := json.Unmarshal(rawEvent.Data, &subscription)
		if err != nil {
			log.Printf("Failed to parse data for subscription event: %v", err)
			return
		}
		log.Printf("Subscription request received: %+v\n", subscription)

	default:
		log.Printf("Unhandled event type: %s", rawEvent.Type)
	}

	rabbit.CreateQueue(ch, subscription.QueueName)

	exchangeName := os.Getenv("QUEUE_SUBSCRIPTION_NAME")
	for _, key := range subscription.Subscriptions.EventTypes {
		rabbit.BindQueue(ch, exchangeName, subscription.QueueName, key)
	}

	log.Printf("✅ Received Event:\n%+v\n", subscription)
}

func HandleDeprecatedEvent(msg []byte, ch *amqp.Channel) {
	var event events.Deprecated
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Printf("Failed to parse event envelope: %v", err)
		return
	}

	log.Printf("Received deprecated Event: %+v", event.Type)
}
