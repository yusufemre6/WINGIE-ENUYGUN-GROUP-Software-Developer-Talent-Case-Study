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
	MinCompletionTime int
	TaskSchedules     []TaskSchedule // sorted by EarliestStart
	ExecutionOrder    []string       // task IDs sorted by start time
	CriticalPath      []string       // longest path through the DAG
}
