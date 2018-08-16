package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
	funcRegistry *FuncRegistry
	tasks        map[int64]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		funcRegistry: NewFuncRegistry(),
		tasks:        make(map[int64]*Task),
	}
}

func (scheduler *Scheduler) RunAt(time time.Time, function Function, params ...Param) (int64, error) {
	funcMeta, err := scheduler.funcRegistry.Add(function)
	if err != nil {
		return 0, err
	}

	task := NewTask(funcMeta, params)
	task.SetTime(time)
	scheduler.tasks[task.TaskID] = task

	go task.Run()
	return task.TaskID, nil
}

func (scheduler *Scheduler) Cancel(taskID int64) error {
	task, found := scheduler.tasks[taskID]
	if !found {
		return fmt.Errorf("Task %v not found", taskID)
	}

	task.Stop()
	delete(scheduler.tasks, taskID)
	return nil
}

func (scheduler *Scheduler) Reschedule(taskID int64, time time.Time) error {
	task, found := scheduler.tasks[taskID]
	if !found {
		return fmt.Errorf("Task %v not found", taskID)
	}

	task.Stop()
	task.SetTime(time)

	go task.Run()
	return nil
}
