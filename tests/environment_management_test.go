package tests

import (
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/stretchr/testify/assert"
)

type EnvironmentManagementTest struct {
	eventBus *core.EventBus
}

func (emt *EnvironmentManagementTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the user inputs "([^"]*)"$`, emt.theUserInputs)
	ctx.Step(`^an "([^"]*)" event is emitted with:$`, emt.anEventIsEmittedWith)
}

func (emt *EnvironmentManagementTest) theUserInputs(input string) error {
	command, args := core.ParseCommand(input)
	emt.eventBus.Publish(core.Event{
		Type:      "CommandReceived",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"command":  command,
			"args":     args,
			"rawInput": input,
		},
	})
	return nil
}

func (emt *EnvironmentManagementTest) anEventIsEmittedWith(eventType string, table *godog.Table) error {
	expectedPayload := make(map[string]interface{})
	for _, row := range table.Rows {
		expectedPayload[row.Cells[0].Value] = row.Cells[1].Value
	}

	var emittedEvent core.Event
	for _, event := range emt.eventBus.GetEvents() {
		if event.Type == eventType {
			emittedEvent = event
			break
		}
	}

	assert.Equal(emt, expectedPayload, emittedEvent.Payload)
	return nil
}

func TestEnvironmentManagement(t *testing.T) {
	emt := &EnvironmentManagementTest{
		eventBus: core.NewEventBus(),
	}

	opts := godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
	}

	status := godog.TestSuite{
		Name:                 "environment_management",
		ScenarioInitializer:  emt.InitializeScenario,
		Options:              &opts,
	}.Run()

	if status != 0 {
		t.Fatalf("non-zero status returned, failed to run feature tests")
	}
}
