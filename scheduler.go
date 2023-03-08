package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	funcRegistry *FuncRegistry
	tasks        map[string]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		funcRegistry: NewFuncRegistry(),
		tasks:        make(map[string]*Task),
	}
}

func (scheduler *Scheduler) RunAt(time time.Time, function Function, params ...Param) (string, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := NewTask(funcMeta, params)
	task.SetNextRun(time)
	scheduler.tasks[task.TaskID] = task

	go task.Run()
	return task.TaskID, nil
}

func (scheduler *Scheduler) RunAfter(duration time.Duration, function Function, params ...Param) (string, error) {
	return scheduler.RunAt(time.Now().Add(duration), function, params...)
}

func (scheduler *Scheduler) RunEvery(runEvery time.Duration, runUntil time.Time, function Function, params ...Param) (string, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return "", err
	}

	task := NewTask(funcMeta, params)
	task.SetInterval(runEvery, runUntil)
	scheduler.tasks[task.TaskID] = task

	go task.Run()
	return task.TaskID, nil
}

func (scheduler *Scheduler) Cancel(taskID string) error {
	task, found := scheduler.tasks[taskID]
	if !found {
		return fmt.Errorf("Task %v not found", taskID)
	}

	task.Stop()
	delete(scheduler.tasks, taskID)
	return nil
}

// func (scheduler *Scheduler) ClearExpired() {
// 	for _, task := range scheduler.tasks {
// 		delete(scheduler.tasks, task.TaskID)
// 	}
// }

func (scheduler *Scheduler) ClearAll() {
	for _, task := range scheduler.tasks {
		task.Stop()
		delete(scheduler.tasks, task.TaskID)
	}
	scheduler.funcRegistry = NewFuncRegistry()
}

func (scheduler *Scheduler) Reschedule(taskID string, time time.Time) error {
	task, found := scheduler.tasks[taskID]
	if !found {
		return fmt.Errorf("Task %v not found", taskID)
	}

	task.Stop()
	task.SetNextRun(time)

	go task.Run()
	return nil
}
