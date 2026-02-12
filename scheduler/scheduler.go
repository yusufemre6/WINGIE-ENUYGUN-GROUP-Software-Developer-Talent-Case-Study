// Package scheduler computes the minimum completion time of a Job
// using the Critical Path Method (CPM).
//
// Algorithm:
//  1. Topological sort (Kahn's algorithm)
//  2. Forward pass — compute EST and EFT for every task
//  3. Minimum completion time = max(EFT)
//  4. Backward trace — extract the critical path
//
// Time complexity: O(V + E)
package scheduler

import (
	"fmt"
	"sort"

	"wingie_case/model"
)

// Scheduler is the interface for job scheduling strategies.
type Scheduler interface {
	Schedule(job *model.Job) (*model.ScheduleResult, error)
}

// CriticalPathScheduler implements the CPM algorithm.
//
// Assumption: all tasks whose dependencies are satisfied may run
// in parallel (unlimited workers). Each individual task uses one worker.
type CriticalPathScheduler struct{}

func NewCriticalPathScheduler() *CriticalPathScheduler {
	return &CriticalPathScheduler{}
}

// Schedule computes EST/EFT for every task and returns the schedule.
func (s *CriticalPathScheduler) Schedule(job *model.Job) (*model.ScheduleResult, error) {
	order, err := s.topologicalOrder(job)
	if err != nil {
		return nil, err
	}

	est := make(map[string]int, job.TaskCount())
	eft := make(map[string]int, job.TaskCount())

	// Forward pass: compute Earliest Start / Earliest Finish
	//   EST(t) = max( EFT(dep) ) for each dependency
	//   EFT(t) = EST(t) + duration(t)
	for _, id := range order {
		task := job.Tasks[id]

		if !task.HasDependencies() {
			est[id] = 0
		} else {
			maxFinish := 0
			for _, depID := range task.Dependencies {
				if eft[depID] > maxFinish {
					maxFinish = eft[depID]
				}
			}
			est[id] = maxFinish
		}

		eft[id] = est[id] + task.Duration
	}

	// The job finishes when the last task finishes
	minCompletion := 0
	for _, f := range eft {
		if f > minCompletion {
			minCompletion = f
		}
	}

	schedules := s.buildSortedSchedules(order, est, eft)

	executionOrder := make([]string, 0, len(schedules))
	for _, ts := range schedules {
		executionOrder = append(executionOrder, ts.TaskID)
	}

	criticalPath := s.findCriticalPath(job, est, eft, minCompletion)

	return &model.ScheduleResult{
		JobName:           job.Name,
		MinCompletionTime: minCompletion,
		TaskSchedules:     schedules,
		ExecutionOrder:    executionOrder,
		CriticalPath:      criticalPath,
	}, nil
}

// findCriticalPath traces back from the latest-finishing task to find
// the chain of tasks that determines the total duration.
func (s *CriticalPathScheduler) findCriticalPath(
	job *model.Job,
	est, eft map[string]int,
	minCompletion int,
) []string {
	// Find the task that finishes last
	var endTaskID string
	for id, f := range eft {
		if f == minCompletion {
			endTaskID = id
			break
		}
	}

	path := []string{endTaskID}
	currentID := endTaskID

	for {
		task := job.Tasks[currentID]
		if !task.HasDependencies() {
			break
		}

		// Pick the dependency whose EFT equals the current task's EST
		found := false
		for _, depID := range task.Dependencies {
			if eft[depID] == est[currentID] {
				path = append([]string{depID}, path...)
				currentID = depID
				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	return path
}

// buildSortedSchedules creates TaskSchedule entries sorted by EST.
func (s *CriticalPathScheduler) buildSortedSchedules(
	order []string,
	est, eft map[string]int,
) []model.TaskSchedule {
	schedules := make([]model.TaskSchedule, 0, len(order))
	for _, id := range order {
		schedules = append(schedules, model.TaskSchedule{
			TaskID:         id,
			EarliestStart:  est[id],
			EarliestFinish: eft[id],
		})
	}

	sort.Slice(schedules, func(i, j int) bool {
		if schedules[i].EarliestStart == schedules[j].EarliestStart {
			return schedules[i].TaskID < schedules[j].TaskID
		}
		return schedules[i].EarliestStart < schedules[j].EarliestStart
	})

	return schedules
}

// topologicalOrder returns tasks in topological order using Kahn's algorithm.
// The queue is sorted alphabetically at each step for deterministic output.
func (s *CriticalPathScheduler) topologicalOrder(job *model.Job) ([]string, error) {
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
	sort.Strings(queue)

	order := make([]string, 0, job.TaskCount())
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		neighbors := adjacency[current]
		sort.Strings(neighbors)

		for _, neighbor := range neighbors {
			indegree[neighbor]--
			if indegree[neighbor] == 0 {
				queue = append(queue, neighbor)
				sort.Strings(queue)
			}
		}
	}

	if len(order) != job.TaskCount() {
		return nil, fmt.Errorf("topological sort failed: graph contains a cycle")
	}

	return order, nil
}
