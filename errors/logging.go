package errors

import (
	"fmt"
	"time"
	"github.com/robmillersoftware/flux-cli/core"
	"github.com/robmillersoftware/flux-cli/events"
)

// LogError logs an error and emits an ErrorOccurred event
func LogError(eventBus *core.EventBus, errorCode, errorMessage, stackTrace string, context map[string]interface{}) {
	fmt.Printf("Error: %s - %s\n", errorCode, errorMessage)
	eventBus.Publish(events.NewErrorOccurredEvent(errorCode, errorMessage, stackTrace, context))
}

// NewErrorOccurredEvent creates a new ErrorOccurred event
func NewErrorOccurredEvent(errorCode, errorMessage, stackTrace string, context map[string]interface{}) core.Event {
	return core.Event{
		Type:      "ErrorOccurred",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: map[string]interface{}{
			"errorCode":    errorCode,
			"errorMessage": errorMessage,
			"stackTrace":   stackTrace,
			"context":      context,
		},
	}
}
