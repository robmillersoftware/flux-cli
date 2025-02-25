package tests

import (
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/robmillersoftware/flux-cli/core"
	"github.com/robmillersoftware/flux-cli/errors"
	"github.com/robmillersoftware/flux-cli/events"
	"github.com/robmillersoftware/flux-cli/plugins"
)

func TestErrorHandling(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "error_handling",
		ScenarioInitializer:  InitializeScenario,
		Options:              &godog.Options{Output: t},
	}
	if suite.Run() != 0 {
		t.Fatal("Test suite failed")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^"([^"]*)" is present with missing dependencies$`, faultyPluginPresent)
	ctx.Step(`^FlexiCLI attempts to load it$`, flexiCLIAttemptsToLoad)
	ctx.Step(`^an "([^"]*)" event is emitted with:$`, eventIsEmittedWith)
	ctx.Step(`^a valid command is received$`, validCommandReceived)
	ctx.Step(`^an internal exception occurs during processing$`, internalExceptionOccurs)
}

var eventBus *core.EventBus
var pluginLoader *plugins.PluginLoader
var receivedEvents []core.Event

func faultyPluginPresent(pluginName string) error {
	eventBus = core.NewEventBus()
	pluginLoader = plugins.NewPluginLoader(eventBus, nil)
	receivedEvents = []core.Event{}

	eventBus.Subscribe("ErrorOccurred", func(event core.Event) {
		receivedEvents = append(receivedEvents, event)
	})

	// Simulate faulty plugin
	plugin := plugins.Plugin{
		ID:          pluginName,
		Version:     "1.0.0",
		Description: "Faulty plugin with missing dependencies",
		Dependencies: []string{"missing-dependency"},
		Init: func(eventBus *core.EventBus, config map[string]interface{}) {
			errors.LogError(eventBus, "PLUGIN_LOAD_FAIL", "Missing dependency: X", "stack trace", map[string]interface{}{"pluginId": pluginName})
		},
	}
	pluginLoader.LoadPlugins = func(folder string) error {
		pluginLoader.plugins = append(pluginLoader.plugins, plugin)
		plugin.Init(eventBus, nil)
		return nil
	}
	return nil
}

func flexiCLIAttemptsToLoad() error {
	return pluginLoader.LoadPlugins("")
}

func eventIsEmittedWith(eventType string, table *godog.Table) error {
	for _, event := range receivedEvents {
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

func validCommandReceived() error {
	eventBus = core.NewEventBus()
	receivedEvents = []core.Event{}

	eventBus.Subscribe("ErrorOccurred", func(event core.Event) {
		receivedEvents = append(receivedEvents, event)
	})

	// Simulate valid command received
	event := core.Event{
		Type:      "CommandReceived",
		Timestamp: time.Now().Format(time.RFC3339),
		Payload: map[string]interface{}{
			"command":  "valid-command",
			"args":     []string{},
			"rawInput": "valid-command",
		},
	}
	eventBus.Publish(event)
	return nil
}

func internalExceptionOccurs() error {
	errors.LogError(eventBus, "COMMAND_EXECUTION_ERROR", "Null reference exception", "stack trace", map[string]interface{}{"command": "valid-command"})
	return nil
}
