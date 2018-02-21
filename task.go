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
	TaskName string `json:"task"`

	// ID is the task UUID
	ID string `json:"id"`

	// Args are task arguments (optional)
	Args []string `json:"args, omitempty"`

	// KWArgs are keyword arguments (optional)
	KWArgs map[string]interface{} `json:"kwargs, omitempty"`

	// Retries is a number of retries to perform if an error occurs (optional)
	Retries int `json:"retries, omitempty"`

	// ETA is the estimated completion time (optional)
	ETA *time.Time `json:"eta, omitempty"`

	// Expires is the time when this task will expire (optional)
	Expires *time.Time `json:"expires, omitempty"`
}

/*
NewTask is a factory function that creates and returns a pointer to a new task object
*/
func NewTask(taskName string, args []string, kwargs map[string]interface{}) (*Task, error) {

	taskUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	newTask, err := NewTaskWithID(taskUUID.String(), taskName, args, kwargs)
	if err != nil {
		return nil, err
	}

	return newTask, nil
}

/*
NewTaskWithID is a factory function that creates and returns a pointer to a new task object, allowing caller to
specify the task ID.
*/
func NewTaskWithID(taskID string, taskName string, args []string, kwargs map[string]interface{}) (*Task, error) {

	newTaskID := taskID

	if newTaskID == "" {
		taskUUID, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		newTaskID = taskUUID.String()
	}

	if args == nil {
		args = []string{}
	}

	newTask := Task {
		TaskName: taskName,
		ID:       newTaskID,
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
	out := Task {
		TaskName: task.TaskName,
		ID:       task.ID,
		Args:     task.Args,
		KWArgs:   task.KWArgs,
		Retries:  task.Retries,
	}

	// Convert time properties to UTC, and format as ISO8601
	if task.ETA != nil && !task.ETA.IsZero() {
		*out.ETA = task.ETA.UTC()
	}

	if task.Expires != nil && !task.Expires.IsZero() {
		*out.Expires = task.Expires.UTC()
	}

	return json.Marshal(out)
}
