package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"example.com/pm-cli/pkg/storage"
)

var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Project management CLI for tracking JIRA stories",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// -- init command --
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize global configuration for pm CLI",
		RunE:  runInit,
	}
	initFlags := initCmd.Flags()
	initFlags.StringP("name", "n", "", "Your name")
	initFlags.StringP("team", "t", "", "Team name")
	initFlags.StringP("framework", "f", "", "Agile framework (scrum|kanban)")
	initFlags.IntP("sprint-number", "", 1, "Current sprint number (positive integer)")
	initFlags.IntP("sprint-length", "", 2, "Sprint length in weeks (positive integer)")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("team")
	_ = initCmd.MarkFlagRequired("framework")
	_ = initCmd.MarkFlagRequired("sprint-number")
	_ = initCmd.MarkFlagRequired("sprint-length")
	rootCmd.AddCommand(initCmd)

	// -- add command --
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a story or a note to the current sprint",
		RunE:  runAdd,
	}
	addFlags := addCmd.Flags()
	addFlags.StringP("id", "i", "", "Story ID (required)")
	addFlags.StringP("title", "t", "", "Title of the story (optional)")
	addFlags.StringP("note", "n", "", "Note text (optional)")
	_ = addCmd.MarkFlagRequired("id")
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
}

// runInit validates flags and writes config file.
func runInit(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	team, _ := cmd.Flags().GetString("team")
	framework, _ := cmd.Flags().GetString("framework")
	sprintNum, _ := cmd.Flags().GetInt("sprint-number")
	sprintLen, _ := cmd.Flags().GetInt("sprint-length")

	if name == "" || team == "" || framework == "" {
		return fmt.Errorf("all of --name, --team, and --framework must be provided")
	}
	if sprintNum <= 0 {
		return fmt.Errorf("--sprint-number must be a positive integer")
	}
	if sprintLen <= 0 {
		return fmt.Errorf("--sprint-length must be a positive integer")
	}
	config := storage.Config{
		Name:           name,
		Team:           team,
		Framework:      framework,
		SprintNumber:   sprintNum,
		SprintLength:   sprintLen,
	}
	if err := storage.SaveConfig(&config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	fmt.Printf("Configuration saved to %s\n", storage.ConfigFile())
	return nil
}

// runAdd handles adding a story or a note.
// It expects the global add flags (--id, --title, --note) to be set via cobra flag parsing.
func runAdd(cmd *cobra.Command, args []string) error {
	id, _ := cmd.Flags().GetString("id")
	title, _ := cmd.Flags().GetString("title")
	note, _ := cmd.Flags().GetString("note")

	if id == "" {
		return fmt.Errorf("flag --id is required")
	}

	// At least one of title or note must be provided
	if strings.TrimSpace(title) == "" && strings.TrimSpace(note) == "" {
		return fmt.Errorf("either --title or --note must be provided")
	}

	// Load config to get sprint number
	cfg, err := storage.LoadConfig()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}
	sprintNumber := cfg.SprintNumber
	if sprintNumber == 0 {
		return fmt.Errorf("sprint-number not set in config")
	}

	if strings.TrimSpace(title) != "" {
		// Add or update story title
		if err := storage.AddStory(sprintNumber, id, title); err != nil {
			return fmt.Errorf("failed to add story: %w", err)
		}
		fmt.Printf("Story %s title set to \"%s\"\n", id, title)
	}
	if strings.TrimSpace(note) != "" {
		// Append note to existing story
		if err := storage.AddNote(sprintNumber, id, note); err != nil {
			return fmt.Errorf("failed to add note: %w", err)
		}
		fmt.Printf("Note added to story %s\n", id)
	}
	return nil
}