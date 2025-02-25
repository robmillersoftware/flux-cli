# FlexiCLI

![CI](https://github.com/robmillersoftware/flux-cli/actions/workflows/ci.yml/badge.svg)

## Introduction & Goals

FlexiCLI is a minimal, cross‑platform CLI tool that auto‑adapts to a project's tech stack via dynamic plugins. Core behavior is driven by event modeling (every state change is an immutable event) and defined via BDD (Gherkin scenarios).

### Goals

- Minimal core
- Tech‑stack agnostic auto‑detection
- Extensible plugin ecosystem
- Robust auditability & error handling
- BDD‑driven development

## System Architecture

### Core Engine

- Parses commands
- Handles global options
- Dispatches events
- Maintains an event store

### Event Bus

- Central hub for publishing/subscribing events
- Supports sync and async

### Plugin Loader & Registry

- Scans designated folders/configurations to detect the tech stack
- Dynamically loads plugins

### Plugins

- Modular extensions (e.g., test runners, environment switchers)
- Implement an init() method and subscribe to events

### Configuration Manager

- Reads settings from a ".flexicli.json" file
- Supports manual overrides and environment-specific options

### Error & Logging Module

- Captures errors (via ErrorOccurred events)
- Provides fallback behavior

## Event Models

Each event has an "event" type, "timestamp", and "payload".

### CommandReceived

Captures user input.

```json
{
  "event": "CommandReceived",
  "timestamp": "ISO8601",
  "payload": {
    "command": "string",
    "args": ["string"],
    "rawInput": "string"
  }
}
```

### TechStackDetected

Reports detected tech stack.

```json
{
  "event": "TechStackDetected",
  "timestamp": "ISO8601",
  "payload": {
    "stack": "string",
    "indicatorFiles": ["string"],
    "metadata": { "version": "string" }
  }
}
```

### PluginLoaded

Confirms a plugin is loaded.

```json
{
  "event": "PluginLoaded",
  "timestamp": "ISO8601",
  "payload": {
    "pluginId": "string",
    "version": "string",
    "description": "string",
    "dependencies": ["string"]
  }
}
```

### TestRunStarted

Marks the start of a test run.

```json
{
  "event": "TestRunStarted",
  "timestamp": "ISO8601",
  "payload": {
    "suite": "string",
    "mode": "string",
    "startTime": "ISO8601"
  }
}
```

### TestRunCompleted

Conveys test results.

```json
{
  "event": "TestRunCompleted",
  "timestamp": "ISO8601",
  "payload": {
    "suite": "string",
    "result": {
      "passed": "number",
      "failed": "number",
      "skipped": "number"
    },
    "duration": "number",
    "errors": ["string"]
  }
}
```

### EnvironmentSwitched

Indicates an environment change.

```json
{
  "event": "EnvironmentSwitched",
  "timestamp": "ISO8601",
  "payload": {
    "environment": "string",
    "config": { "key": "value" }
  }
}
```

### ErrorOccurred

Logs errors.

```json
{
  "event": "ErrorOccurred",
  "timestamp": "ISO8601",
  "payload": {
    "errorCode": "string",
    "errorMessage": "string",
    "stackTrace": "string",
    "context": { "command": "string", "pluginId": "string" }
  }
}
```

### ConfigUpdated

Notifies configuration changes.

```json
{
  "event": "ConfigUpdated",
  "timestamp": "ISO8601",
  "payload": {
    "configFile": "string",
    "changes": { "key": "newValue" }
  }
}
```

## BDD Specifications (Gherkin Scenarios)

### Feature: Command Processing

#### Scenario: Valid command processing

Given FlexiCLI is running
When the user inputs "run-tests --verbose"
Then a "CommandReceived" event is emitted with:
  | command | run-tests |
  | args    | --verbose |
And the command is validated successfully

#### Scenario: Invalid command input

Given FlexiCLI is running
When the user inputs "run-magic"
Then a "CommandReceived" event is emitted with:
  | command | run-magic |
And the system responds with error "Unknown command: run-magic"
And an "ErrorOccurred" event is emitted with errorCode "INVALID_COMMAND"

### Feature: Tech Stack Detection & Plugin Loading

#### Scenario: Detect Node.js project

Given the project contains "package.json"
When FlexiCLI starts
Then a "TechStackDetected" event is emitted with:
  | stack          | Node.js      |
  | indicatorFiles | package.json |
And the "node-test-runner" plugin is loaded (emits "PluginLoaded" event)

#### Scenario: Unknown tech stack fallback

Given no recognizable tech stack files exist
When FlexiCLI starts
Then a "TechStackDetected" event is emitted with:
  | stack | unknown |
And only general-purpose commands are available

### Feature: Testing Workflow Integration

#### Scenario: Run tests in Node.js with verbose output

Given a Node.js project is detected and the "node-test-runner" plugin is loaded
When the user inputs "flexicli run-tests --verbose"
Then a "TestRunStarted" event is emitted with:
  | suite | unit    |
  | mode  | verbose |
And tests execute, followed by a "TestRunCompleted" event with:
  | result.passed | 42 |
  | result.failed | 0  |
  | duration      | (numeric value) |

#### Scenario: Test run error handling

Given a Node.js project is detected
When the user inputs "flexicli run-tests" and tests fail
Then a "TestRunStarted" event is emitted
And an "ErrorOccurred" event is emitted with:
  | errorMessage | Test suite failed |
And the system displays "Test run failed. Check logs for details."

### Feature: Environment Management

#### Scenario: Switch to development environment

Given environments "development" and "production" exist
When the user inputs "flexicli switch-env dev"
Then an "EnvironmentSwitched" event is emitted with:
  | environment | development |
And the corresponding environment variables are loaded

#### Scenario: Switch to a non-existent environment

Given only "dev" and "prod" are available
When the user inputs "flexicli switch-env staging"
Then an "ErrorOccurred" event is emitted with:
  | errorMessage | Environment 'staging' not found |
And the system displays "Environment 'staging' not found"

### Feature: Error Handling

#### Scenario: Plugin loading error due to missing dependency

Given "faulty-plugin.js" is present with missing dependencies
When FlexiCLI attempts to load it
Then an "ErrorOccurred" event is emitted with:
  | errorCode    | PLUGIN_LOAD_FAIL |
  | errorMessage | Missing dependency: X |
And the system logs the error and continues

#### Scenario: Command execution error due to internal exception

Given a valid command is received
When an internal exception occurs during processing
Then an "ErrorOccurred" event is emitted with:
  | errorCode    | COMMAND_EXECUTION_ERROR |
  | errorMessage | Null reference exception |
And the system shows "An error occurred. Please try again."

## Core Components

### Event Bus & Store

- API: publish(event), subscribe(eventType, handler), unsubscribe(eventType, handler), replayEvents(filter)
- Event Store: Append-only log (in-memory or persistent) with metadata for auditability and state recovery.

### Plugin Loader & API

- Responsibilities: Scan designated folders (e.g., for packages named "flexicli-plugin-*"), auto-detect tech stack, and initialize plugins via init(eventBus, config).
- Example Plugin Interface:

```javascript
module.exports = {
  pluginId: "node-test-runner",
  version: "1.0.0",
  description: "Runs Node.js tests via Jest/Mocha",
  dependencies: ["package.json"],
  init: function(eventBus, config) {
    eventBus.subscribe("CommandReceived", (event) => {
      if (event.payload.command === "run-tests") {
        eventBus.publish({
          event: "TestRunStarted",
          timestamp: new Date().toISOString(),
          payload: {
            suite: "unit",
            mode: event.payload.args.includes("--verbose") ? "verbose" : "normal",
            startTime: new Date().toISOString()
          }
        });
        // Execute tests...
        eventBus.publish({
          event: "TestRunCompleted",
          timestamp: new Date().toISOString(),
          payload: {
            suite: "unit",
            result: { passed: 42, failed: 0, skipped: 0 },
            duration: 1234,
            errors: []
          }
        });
      }
    });
  }
};
```

- Plugin API Helpers: registerCommand, emitEvent, getConfig, log.

### Command Parser & Core Engine

- Command Parser: Splits user input into command and arguments, validates against a registry, and emits CommandReceived events.
- Core Flow:
  1. Load config from ".flexicli.json".
  2. Scan project directory and emit TechStackDetected.
  3. Load plugins via the Plugin Loader.
  4. Enter REPL or parse CLI arguments.
  5. For each command, emit CommandReceived and dispatch to handlers.
  6. Listen for events (e.g., TestRunStarted, TestRunCompleted) to update output.

### Configuration & Deployment

- Sample .flexicli.json:

```json
{
  "plugins": {
    "enabled": ["node-test-runner", "python-env-manager"],
    "disabled": []
  },
  "environments": {
    "development": { "VAR1": "devValue", "VAR2": "devSetting" },
    "production": { "VAR1": "prodValue", "VAR2": "prodSetting" }
  },
  "globalSettings": {
    "loggingLevel": "info",
    "commandTimeout": 30000
  }
}
```

- Use a cross‑platform language (Node.js, Go, or Rust), package as a binary or npm package, and integrate BDD tests in CI/CD.

## Additional Notes

- The event store supports replay for debugging and state recovery.
- Extensible via the Plugin API.
- Log all critical events; integration with local/remote loggers is optional.
- Auto‑generate documentation from event schemas and the Plugin API.
- Integrate Gherkin tests into CI pipelines for continuous validation.

## Conclusion

This specification leverages event modeling and BDD to define FlexiCLI’s behavior precisely. With detailed event schemas, comprehensive Gherkin scenarios, and clear documentation of core components, an AI builder can autonomously construct a robust, flexible, cross‑platform CLI application.

## Continuous Integration (CI)

To ensure the quality and reliability of the FlexiCLI project, we have integrated Continuous Integration (CI) using GitHub Actions. The CI pipeline is defined in the `.github/workflows/ci.yml` file.

### CI Status Badge

The current status of the CI pipeline can be seen with the following badge:

![CI](https://github.com/robmillersoftware/flux-cli/actions/workflows/ci.yml/badge.svg)

### Running CI

The CI pipeline is triggered automatically on every push to the `main` branch and on every pull request targeting the `main` branch. The pipeline includes the following steps:

1. **Checkout code**: Retrieves the latest code from the repository.
2. **Set up Go**: Configures the Go environment.
3. **Install dependencies**: Installs the required Go modules.
4. **Build**: Compiles the project.
5. **Run tests**: Executes the unit tests.
6. **Run BDD tests**: Executes the BDD tests using Gherkin scenarios.

To manually trigger the CI pipeline, you can push changes to the `main` branch or create a pull request targeting the `main` branch.
