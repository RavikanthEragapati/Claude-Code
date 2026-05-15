# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview & Tech Stack
- **Language / Runtime**: Go 1.20
- **CLI framework**: Cobra (`github.com/spf13/cobra`)
- **Package layout**:
  - `cmd/` – root command (`main.go`) and sub‑commands (`init`, `add`, `list`)
  - `pkg/storage/` – persistence layer (`config.go`, `task.go`)
- **Data storage**: JSON files under `$HOME/.pm/`
  - Config: `$HOME/.pm/config.json`
  - Sprint data: `$HOME/.pm/data/<sprint-number>/task.json`

## Development Commands
- **Build**: `go build -o pm ./cmd` – produces the `pm` executable in the repository root.
- **Run the CLI**: `./pm <command> [flags]`
- **Lint**: `golangci-lint run ./...` (if the linter is installed).
- **Test**: `go test ./...` – runs all tests.
- **Run a single test**: `go test ./... -run TestFunctionName`

## Command Overview
| Command | Purpose | Important Flags |
|---------|---------|-----------------|
| `init` | Initialise global configuration | `--name`, `--team`, `--framework`, `--sprint-number`, `--sprint-length` |
| `add`  | Add a story or a note to the current sprint | `--id`, `--title`, `--note` |
| `list` | List all tasks for the current sprint or show a specific story | `--id` (optional) |

## High‑Level Architecture
- **Root command (`cmd/main.go`)** creates a `cobra.Command` called `pm` and registers sub‑commands via `init()` functions.
- **Sub‑commands** are defined in their own files (e.g. `cmd/init.go`, `cmd/add.go`, `cmd/list.go`). Each command validates its flags, loads the configuration with `storage.LoadConfig()`, and delegates to the storage layer.
- **Storage package (`pkg/storage`)** abstracts all file I/O:
  - `config.go` – `Config` struct, `LoadConfig`, `SaveConfig`, helper paths.
  - `task.go` – `TaskEntry` struct, `LoadTasks`, `SaveTasks`, `AddStory`, `AddNote` and validation helpers.
- **Data flow**: CLI → command handler → storage layer → JSON on disk.
- **Error handling**: Functions return `error`; the CLI prints the error and exits with status 1.

## Extending the Tool
1. **Add a new sub‑command** – create a `*cobra.Command` in `cmd/`, implement its `RunE`, and register it with `rootCmd.AddCommand(newCmd)`.
2. **Persist additional data** – add structs and read/write helpers in `pkg/storage` and use the same JSON‑file pattern.
3. **Validation** – follow existing patterns: use `MarkFlagRequired`, regular expressions, and explicit checks before persisting.
4. **Testing** – place test files alongside the package (e.g. `pkg/storage/config_test.go`) and run `go test ./...`.

## Detected Patterns & Rules
- Commands are defined with `Use`, `Short`, and `RunE` fields and registered in `init()`.
- Flags are added via `cmd.Flags().StringP/IntP` and required flags are enforced with `MarkFlagRequired`.
- Configuration and task data are stored as JSON; directories are created on‑demand.
- Story IDs are validated against the regex `[A-Z]+-\d+`.
- All public functions return `error` for callers to handle.

## Typical Workflow
1. `./pm init --name="Alice" --team="Backend" --framework="scrum" --sprint-number=1 --sprint-length=2`
2. `./pm add --id="PROJ-101" --title="Implement auth"`
3. `./pm add --id="PROJ-101" --note="Reviewed by security"`
4. `./pm list` – shows all stories for the current sprint.
5. `./pm list --id=PROJ-101` – shows details for a single story.

## Tasks (current state)
- [x] Create Go module, root command, and sub‑commands (`init`, `add`).
- [x] Implement config and task storage.
- [x] Add `list` sub‑command.
- [x] Add datetime field to TaskEntry json at root level.
- [x] Add datetime field to Notes[] in TaskEntry json make notes an object.
- [ ] Update the command output to display in table format.
- [ ] Add `brag` sub‑command.
- [ ] Add JSON schema validation for inputs.
- [ ] Write unit tests for storage and command handlers.
- [ ] End‑to‑end verification of the binary.
