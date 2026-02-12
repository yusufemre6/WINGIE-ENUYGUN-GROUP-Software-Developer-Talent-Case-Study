package model

import "fmt"

// Job holds a collection of Tasks organized as a DAG (Directed Acyclic Graph).
type Job struct {
	Name  string
	Tasks map[string]*Task
}

// NewJob creates a new Job. Falls back to "Job" when name is empty.
func NewJob(name string) *Job {
	if name == "" {
		name = "Job"
	}
	return &Job{
		Name:  name,
		Tasks: make(map[string]*Task),
	}
}

// AddTask inserts a task into the job. Returns an error on duplicate IDs.
func (j *Job) AddTask(task *Task) error {
	if task == nil {
		return fmt.Errorf("cannot add nil task")
	}
	if _, exists := j.Tasks[task.ID]; exists {
		return fmt.Errorf("duplicate task ID: '%s'", task.ID)
	}
	j.Tasks[task.ID] = task
	return nil
}

// GetTask returns the task with the given ID, or nil if not found.
func (j *Job) GetTask(id string) (*Task, bool) {
	task, ok := j.Tasks[id]
	return task, ok
}

// TaskCount returns the number of tasks in the job.
func (j *Job) TaskCount() int {
	return len(j.Tasks)
}

// IndependentTasks returns all tasks that have no dependencies.
func (j *Job) IndependentTasks() []*Task {
	var result []*Task
	for _, task := range j.Tasks {
		if !task.HasDependencies() {
			result = append(result, task)
		}
	}
	return result
}
