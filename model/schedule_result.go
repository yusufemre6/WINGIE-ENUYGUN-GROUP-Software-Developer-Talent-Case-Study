package model

// TaskSchedule holds the computed timing for a single task.
//
// EarliestFinish = EarliestStart + Duration
type TaskSchedule struct {
	TaskID         string
	EarliestStart  int
	EarliestFinish int
}

// ScheduleResult contains the full output of the scheduling algorithm.
type ScheduleResult struct {
	JobName           string
	Workers           int             // number of workers used
	MinCompletionTime int
	TaskSchedules     []TaskSchedule  // sorted by start time
	ExecutionOrder    []string        // task IDs in order they were started
	CriticalPath      []string        // longest path (only when workers >= task count)
}
