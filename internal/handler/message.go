package handler

import (
	"encoding/json"
	"log"
	"orchestrator/internal/rabbit"
	"orchestrator/pkg/events"

	"github.com/streadway/amqp"
)

func HandleMessage(msg []byte, ch *amqp.Channel) {
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

	for _, key := range subscription.Subscriptions.EventTypes {
		rabbit.BindQueue(ch, subscription.QueueName, key)
	}

	log.Printf("✅ Received Event:\n%+v\n", subscription)
}
