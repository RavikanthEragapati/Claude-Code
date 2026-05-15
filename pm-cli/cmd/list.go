package main

import (
	"fmt"
	"time"

	"example.com/pm-cli/pkg/storage"
	"github.com/spf13/cobra"
)

var listID string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for the current sprint or show a specific story",
	RunE:  runList,
}

func init() {
	listCmd.Flags().StringVar(&listID, "id", "", "Filter by story ID")
}

func runList(cmd *cobra.Command, args []string) error {
	// Load configuration to get current sprint number
	cfg, err := storage.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	sprint := cfg.SprintNumber
	if sprint == 0 {
		return fmt.Errorf("sprint number not set in config")
	}

	// Load tasks for the sprint
	tasks, err := storage.LoadTasks(sprint)
	if err != nil {
		return fmt.Errorf("could not load tasks for sprint %d: %w", sprint, err)
	}

	// If an ID flag is provided, find that story
	if listID != "" {
		for _, t := range tasks {
			if t.StoryID == listID {
				fmt.Printf("%-30s %-12s %s\n", "CreatedAt", "StoryID", "Title")
				fmt.Printf("%-30s %-12s %s\n", t.CreatedAt.Format(time.RFC3339), t.StoryID, t.Title)
				if len(t.Notes) > 0 {
					fmt.Println("\nNotes:")
					for _, n := range t.Notes {
						fmt.Printf("- %-30s %s\n", n.CreatedAt.Format(time.RFC3339), n.Text)
					}
				}
				return nil
			}
		}
		return fmt.Errorf("story ID %s not found in sprint %d", listID, sprint)
	}

	// No ID filter – list all tasks
	if len(tasks) == 0 {
		fmt.Printf("No tasks found for sprint %d.\n", sprint)
		return nil
	}
	fmt.Printf("Tasks for sprint %d:\n", sprint)
	// Table header
	fmt.Printf("%-30s %-12s %s\n", "CreatedAt", "StoryID", "Title")
	// Table rows
	for _, t := range tasks {
		fmt.Printf("%-30s %-12s %s\n", t.CreatedAt.Format(time.RFC3339), t.StoryID, t.Title)
	}
	return nil
}
