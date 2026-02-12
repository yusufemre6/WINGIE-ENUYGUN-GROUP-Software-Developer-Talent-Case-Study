// Package validator checks that a Job forms a valid DAG before scheduling.
package validator

import (
	"fmt"

	"wingie_case/model"
)

// ValidationError represents an invalid field in the job definition.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s", e.Field, e.Message)
	}
	return e.Message
}

// CycleError is returned when the dependency graph contains a cycle.
type CycleError struct {
	Message string
}

func (e *CycleError) Error() string {
	return e.Message
}

// Validator defines the contract for job validation.
type Validator interface {
	Validate(job *model.Job) error
}

// GraphValidator validates the dependency graph of a job.
// It checks for empty jobs, invalid durations, undefined or self
// dependencies, and cycles.
type GraphValidator struct{}

func NewGraphValidator() *GraphValidator {
	return &GraphValidator{}
}

// Validate runs all checks and returns the first error encountered.
func (v *GraphValidator) Validate(job *model.Job) error {
	if job.TaskCount() == 0 {
		return &ValidationError{
			Field:   "job.tasks",
			Message: "job has no tasks",
		}
	}

	for id, task := range job.Tasks {
		if task.Duration <= 0 {
			return &ValidationError{
				Field:   fmt.Sprintf("task.%s.duration", id),
				Message: fmt.Sprintf("duration must be positive, got %d", task.Duration),
			}
		}

		for _, depID := range task.Dependencies {
			if depID == id {
				return &ValidationError{
					Field:   fmt.Sprintf("task.%s.dependencies", id),
					Message: fmt.Sprintf("task '%s' cannot depend on itself", id),
				}
			}
			if _, exists := job.Tasks[depID]; !exists {
				return &ValidationError{
					Field:   fmt.Sprintf("task.%s.dependencies", id),
					Message: fmt.Sprintf("task '%s' depends on undefined task '%s'", id, depID),
				}
			}
		}
	}

	return v.detectCycle(job)
}

// detectCycle uses Kahn's algorithm (BFS topological sort) to detect cycles.
// If not all nodes are visited, the graph contains a cycle.
// Time complexity: O(V + E)
func (v *GraphValidator) detectCycle(job *model.Job) error {
	indegree := make(map[string]int, job.TaskCount())
	adjacency := make(map[string][]string, job.TaskCount())

	for id := range job.Tasks {
		indegree[id] = 0
	}

	for id, task := range job.Tasks {
		for _, depID := range task.Dependencies {
			adjacency[depID] = append(adjacency[depID], id)
			indegree[id]++
		}
	}

	queue := make([]string, 0, job.TaskCount())
	for id, deg := range indegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	visitedCount := 0
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visitedCount++

		for _, neighbor := range adjacency[current] {
			indegree[neighbor]--
			if indegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if visitedCount != job.TaskCount() {
		return &CycleError{
			Message: "dependency graph contains a cycle, job cannot be completed",
		}
	}

	return nil
}
