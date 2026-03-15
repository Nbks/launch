# AGENTS.md - Development Guidelines for Launch

This document provides guidelines for AI agents working on this codebase.

## Project Overview

**Launch** is a Go CLI tool to manage and open development projects with predefined tool profiles.
- **Module**: `launch`
- **Go Version**: 1.25.0
- **Dependencies**: gopkg.in/yaml.v3

## Build, Lint, and Test Commands

```bash
# Build the CLI binary
go build ./cmd/launch

# Install the CLI to $GOPATH/bin
go install ./cmd/launch

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/config/...

# Run a single test (verbose)
go test -v -run TestFunctionName ./internal/package/...

# Format code
go fmt ./...

# Run go vet
go vet ./...

# Tidy go.mod
go mod tidy

# List all packages
go list ./...
```

## Code Style Guidelines

### Imports

Group imports: 1) Standard library, 2) External packages, 3) Internal packages.

```go
import (
    "flag"
    "fmt"
    "os"
    
    "launch/internal/config"
    "launch/internal/logger"
    
    "gopkg.in/yaml.v3"
)
```

### Naming Conventions

- **Packages**: lowercase, short, no underscores (e.g., `config`, `project`)
- **Exported functions/types**: PascalCase (e.g., `Load()`, `Config`)
- **Unexported functions/variables**: camelCase (e.g., `defaultPath()`)
- **File names**: lowercase with underscores for platform-specific files (`config_windows.go`)

### Structs and Types

```go
type Config struct {
    Projects map[string]Project `yaml:"projects"`
}

type Project struct {
    Path     string             `yaml:"path"`
    Profiles map[string]Profile `yaml:"profiles"`
}
```

### Error Handling

Return errors early, use `fmt.Errorf` with %w for wrapped errors.

```go
func Load() (*Config, error) {
    data, err := os.ReadFile(DefaultPath())
    if err != nil {
        return nil, fmt.Errorf("reading config: %w", err)
    }
    // ...
}
```

### Logging

Use `internal/logger` with structured key-value logging.

```go
logger.Log.Debug("starting tool",
    "name", t.Name,
    "type", t.Type,
)
```

### Platform-Specific Code

Use build tags:

```go
//go:build windows

package config

func DefaultPath() string {
    appData := os.Getenv("APPDATA")
    return filepath.Join(appData, "launch", "config.yml")
}
```

### Concurrency

Use `sync.WaitGroup` and `sync.Mutex`. Always defer `wg.Done()`.

```go
var wg sync.WaitGroup
var mu sync.Mutex
var errs []error

for _, tool := range profile.Tools {
    wg.Add(1)
    go func(t config.Tool) {
        defer wg.Done()
        // ... work ...
        mu.Lock()
        errs = append(errs, err)
        mu.Unlock()
    }(tool)
}
wg.Wait()
```

### Code Organization

- `cmd/launch/main.go` - Entry point
- `internal/` - Private application code
  - `config/` - Configuration handling
  - `project/` - Project management
  - `launcher/` - Tool launching logic
  - `logger/` - Logging utilities
  - `workspace/` - Window arrangement (Linux)

### Testing

Place test files in the same package, name them `*_test.go`.

## Common Patterns

### CLI Flags

```go
verbose := flag.Bool("v", false, "verbose output")
flag.Parse()
args := flag.Args()
```

### File Operations

Use `os.ReadFile` / `os.WriteFile`, `os.MkdirAll`, and appropriate permissions (0755, 0644).

## Working on This Project

1. Run `go mod tidy` after adding dependencies
2. Run `go fmt ./...` before committing
3. Run `go vet ./...` to catch potential issues
4. Test on both Windows and Linux for platform-specific changes