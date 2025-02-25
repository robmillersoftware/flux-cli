package events

import "time"

// CommandReceived event captures user input
type CommandReceived struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		Command  string   `json:"command"`
		Args     []string `json:"args"`
		RawInput string   `json:"rawInput"`
	} `json:"payload"`
}

// TechStackDetected event reports detected tech stack
type TechStackDetected struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		Stack          string            `json:"stack"`
		IndicatorFiles []string          `json:"indicatorFiles"`
		Metadata       map[string]string `json:"metadata"`
	} `json:"payload"`
}

// PluginLoaded event confirms a plugin is loaded
type PluginLoaded struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		PluginID     string   `json:"pluginId"`
		Version      string   `json:"version"`
		Description  string   `json:"description"`
		Dependencies []string `json:"dependencies"`
	} `json:"payload"`
}

// TestRunStarted event marks the start of a test run
type TestRunStarted struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		Suite     string    `json:"suite"`
		Mode      string    `json:"mode"`
		StartTime time.Time `json:"startTime"`
	} `json:"payload"`
}

// TestRunCompleted event conveys test results
type TestRunCompleted struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		Suite    string `json:"suite"`
		Result   struct {
			Passed  int `json:"passed"`
			Failed  int `json:"failed"`
			Skipped int `json:"skipped"`
		} `json:"result"`
		Duration int      `json:"duration"`
		Errors   []string `json:"errors"`
	} `json:"payload"`
}

// EnvironmentSwitched event indicates an environment change
type EnvironmentSwitched struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		Environment string            `json:"environment"`
		Config      map[string]string `json:"config"`
	} `json:"payload"`
}

// ErrorOccurred event logs errors
type ErrorOccurred struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		ErrorCode    string            `json:"errorCode"`
		ErrorMessage string            `json:"errorMessage"`
		StackTrace   string            `json:"stackTrace"`
		Context      map[string]string `json:"context"`
	} `json:"payload"`
}

// ConfigUpdated event notifies configuration changes
type ConfigUpdated struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	Payload   struct {
		ConfigFile string            `json:"configFile"`
		Changes    map[string]string `json:"changes"`
	} `json:"payload"`
}
