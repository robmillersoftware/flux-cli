package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ConfigManager manages the configuration settings for FlexiCLI
type ConfigManager struct {
	configFile string
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(configFile string) *ConfigManager {
	return &ConfigManager{
		configFile: configFile,
	}
}

// LoadConfig loads the configuration from the config file
func (manager *ConfigManager) LoadConfig() (map[string]interface{}, error) {
	file, err := os.Open(manager.configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to the config file
func (manager *ConfigManager) SaveConfig(config map[string]interface{}) error {
	bytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := ioutil.WriteFile(manager.configFile, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
