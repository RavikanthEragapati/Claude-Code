# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

- **Build**: `go build -o pm ./cmd`
- **Run**: `./pm <subcommand> <flags>`
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

- **Configuration** is stored in `~/.pm/config.json`.
- **Sprint data** lives under `~/.pm/data/<sprint-number>/task.json` as a JSON array of task entries.
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

## Typical Workflow

1. Run `./pm init ...` to set up your global config.
2. Run `./pm add ...` to record stories or notes.
3. stories and notes are persisted to `~/.pm/data/<sprint-number>/task.json`.
4. To review stored data, inspect the JSON files directly or extend the tool to provide a `list` command.
