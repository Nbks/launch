# Launch 🚀

A simple CLI tool to manage and open development projects with predefined tool profiles.

## Features

- Manage multiple development projects
- Define profiles with different tool configurations
- Launch multiple tools simultaneously
- Support for GUI, terminal, and background applications
- Window arrangement support (Linux)
- Environment variable support per profile

## Installation

```bash
go install ./cmd/launch
```

On Linux, you'll need `wmctrl` for window arrangement:

```bash
# Debian/Ubuntu
sudo apt install wmctrl

# Fedora
sudo dnf install wmctrl
```

## Quick Start

### 1. Initialize config

```bash
launch init
```

This creates a config file at the appropriate location:
- **Windows**: `%APPDATA%\launch\config.yml`
- **Linux**: `~/.config/launch/config.yml`

### 2. Add a project

```bash
launch add myproject /path/to/your/project
```

### 3. Open a project

```bash
launch open myproject dev
```

This opens the project with the default "dev" profile.

## Usage

```
launch <command>

Commands:
  init           Initialize config file
  list           List all projects
  add <name> <path>    Add a new project
  open <project> <profile>  Open project with profile
  -v             Enable verbose output
```

## Configuration

Edit your config file to customize projects and profiles:

**Windows**: `%APPDATA%\launch\config.yml`
**Linux**: `~/.config/launch/config.yml`

```yaml
projects:
  myproject:
    path: /home/user/code/myproject
    profiles:
      dev:
        tools:
          - name: vscode
            path: code
            args: ["."]
      full:
        env:
          NODE_ENV: development
        tools:
          - name: vscode
            path: code
            args: ["."]
          - name: terminal
            path: gnome-terminal
            args: ["--"]
            type: terminal
          - name: browser
            path: firefox
            args: ["http://localhost:3000"]
            type: gui
```

### Tool Types

- `gui` - GUI application (default, opens in foreground)
- `terminal` - Terminal application (opens new terminal window)
- `background` - Background process (runs without visible window)

### Layout (Linux only)

Window positions for multi-monitor setups:
- `fullscreen`, `top-half`, `bottom-half`
- `left-half`, `right-half`

### Display

Specify which display to open on:
- `display: "1"` - Display 1
- `display: "2"` - Display 2

## License

MIT License - See LICENSE file
