package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

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

// PluginLoader is responsible for loading plugins
type PluginLoader struct {
	plugins []Plugin
	eventBus *core.EventBus
	config   map[string]interface{}
}

// NewPluginLoader creates a new PluginLoader
func NewPluginLoader(eventBus *core.EventBus, config map[string]interface{}) *PluginLoader {
	return &PluginLoader{
		eventBus: eventBus,
		config:   config,
	}
}

// LoadPlugins scans the designated folder and loads plugins
func (loader *PluginLoader) LoadPlugins(folder string) error {
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "flexicli-plugin-") && strings.HasSuffix(info.Name(), ".so") {
			p, err := plugin.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open plugin %s: %v", path, err)
			}
			symPlugin, err := p.Lookup("Plugin")
			if err != nil {
				return fmt.Errorf("failed to find Plugin symbol in %s: %v", path, err)
			}
			pluginInstance, ok := symPlugin.(*Plugin)
			if !ok {
				return fmt.Errorf("invalid Plugin type in %s", path)
			}
			loader.plugins = append(loader.plugins, *pluginInstance)
			pluginInstance.Init(loader.eventBus, loader.config)
			loader.eventBus.Publish(core.Event{
				Type:      "PluginLoaded",
				Timestamp: core.GetCurrentTimestamp(),
				Payload: map[string]interface{}{
					"pluginId":     pluginInstance.ID,
					"version":      pluginInstance.Version,
					"description":  pluginInstance.Description,
					"dependencies": pluginInstance.Dependencies,
				},
			})
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load plugins: %v", err)
	}
	return nil
}

// DetectTechStack scans the project directory to detect the tech stack
func (loader *PluginLoader) DetectTechStack(projectDir string) (string, []string, error) {
	var stack string
	var indicatorFiles []string

	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			switch info.Name() {
			case "package.json":
				stack = "Node.js"
				indicatorFiles = append(indicatorFiles, "package.json")
			case "requirements.txt":
				stack = "Python"
				indicatorFiles = append(indicatorFiles, "requirements.txt")
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to detect tech stack: %v", err)
	}

	if stack == "" {
		stack = "unknown"
	}

	loader.eventBus.Publish(core.Event{
		Type:      "TechStackDetected",
		Timestamp: core.GetCurrentTimestamp(),
		Payload: map[string]interface{}{
			"stack":          stack,
			"indicatorFiles": indicatorFiles,
			"metadata":       map[string]interface{}{"version": "1.0.0"},
		},
	})

	return stack, indicatorFiles, nil
}
