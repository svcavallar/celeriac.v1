package celeriac

import (
	"encoding/json"
	"time"

	// Package dependencies
	"github.com/nu7hatch/gouuid"
)

/*
Task is a representation of a Celery task
*/
type Task struct {
	// TaskName is the name of the task
	TaskName string

	// ID is the task UUID
	ID string

	// Args are task arguments (optional)
	Args []string

	// KWArgs are keyword arguments (optional)
	KWArgs map[string]interface{}

	// Retries is a number of retries to perform if an error occurs (optional)
	Retries int

	// ETA is the estimated completion time (optional)
	ETA time.Time

	// Expires is the time when this task will expire (optional)
	Expires time.Time
}

/*
NewTask is a factory function that creates and returns a pointer to a new task object
*/
func NewTask(taskName string, args []string, kwargs map[string]interface{}) (*Task, error) {
	newTaskID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	newTask := Task{
		TaskName: taskName,
		ID:       newTaskID.String(),
		Args:     args,
		KWArgs:   kwargs,
		Retries:  0,
	}

	return &newTask, nil
}

/*
MarshalJSON marshals a Task object into a json bytes array

Time properties are converted to UTC and formatted in ISO8601
*/
func (task *Task) MarshalJSON() ([]byte, error) {
	type _jsonTask struct {
		TaskName string                 `json:"task"`
		ID       string                 `json:"id"`
		Args     []string               `json:"args, omitempty"`
		KWArgs   map[string]interface{} `json:"kwargs, omitempty"`
		Retries  int                    `json:"retries, omitempty"`
		ETA      string                 `json:"eta, omitempty"`
		Expires  string                 `json:"expires, omitempty"`
	}

	out := _jsonTask{
		TaskName: task.TaskName,
		ID:       task.ID,
		Args:     task.Args,
		KWArgs:   task.KWArgs,
		Retries:  task.Retries,
	}

	// Convert time properties to UTC, and format as ISO8601
	if !task.ETA.IsZero() {
		out.ETA = task.ETA.UTC().Format(ConstTimeFormat)
	}

	if !task.Expires.IsZero() {
		out.Expires = task.Expires.UTC().Format(ConstTimeFormat)
	}

	return json.Marshal(out)
}
