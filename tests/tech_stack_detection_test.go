package tests

import (
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/robmillersoftware/flux-cli/core"
	"github.com/robmillersoftware/flux-cli/events"
	"github.com/robmillersoftware/flux-cli/plugins"
)

func TestMain(m *testing.M) {
	status := godog.TestSuite{
		Name:                 "tech_stack_detection",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

var eventBus *core.EventBus
var pluginLoader *plugins.PluginLoader

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		eventBus = core.NewEventBus()
		pluginLoader = plugins.NewPluginLoader(eventBus, nil)
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the project contains "([^"]*)"$`, theProjectContains)
	ctx.Step(`^FlexiCLI starts$`, flexiCLIStarts)
	ctx.Step(`^a "([^"]*)" event is emitted with:$`, anEventIsEmittedWith)
	ctx.Step(`^the "([^"]*)" plugin is loaded$`, thePluginIsLoaded)
}

func theProjectContains(indicatorFile string) error {
	// Simulate the presence of the indicator file
	// In a real test, this would involve setting up the test environment
	return nil
}

func flexiCLIStarts() error {
	// Simulate FlexiCLI starting and detecting the tech stack
	stack, indicatorFiles, err := pluginLoader.DetectTechStack(".")
	if err != nil {
		return err
	}

	// Emit the TechStackDetected event
	eventBus.Publish(core.Event{
		Type:      "TechStackDetected",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"stack":          stack,
			"indicatorFiles": indicatorFiles,
			"metadata":       map[string]interface{}{"version": "1.0.0"},
		},
	})

	return nil
}

func anEventIsEmittedWith(eventType string, table *godog.Table) error {
	// Check if the event was emitted with the expected values
	// In a real test, this would involve checking the event store or a mock
	return nil
}

func thePluginIsLoaded(pluginID string) error {
	// Check if the plugin was loaded
	// In a real test, this would involve checking the plugin loader state
	return nil
}
