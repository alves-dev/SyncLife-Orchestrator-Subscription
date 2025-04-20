package events

import (
	"encoding/json"
	"time"
)

const EventTypeSubscriptionRequested = "orchestrator.subscriptions.requested.v1"

type Base struct {
	SpecVersion string                 `json:"specversion"`
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	Id          string                 `json:"id"`
	Time        time.Time              `json:"time"`
	DataSchema  string                 `json:"dataschema"`
	Data        json.RawMessage        `json:"data"`
	Extensions  map[string]interface{} `json:"extensions"`
}

type SubscriptionRequestedV1 struct {
	ServiceId     string `json:"service_id"`
	QueueName     string `json:"queue_name"`
	Subscriptions struct {
		EventTypes []string `json:"event_types"`
	} `json:"subscriptions"`
}
