# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Project Overview & Tech Stack
- **Language/Run-time**: Go 1.20
- **Primary Framework**: Cobra (github.com/spf13/cobra) for CLI structuring
- **Package Layout**: 
  - `cmd/` contains the root command and subcommands (`init`, `add`)
  - `pkg/storage/` handles configuration persistence and story/note storage
- **Dependency Management**: Modules (`go.mod`) with indirect dependencies on `github.com/inconshreveable/mousetrap` and `github.com/spf13/pflag`
- **Styling/Approach**: Standard Go formatting (`gofmt`), no CSS framework; command-line interface only
- **State**: CLI maintains runtime config via `storage.Config` struct persisted to a file
- **Key Components**: 
  - `main.go` defines root command and bootstraps CLI
  - `storage` package provides `SaveConfig`, `LoadConfig`, `AddStory`, `AddNote` methods
  - Command flags enforce required inputs (e.g., `--name`, `--team`, `--framework`)


## Development Commands

- **Build**: `go build -o pm ./cmd` (produces the `pm` executable)
- **Run CLI**: `./pm [command] <flags>`
  - Example: `./pm init --name="John Doe" --team="AgileTeam" --framework="scrum" --sprint-number=1 --sprint-length=2`
  - Example: `./pm add --id="PROJ-123" --title="Fix login" --note="Feature tested and story complete"` 
- **Test**: The project does not have a test suite yet; when tests are added, use `go test ./...`
- **Lint**: Use `golangci-lint run ./...` if installed

## Command Overview

| Command | Purpose | Important Flags | Example |
|---------|---------|----------------|---------|
| `init` | Initialize global configuration | `--name`, `--team`, `--framework`, `--sprint-number`, `--sprint-length` | `pm init --name="Jane" --team="Design" --framework="kanban" --sprint-number=2 --sprint-length=1` |
| `add` | Add a story or a note to the current sprint | `--id`, `--title`, `--note` | `pm add --id="PROJ-456" --title="Implement auth"` or `pm add --id="PROJ-456" --note="Discuss with security"` |

### Init Flags
- `--name` – Your name
- `--team` – Team name
- `--framework` – Either `scrum` or `kanban`
- `--sprint-number` – Current sprint number (positive integer)
- `--sprint-length` – Length of a sprint in weeks (positive integer)

### Add Flags
- `--id` – JIRA story identifier (e.g., `PROJ-123`) – required for all modes
- `--title` – Title of the story
- `--note` – Optional note

## Data Storage

- **Configuration** is stored in `$HOME/.pm/config.json`.
- **Sprint data** lives under `$HOME/.pm/data/<sprint-number>/task.json` as a JSON array of task entries.
  - Each task entry contains `storyID`, `title`, and optional array of `notes`.

## Validation Rules

- Story IDs must match the pattern `[A-Z]+-\d+` (e.g., `PROJ-123`).
- All required flags must be provided; missing required fields return descriptive errors.
- Positive integers are enforced for sprint number and length.

## Important Packages

- `github.com/spf13/cobra` – CLI framework.
- `github.com/spf13/pflag` – Flag parsing.
- Local packages:
  - `pkg/storage/config.go` – Handles config file I/O.
  - `pkg/storage/task.go` – Handles sprint task read/write and validation.

## How to Extend

1. **Add New Subcommands** – Create a new `*cobra.Command` in `cmd/main.go`, add it to the root command, define flags, and implement logic.
2. **Persist Additional Data** – Use `pkg/storage` utilities to read/write JSON files.
3. **Validate Input** – Follow existing validation patterns; add new regexes or checks as needed.
4. **Write Tests** – Place test files alongside the package they test; run `go test ./...` to execute.

# Detected Code Patterns & Architecture Rules
- **Command Definition**: All CLI commands are Cobra commands registered with `Use`, `Short`, and `RunE` fields
- **Flag Handling**: Flags are parsed via Cobra flag sets; required flags are enforced with `MarkFlagRequired`
- **Error Handling**: Functions return `error` values; errors are printed to stdout and exit with status 1
- **Configuration Persistence**: JSON or similar format used in `storage.ConfigFile()`; config struct includes `Name`, `Team`, `Framework`, `SprintNumber`, `SprintLength`
- **Export Style**: Commands are exported via `init()` functions that add them to the root command
- **Validation**: Input validation occurs before proceeding (e.g., required fields, positive integers)

## Typical Workflow

1. Run `./pm init ...` to set up your global config.
2. Run `./pm add ...` to record stories or notes.
3. stories and notes are persisted to `$HOME$/.pm/data/<sprint-number>/task.json`.
4. To review stored data, inspect the JSON files directly or extend the tool to provide a `list` command.

# Tasks
- [x] Create Go module structure (`go.mod`, `main.go`).
- [x] Implement Cobra root command and subcommands (`init`, updated `add`).
- [x] Build config storage (`pkg/storage/config.go`).
- [x] Build task storage (`pkg/storage/task.go`).
- [x] Implement `ps init` command logic and validation.
- [x] Implement `ps add` command logic and validation.
- [ ] Implement `ps list` command logic and validate.
- [ ] Implement `ps brag` command logic and validate.
- [ ] Add JSON schema validation for inputs.
- [ ] Write unit tests for storage and command handlers.
- [ ] Build binary and perform end‑to‑end verification.