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
	task.SetTime(time)
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

func (scheduler *Scheduler) Reschedule(taskID string, time time.Time) error {
	task, found := scheduler.tasks[taskID]
	if !found {
		return fmt.Errorf("Task %v not found", taskID)
	}

	task.Stop()
	task.SetTime(time)

	go task.Run()
	return nil
}
