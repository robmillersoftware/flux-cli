package plugins

import (
	"github.com/robmillersoftware/flux-cli/core"
)

// Plugin represents a dynamically loaded plugin
type Plugin struct {
	ID          string
	Version     string
	Description string
	Dependencies []string
	Init        func(eventBus *core.EventBus, config map[string]interface{})
}

// RegisterCommand registers a command with the event bus
func RegisterCommand(eventBus *core.EventBus, command string, handler func(args []string)) {
	eventBus.Subscribe("CommandReceived", func(event core.Event) {
		if event.Payload["command"] == command {
			handler(event.Payload["args"].([]string))
		}
	})
}

// EmitEvent emits an event to the event bus
func EmitEvent(eventBus *core.EventBus, eventType string, payload map[string]interface{}) {
	eventBus.Publish(core.Event{
		Type:      eventType,
		Timestamp: core.GetCurrentTimestamp(),
		Payload:   payload,
	})
}

// Log logs a message to the event bus
func Log(eventBus *core.EventBus, message string) {
	eventBus.Publish(core.Event{
		Type:      "Log",
		Timestamp: core.GetCurrentTimestamp(),
		Payload: map[string]interface{}{
			"message": message,
		},
	})
}
