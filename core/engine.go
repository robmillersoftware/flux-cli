package core

import (
	"fmt"
	"strings"
	"time"
)

// Event represents a generic event in the system
type Event struct {
	Type      string
	Timestamp time.Time
	Payload   map[string]interface{}
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

// CoreEngine represents the core engine of FlexiCLI
type CoreEngine struct {
	eventStore *EventStore
}

// NewCoreEngine creates a new core engine
func NewCoreEngine() *CoreEngine {
	return &CoreEngine{
		eventStore: NewEventStore(),
	}
}

// ParseCommand parses a command and its arguments
func (engine *CoreEngine) ParseCommand(input string) (string, []string) {
	// For simplicity, split the input by spaces
	parts := strings.Split(input, " ")
	command := parts[0]
	args := parts[1:]
	return command, args
}

// HandleCommand handles a command and dispatches events
func (engine *CoreEngine) HandleCommand(command string, args []string) {
	// Create a CommandReceived event
	event := Event{
		Type:      "CommandReceived",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"command":  command,
			"args":     args,
			"rawInput": fmt.Sprintf("%s %s", command, strings.Join(args, " ")),
		},
	}

	// Add the event to the store
	engine.eventStore.AddEvent(event)

	// Dispatch the event (for now, just print it)
	fmt.Printf("Event dispatched: %+v\n", event)
}
