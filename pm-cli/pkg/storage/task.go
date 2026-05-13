package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var storyIDRegexp = regexp.MustCompile(`^[A-Z]+-\d+$`)

type TaskEntry struct {
	StoryID string   `json:"storyID"`
	Title   string   `json:"title"`
	Notes   []string `json:"notes,omitempty"`
}

func BaseDir() string {
	return filepath.Join("~", ".pm", "data")
}

func SprintDir(sprintNumber int) string {
	return filepath.Join(BaseDir(), fmt.Sprintf("%d", sprintNumber))
}

func TaskFile(sprintNumber int) string {
	return filepath.Join(SprintDir(sprintNumber), "task.json")
}

func ensureSprintDir(sprintNumber int) error {
	return os.MkdirAll(SprintDir(sprintNumber), 0755)
}

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

func AddStory(sprintNumber int, storyID, title string) error {
	if strings.TrimSpace(storyID) == "" || strings.TrimSpace(title) == "" {
		return fmt.Errorf("storyID and title must be non‑empty")
	}
	if !storyIDRegexp.MatchString(storyID) {
		return fmt.Errorf("invalid storyID format: %s", storyID)
	}
	tasks, err := LoadTasks(sprintNumber)
	if err != nil {
		return err
	}
	tasks = append(tasks, TaskEntry{StoryID: storyID, Title: title})
	return SaveTasks(sprintNumber, tasks)
}

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