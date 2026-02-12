// Package input handles reading Job definitions from external sources.
package input

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"wingie_case/model"
)

// JobInput holds the result of reading: a job and the number of workers.
type JobInput struct {
	Job     *model.Job
	Workers int
}

// Reader is the interface for reading a Job and worker count from any source.
type Reader interface {
	ReadJob() (*JobInput, error)
}

// CLIReader reads job definitions interactively from a terminal.
// It reads from an io.Reader (e.g. os.Stdin).
type CLIReader struct {
	scanner *bufio.Scanner
}

// NewCLIReader creates a CLIReader backed by the given reader (e.g. os.Stdin).
func NewCLIReader(r io.Reader) *CLIReader {
	return &CLIReader{
		scanner: bufio.NewScanner(r),
	}
}

// ReadJob prompts for job name, task count, each task's data, and number of workers.
func (c *CLIReader) ReadJob() (*JobInput, error) {
	jobName, err := c.promptString("Enter job name (e.g. J)")
	if err != nil {
		return nil, fmt.Errorf("could not read job name: %w", err)
	}

	job := model.NewJob(jobName)

	taskCount, err := c.promptInt("How many tasks?")
	if err != nil {
		return nil, fmt.Errorf("could not read task count: %w", err)
	}
	if taskCount <= 0 {
		return nil, fmt.Errorf("task count must be positive, got %d", taskCount)
	}

	for i := 0; i < taskCount; i++ {
		task, err := c.readTask(i + 1)
		if err != nil {
			return nil, fmt.Errorf("task %d: %w", i+1, err)
		}

		if err := job.AddTask(task); err != nil {
			return nil, err
		}
	}

	workers, err := c.promptInt("How many workers?")
	if err != nil {
		return nil, fmt.Errorf("could not read worker count: %w", err)
	}
	if workers <= 0 {
		return nil, fmt.Errorf("worker count must be positive, got %d", workers)
	}

	return &JobInput{Job: job, Workers: workers}, nil
}

// readTask reads a single task definition from the user.
func (c *CLIReader) readTask(index int) (*model.Task, error) {
	fmt.Printf("\n--- Task %d ---\n", index)

	id, err := c.promptString("Task ID (e.g. A)")
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	duration, err := c.promptInt(fmt.Sprintf("Duration for task '%s' (positive integer)", id))
	if err != nil {
		return nil, err
	}

	depsStr, err := c.promptString(
		fmt.Sprintf("Dependencies for task '%s' (comma-separated, or leave empty)", id))
	if err != nil {
		return nil, err
	}

	deps := parseDependencies(depsStr, id)

	return model.NewTask(id, duration, deps)
}

// parseDependencies splits a comma-separated string into dependency IDs.
// It filters out blanks, duplicates, and self-references.
func parseDependencies(input string, selfID string) []string {
	if input == "" {
		return nil
	}

	var deps []string
	seen := make(map[string]bool)

	for _, p := range strings.Split(input, ",") {
		dep := strings.TrimSpace(p)
		if dep == "" || dep == selfID || seen[dep] {
			continue
		}
		seen[dep] = true
		deps = append(deps, dep)
	}

	return deps
}

func (c *CLIReader) promptString(message string) (string, error) {
	fmt.Printf("%s: ", message)
	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return "", fmt.Errorf("read error: %w", err)
		}
		return "", fmt.Errorf("unexpected end of input (EOF)")
	}
	return strings.TrimSpace(c.scanner.Text()), nil
}

func (c *CLIReader) promptInt(message string) (int, error) {
	str, err := c.promptString(message)
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("invalid number: '%s' (integer expected)", str)
	}
	return val, nil
}
