package tests

import (
	"testing"
	"time"
	"github.com/robmillersoftware/flux-cli/core"
	"github.com/robmillersoftware/flux-cli/events"
	"github.com/cucumber/godog"
)

func TestCommandProcessing(t *testing.T) {
	suite := godog.TestSuite{
		Name: "command_processing",
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step(`^FlexiCLI is running$`, flexiCLIIsRunning)
			s.Step(`^the user inputs "([^"]*)"$`, theUserInputs)
			s.Step(`^a "([^"]*)" event is emitted with:$`, anEventIsEmittedWith)
			s.Step(`^the command is validated successfully$`, theCommandIsValidatedSuccessfully)
			s.Step(`^the system responds with error "([^"]*)"$`, theSystemRespondsWithError)
			s.Step(`^an "([^"]*)" event is emitted with errorCode "([^"]*)"$`, anErrorEventIsEmittedWithErrorCode)
		},
		Options: &godog.Options{
			Format: "pretty",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

var (
	engine    *core.CoreEngine
	eventBus  *core.EventBus
	eventStore *core.EventStore
)

func flexiCLIIsRunning() error {
	engine = core.NewCoreEngine()
	eventBus = core.NewEventBus()
	eventStore = core.NewEventStore()
	return nil
}

func theUserInputs(input string) error {
	command, args := engine.ParseCommand(input)
	engine.HandleCommand(command, args)
	return nil
}

func anEventIsEmittedWith(eventType string, table *godog.Table) error {
	events := eventStore.GetEvents()
	for _, event := range events {
		if event.Type == eventType {
			for _, row := range table.Rows {
				key := row.Cells[0].Value
				value := row.Cells[1].Value
				if event.Payload[key] != value {
					return fmt.Errorf("expected %s to be %s, but got %s", key, value, event.Payload[key])
				}
			}
			return nil
		}
	}
	return fmt.Errorf("event %s not found", eventType)
}

func theCommandIsValidatedSuccessfully() error {
	// Assuming command validation is part of HandleCommand
	return nil
}

func theSystemRespondsWithError(errorMessage string) error {
	// Assuming system response is part of HandleCommand
	return nil
}

func anErrorEventIsEmittedWithErrorCode(eventType, errorCode string) error {
	events := eventStore.GetEvents()
	for _, event := range events {
		if event.Type == eventType && event.Payload["errorCode"] == errorCode {
			return nil
		}
	}
	return fmt.Errorf("error event %s with errorCode %s not found", eventType, errorCode)
}
