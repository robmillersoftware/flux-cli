package main

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/robmillersoftware/flux-cli/core"
	"github.com/robmillersoftware/flux-cli/plugins"
	"github.com/robmillersoftware/flux-cli/config"
	"github.com/robmillersoftware/flux-cli/errors"
	"github.com/robmillersoftware/flux-cli/events"
)

func main() {
	// Initialize core engine, event bus, and configuration manager
	eventBus := core.NewEventBus()
	eventStore := core.NewEventStore()
	configManager := config.NewConfigManager(".flexicli.json")

	// Load configuration
	config, err := configManager.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		eventBus.Publish(events.NewErrorOccurredEvent("CONFIG_LOAD_FAIL", err.Error(), "", nil))
		return
	}

	// Detect tech stack and load plugins
	techStack, err := plugins.DetectTechStack(config)
	if err != nil {
		fmt.Println("Error detecting tech stack:", err)
		eventBus.Publish(events.NewErrorOccurredEvent("TECH_STACK_DETECT_FAIL", err.Error(), "", nil))
		return
	}
	eventBus.Publish(events.NewTechStackDetectedEvent(techStack))

	plugins, err := plugins.LoadPlugins(config, eventBus)
	if err != nil {
		fmt.Println("Error loading plugins:", err)
		eventBus.Publish(events.NewErrorOccurredEvent("PLUGIN_LOAD_FAIL", err.Error(), "", nil))
		return
	}

	// Start REPL or parse CLI arguments
	if len(os.Args) > 1 {
		command := os.Args[1]
		args := os.Args[2:]
		rawInput := fmt.Sprintf("%s %s", command, args)
		eventBus.Publish(events.NewCommandReceivedEvent(command, args, rawInput))
	} else {
		startREPL(eventBus)
	}
}

func startREPL(eventBus *core.EventBus) {
	for {
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		eventBus.Publish(events.NewCommandReceivedEvent(input, nil, input))
	}
}
