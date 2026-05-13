package storage

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

var storyIDRegexp = regexp.MustCompile(`^[0-9]+$`)

type TaskEntry struct {
    StoryID string   `json:"storyID"`
    Title   string   `json:"title"`
    Notes   []string `json:"notes,omitempty"`
}

// BaseDir returns the base directory for storing sprint data.
func BaseDir() string {
    return filepath.Join("~", ".pm", "data")
}

// SprintDir returns the directory for a specific sprint number.
func SprintDir(sprintNumber int) string {
    return filepath.Join(BaseDir(), fmt.Sprintf("%d", sprintNumber))
}

// TaskFile returns the file path for a sprint's task.json.
func TaskFile(sprintNumber int) string {
    return filepath.Join(SprintDir(sprintNumber), "task.json")
}

// ensureSprintDir ensures the sprint directory exists.
func ensureSprintDir(sprintNumber int) error {
    return os.MkdirAll(SprintDir(sprintNumber), 0755)
}

// LoadTasks reads the task.json file for a sprint and unmarshals it into a slice of TaskEntry.
func LoadTasks(sprintNumber int) ([]TaskEntry, error) {
    if err := ensureSprintDir(sprintNumber); err != nil {
        return nil, err
    }
    data, err := os.ReadFile(TaskFile(sprintNumber))
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil
        }
        return nil, err
    }
    var tasks []TaskEntry
    if len(data) > 0 {
        if err := json.Unmarshal(data, &tasks); err != nil {
            return nil, err
        }
    }
    return tasks, nil
}

// SaveTasks marshals the slice of TaskEntry into JSON and writes it to task.json.
func SaveTasks(sprintNumber int, tasks []TaskEntry) error {
    if err := ensureSprintDir(sprintNumber); err != nil {
        return err
    }
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(TaskFile(sprintNumber), data, 0644)
}

// validateIDUniqueness checks that there are no duplicate StoryID values in the provided slice.
func validateIDUniqueness(tasks []TaskEntry) error {
    seen := make(map[string]bool)
    for _, t := range tasks {
        if seen[t.StoryID] {
            return fmt.Errorf("duplicate storyID found: %s", t.StoryID)
        }
        seen[t.StoryID] = true
    }
    return nil
}

// AddStory adds a new story or updates an existing story's title if the storyID already exists.
// It enforces numeric storyID and ensures uniqueness within the sprint.
func AddStory(sprintNumber int, storyID, title string) error {
    if strings.TrimSpace(storyID) == "" || strings.TrimSpace(title) == "" {
        return fmt.Errorf("storyID and title must be non‑empty")
    }
    // Enforce numeric storyID.
    if !storyIDRegexp.MatchString(storyID) {
        return fmt.Errorf("storyID must be numeric")
    }
    // Load existing tasks.
    tasks, err := LoadTasks(sprintNumber)
    if err != nil {
        return err
    }
    // Validate uniqueness of IDs in the current slice.
    if err := validateIDUniqueness(tasks); err != nil {
        return err
    }
    // Check if the storyID already exists.
    for i, t := range tasks {
        if t.StoryID == storyID {
            // Update the existing story's title.
            tasks[i].Title = title
            return SaveTasks(sprintNumber, tasks)
        }
    }
    // If not found, append a new story entry.
    tasks = append(tasks, TaskEntry{StoryID: storyID, Title: title})
    return SaveTasks(sprintNumber, tasks)
}

// AddNote appends a note to an existing story identified by storyID.
func AddNote(sprintNumber int, storyID, note string) error {
    if strings.TrimSpace(storyID) == "" || strings.TrimSpace(note) == "" {
        return fmt.Errorf("storyID and note must be non‑empty")
    }
    tasks, err := LoadTasks(sprintNumber)
    if err != nil {
        return err
    }
    for i := range tasks {
        if tasks[i].StoryID == storyID {
            tasks[i].Notes = append(tasks[i].Notes, note)
            return SaveTasks(sprintNumber, tasks)
        }
    }
    return fmt.Errorf("storyID %s not found in sprint %d", storyID, sprintNumber)
}