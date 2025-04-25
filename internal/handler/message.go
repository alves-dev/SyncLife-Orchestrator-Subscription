package handler

import (
	"encoding/json"
	"log"
	"orchestrator/internal/counter"
	"orchestrator/internal/rabbit"
	"orchestrator/pkg/events"
	"os"

	"github.com/streadway/amqp"
)

var eventCounter = counter.NewDailyCounter()

func HandleCountEvent(msg []byte, ch *amqp.Channel) {
	var event events.Deprecated
	err := json.Unmarshal(msg, &event)
	if err != nil {
		eventCounter.Increment("")
		return
	}

	eventCounter.Increment(event.Type)
}

func HandleSubscriptionEvent(msg []byte, ch *amqp.Channel) {
	var rawEvent events.Base
	err := json.Unmarshal(msg, &rawEvent)
	if err != nil {
		log.Printf("[HandleSubscriptionEvent] Failed to parse event envelope: %v", err)
		return
	}

	if rawEvent.Type != events.EventTypeSubscriptionRequested {
		return
	}

	var subscription events.SubscriptionRequestedV1
	err = json.Unmarshal(rawEvent.Data, &subscription)
	if err != nil {
		log.Printf("[HandleSubscriptionEvent] Failed to parse data for subscription event: %v", err)
		return
	}

	log.Printf("[HandleSubscriptionEvent] Received Event: %+v\n", subscription)

	rabbit.CreateQueue(ch, subscription.QueueName)

	exchangeName := os.Getenv("EXCHANGE_EVENTS_NAME")
	for _, key := range subscription.Subscriptions.EventTypes {
		rabbit.BindQueue(ch, exchangeName, subscription.QueueName, key)
	}
}

func HandleDeprecatedEvent(msg []byte, ch *amqp.Channel) {
	var event events.Deprecated
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Printf("[HandleDeprecatedEvent] Received deprecated Event")
		return
	}

	log.Printf("[HandleDeprecatedEvent] Received deprecated Event: %+v", event.Type)
}
