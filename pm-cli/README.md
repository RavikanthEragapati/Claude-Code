# Context
We are building a CLI `pm` for project management. The tool will allow engineers to track JIRA stories and later consolidate work. First phase includes `ps init` to store global configuration (~/.pm/.config) and `ps add` to add stories and notes, persisting sprint data under ~/.pm/data/<sprint-number>/task.json.

# Approach
1. Implement in Go for a compiled binary, single binary distribution.
2. Use Cobra for CLI parsing (powerful subcommands).
3. Persistence layer:
   - Config stored at `$HOME/.pm/.config/config.json`.
   - Sprint data at `~HOME~/.pm/data/<sprint-number>/task.json` (JSON array).
4. Simplified `ps add` command:
   - Flags: `--id` (required, Unique per sprint), `--title` (optional when adding a new story), `--note` (optional additional note).
   - The presence of `--title` creates or updates a stories title; `--note` appends a note to an existing story.
5. Data structure:
   - Config: `{name, team, framework, sprintNumber, sprintLength}`.
   - Task entry: `[{id, title, notes[]},{id, title, notes[]}]` stored in sprint folder.
6. Validation: ensure required flags provided, sprint number positive integer, story ID matches JIRA pattern, etc.
7. Error handling: create directories if missing, atomic writes.

# Critical Files
- `cmd/main.go`: entry point, root command.
- `cmd/init.go`: implements `ps init`.
- `cmd/add.go`: implements simplified `ps add` command logic.
- `pkg/storage/config.go`: config file I/O.
- `pkg/storage/task.go`: task JSON read/write, sprint folder handling.

# Tasks
- [x] Create Go module structure (`go.mod`, `main.go`).
- [x] Implement Cobra root command and subcommands (`init`, updated `add`).
- [x] Build config storage (`pkg/storage/config.go`).
- [x] Build task storage (`pkg/storage/task.go`).
- [x] Implement `ps init` command logic and validation.
- [x] Implement simplified `ps add` command that removes `--mode` flag, makes `--id` always required, and appends notes if story exists.
- [ ] BugFix: Update data structure of Task entry make is an array of objects
- [ ] BugFix: Implement validation on id field in task.json file no to objects with have same id. 
- [ ] Add JSON schema validation for inputs.
- [ ] Write unit tests for storage and command handlers.
- [ ] Build binary and perform end‑to‑end verification.

# Verification
- Run `pm init ...` and confirm `~/.pm/.config/config.json` exists with correct fields.
- Run `pm add --id="PROJ-123" --title="Fix login"` and verify entry appended to `$HOME$/.pm/data/1/task.json`.
- Run `pm add --id="PROJ-123" --note="Review with security"` and verify note added to the array inside an existing entry.
- Run `pm add --id="PROJ-123" --title="Implement UI"` and verify title is updated if story with id already exist or add it .
- Check JSON syntax validity.

# Next Steps
- Update Code to implement BugFix listed above.
- Write unit tests for the new add logic.
- Add proper logging and flag completion.