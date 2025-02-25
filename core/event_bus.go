package core

import (
	"sync"
)

// Event represents a generic event in the system
type Event struct {
	Type      string
	Timestamp string
	Payload   map[string]interface{}
}

// EventHandler is a function that handles an event
type EventHandler func(event Event)

// EventBus is the central hub for publishing and subscribing to events
type EventBus struct {
	subscribers map[string][]EventHandler
	mutex       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Publish publishes an event to all subscribers
func (bus *EventBus) Publish(event Event) {
	bus.mutex.RLock()
	defer bus.mutex.RUnlock()

	if handlers, found := bus.subscribers[event.Type]; found {
		for _, handler := range handlers {
			go handler(event)
		}
	}
}

// Subscribe subscribes a handler to an event type
func (bus *EventBus) Subscribe(eventType string, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
}

// Unsubscribe unsubscribes a handler from an event type
func (bus *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	if handlers, found := bus.subscribers[eventType]; found {
		for i, h := range handlers {
			if h == handler {
				bus.subscribers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}
