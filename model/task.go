// Package model defines the core domain types for the job scheduling system.
package model

import "fmt"

// Task represents a single unit of work within a Job.
type Task struct {
	ID           string
	Duration     int
	Dependencies []string
}

// NewTask creates a Task with the given parameters.
// Returns an error if id is empty or duration is not positive.
func NewTask(id string, duration int, dependencies []string) (*Task, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if duration <= 0 {
		return nil, fmt.Errorf("task '%s': duration must be positive, got %d", id, duration)
	}

	deps := dependencies
	if deps == nil {
		deps = []string{}
	}

	return &Task{
		ID:           id,
		Duration:     duration,
		Dependencies: deps,
	}, nil
}

// HasDependencies returns true if the task depends on other tasks.
func (t *Task) HasDependencies() bool {
	return len(t.Dependencies) > 0
}

// DependsOn checks whether the task depends on the given task ID.
func (t *Task) DependsOn(taskID string) bool {
	for _, dep := range t.Dependencies {
		if dep == taskID {
			return true
		}
	}
	return false
}
