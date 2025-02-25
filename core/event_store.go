package core

import (
	"time"
)

// Event represents a generic event in the system
type Event struct {
	Type      string
	Timestamp time.Time
	Payload   map[string]interface{}
	Metadata  map[string]interface{}
}

// EventStore maintains a list of events
type EventStore struct {
	events []Event
}

// NewEventStore creates a new event store
func NewEventStore() *EventStore {
	return &EventStore{
		events: []Event{},
	}
}

// AddEvent adds an event to the store
func (store *EventStore) AddEvent(event Event) {
	store.events = append(store.events, event)
}

// GetEvents returns all events in the store
func (store *EventStore) GetEvents() []Event {
	return store.events
}

// ReplayEvents replays events based on a filter
func (store *EventStore) ReplayEvents(filter func(Event) bool) []Event {
	var filteredEvents []Event
	for _, event := range store.events {
		if filter(event) {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents
}
