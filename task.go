package scheduler

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const charsetLength = len(charset)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

type Schedule struct {
	Timer        *time.Timer
	RunAt        time.Time
	IntervalTime time.Time
}

type Task struct {
	TaskID   string
	Schedule Schedule
	Func     FunctionMeta
	Params   []Param
}

func NewTask(function FunctionMeta, params []Param) *Task {
	str := make([]byte, 4)
	for i := range str {
		str[i] = charset[seededRand.Intn(charsetLength)]
	}

	return &Task{
		TaskID: fmt.Sprint(string(str), time.Now().UnixNano()),
		Func:   function,
		Params: params,
	}
}

func (task *Task) SetTime(runAt time.Time) {
	task.Schedule.RunAt = runAt
	task.Schedule.Timer = time.NewTimer(runAt.Sub(time.Now()))
}

func (task *Task) Run() {
	<-task.Schedule.Timer.C
	function := reflect.ValueOf(task.Func.function)
	params := make([]reflect.Value, len(task.Params))
	for i, param := range task.Params {
		params[i] = reflect.ValueOf(param)
	}
	function.Call(params)
}

func (task *Task) Stop() {
	task.Schedule.Timer.Stop()
}
