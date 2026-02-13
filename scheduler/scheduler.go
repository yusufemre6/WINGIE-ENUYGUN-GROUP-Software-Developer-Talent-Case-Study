// Package scheduler computes a schedule for a job with a fixed number of workers.
// When workers >= number of tasks, the result matches CPM (unlimited parallelism).
// Otherwise a discrete-event simulation assigns tasks to workers as they become free.
package scheduler

import (
	"fmt"
	"sort"

	"wingie_case/model"
)

// Scheduler is the interface for job scheduling with a given number of workers.
type Scheduler interface {
	Schedule(job *model.Job, workers int) (*model.ScheduleResult, error)
}

// WorkerScheduler schedules tasks with a limited number of workers.
type WorkerScheduler struct{}

func NewWorkerScheduler() *WorkerScheduler {
	return &WorkerScheduler{}
}

// Schedule returns a schedule for the job using the given number of workers.
// When workers >= task count, uses CPM (minimum completion time).
// Otherwise simulates time and assigns ready tasks to free workers.
func (s *WorkerScheduler) Schedule(job *model.Job, workers int) (*model.ScheduleResult, error) {
	if workers <= 0 {
		return nil, fmt.Errorf("workers must be positive, got %d", workers)
	}

	// When we have at least as many workers as tasks, unlimited parallelism applies.
	if workers >= job.TaskCount() {
		return s.scheduleUnlimited(job, workers)
	}

	return s.scheduleLimited(job, workers)
}

// scheduleUnlimited runs CPM and sets Workers on the result.
func (s *WorkerScheduler) scheduleUnlimited(job *model.Job, workers int) (*model.ScheduleResult, error) {
	order, err := s.topologicalOrder(job)
	if err != nil {
		return nil, err
	}

	est := make(map[string]int, job.TaskCount())
	eft := make(map[string]int, job.TaskCount())

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
		Workers:           workers,
		MinCompletionTime: minCompletion,
		TaskSchedules:     schedules,
		ExecutionOrder:    executionOrder,
		CriticalPath:      criticalPath,
	}, nil
}

// scheduleLimited runs a discrete-event simulation with a fixed number of workers.
func (s *WorkerScheduler) scheduleLimited(job *model.Job, workers int) (*model.ScheduleResult, error) {
	_, err := s.topologicalOrder(job)
	if err != nil {
		return nil, err
	}

	// reverse[taskID] = tasks that depend on taskID
	reverse := make(map[string][]string, job.TaskCount())
	for id, task := range job.Tasks {
		for _, depID := range task.Dependencies {
			reverse[depID] = append(reverse[depID], id)
		}
	}

	finished := make(map[string]int)
	startTime := make(map[string]int)
	ready := make(map[string]bool)

	for id, task := range job.Tasks {
		if !task.HasDependencies() {
			ready[id] = true
		}
	}

	type slot struct {
		taskID     string
		finishTime int
	}
	running := make([]slot, 0, workers)
	var executionOrder []string
	currentTime := 0

	for {
		// Assign as many ready tasks as we have free workers
		readyList := make([]string, 0, len(ready))
		for id := range ready {
			readyList = append(readyList, id)
		}
		sort.Strings(readyList)

		for len(running) < workers && len(readyList) > 0 {
			id := readyList[0]
			readyList = readyList[1:]
			delete(ready, id)

			task := job.Tasks[id]
			finish := currentTime + task.Duration
			startTime[id] = currentTime
			running = append(running, slot{taskID: id, finishTime: finish})
			executionOrder = append(executionOrder, id)
		}

		if len(running) == 0 {
			break
		}

		// Advance to next completion event
		minFinish := running[0].finishTime
		for _, sl := range running[1:] {
			if sl.finishTime < minFinish {
				minFinish = sl.finishTime
			}
		}
		currentTime = minFinish

		// Complete all tasks that finish at currentTime
		newRunning := running[:0]
		for _, sl := range running {
			if sl.finishTime == currentTime {
				finished[sl.taskID] = currentTime
				for _, nextID := range reverse[sl.taskID] {
					task := job.Tasks[nextID]
					allDone := true
					for _, depID := range task.Dependencies {
						if _, ok := finished[depID]; !ok {
							allDone = false
							break
						}
					}
					if allDone {
						ready[nextID] = true
					}
				}
			} else {
				newRunning = append(newRunning, sl)
			}
		}
		running = newRunning
	}

	// Build TaskSchedules sorted by start time
	schedules := make([]model.TaskSchedule, 0, job.TaskCount())
	for id := range startTime {
		start := startTime[id]
		finish := finished[id]
		schedules = append(schedules, model.TaskSchedule{
			TaskID:         id,
			EarliestStart:  start,
			EarliestFinish: finish,
		})
	}
	sort.Slice(schedules, func(i, j int) bool {
		if schedules[i].EarliestStart != schedules[j].EarliestStart {
			return schedules[i].EarliestStart < schedules[j].EarliestStart
		}
		return schedules[i].TaskID < schedules[j].TaskID
	})

	// ExecutionOrder: sorted by start time (same time = alphabetical)
	executionOrderSorted := make([]string, 0, len(schedules))
	for _, ts := range schedules {
		executionOrderSorted = append(executionOrderSorted, ts.TaskID)
	}

	return &model.ScheduleResult{
		JobName:           job.Name,
		Workers:           workers,
		MinCompletionTime: currentTime,
		TaskSchedules:     schedules,
		ExecutionOrder:    executionOrderSorted,
		CriticalPath:      nil, // not computed for limited workers
	}, nil
}

func (s *WorkerScheduler) findCriticalPath(job *model.Job, est, eft map[string]int, minCompletion int) []string {
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

func (s *WorkerScheduler) buildSortedSchedules(order []string, est, eft map[string]int) []model.TaskSchedule {
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

func (s *WorkerScheduler) topologicalOrder(job *model.Job) ([]string, error) {
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
