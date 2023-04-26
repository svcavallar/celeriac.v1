package celeriac

import (
	"time"
)

/*
Event defines a base event emitted by Celery workers.
*/
type Event struct {
	// Type is the Celery event type. See supported events listed in "constants.go"
	Type string `json:"type"`

	// Hostname is the name of the host on which the Celery worker is operating
	Hostname string `json:"hostname"`

	// Timestamp is the current time of the event
	Timestamp float32 `json:"timestamp"`

	// PID is the process ID
	PID int `json:"pid"`

	// Clock is the current clock time
	Clock int `json:"clock"`

	// UTCOffset is the offset from UTC for the time when this event is valid
	UTCOffset int `json:"utcoffset"`

	// Data is a property allowing extra data to be sent through for custom events from a Celery worker
	Data interface{} `json:"data, omitempty"`
}

/*
NewEvent is a factory function to create a new Event object
*/
func NewEvent() *Event {
	return &Event{
		Type:      "",
		Hostname:  "",
		Timestamp: 0,
		PID:       0,
		Clock:     0,
		UTCOffset: 0,
	}
}

/*
TimestampFormatted returns a formatted string representation of the task event timestamp
*/
func (event *Event) TimestampFormatted() string {
	// First convert float64 to a time.Time
	dt := time.Unix(int64(event.Timestamp), 0)

	// Convert the timestamp to our own string representation
	return dt.UTC().Format(ConstTimeFormat)
}

//------------------------------------------------------------------------------

/*
WorkerEvent defines an event emitted by workers, specific to its operation. Event "types" emitted are:
  - "worker-online"
  - "worker-offline"
  - "worker-heartbeat"

Example worker event json:

	{
		"sw_sys": "Darwin",
		"clock": 74,
		"timestamp": 1843965659.580637,
		"hostname": "celery@worker1.My-Mac.local",
		"pid": 10837,
		"sw_ver": "3.1.18",
		"utcoffset": -11,
		"loadavg": [2.0, 2.41, 2.54],
		"processed": 6,
		"active": 0,
		"freq": 2.0,
		"type": "worker-offline",
		"sw_ident": "py-celery"
	}
*/
type WorkerEvent struct {
	Type      string  `json:"type"`
	Hostname  string  `json:"hostname"`
	Timestamp float32 `json:"timestamp"`
	PID       int     `json:"pid"`
	Clock     int     `json:"clock"`
	UTCOffset int     `json:"utcoffset"`

	// SWSystem is the software system being used
	SWSystem string `json:"sw_sys"`

	// SWVersion is the software version being used
	SWVersion string `json:"sw_ver"`

	// LoadAverage is an array of average CPU loadings for the worker
	LoadAverage []float32 `json:"loadavg"`

	// Freq is the worker frequency use
	Freq float32 `json:"freq"`

	// SWIdentity is the software identity
	SWIdentity string `json:"sw_ident"`

	// Processed is the number of items processed
	Processed int `json:"processed, omitempt"`

	// Active is the active number of workers
	Active int `json:"active, omitempty"`
}

/*
NewWorkerEvent is a factory function to create a new WorkerEvent object
*/
func NewWorkerEvent() *WorkerEvent {
	return &WorkerEvent{
		Type:      "",
		Hostname:  "",
		Timestamp: 0,
		PID:       0,
		Clock:     0,
		UTCOffset: 0,

		SWSystem:    "",
		SWVersion:   "",
		LoadAverage: []float32{},
		Processed:   0,
		Active:      0,
		Freq:        0,
		SWIdentity:  "",
	}
}

//------------------------------------------------------------------------------

/*
TaskEvent defines an event emitted by workers reponding to tasks. Event "types" emitted are:
	- "task-received"
	- "task-started"
	- "task-succeeded"
	- "task-failed"
	- "task-revoked"
	- "task-retried"

Example event json (Task received):
{
	"retries": 0,
	"args": "[]",
	"name": "workers.ingester.ingest",
	"clock": 59,
	"timestamp": 1463965647.276326,
	"pid": 10837,
	"utcoffset": -11,
	"eta": "2018-02-08T08:36:18.186451+00:00",
	"expires": "2018-02-08T09:36:18.186451+00:00",
	"kwargs": "{u'foo': u'bar'}",
	"type": "task-received",
	"hostname": "celery@worker1.My-Mac.local",
	"uuid": "8e42b71d-175b-47c1-52e8-0f2e82ddd9ba"
}
NOTES:
	* Observed that "eta" is sometimes a string (newer versions of Celery shown above), others times an int (older versions of Celery)
	* Observed that "expires" is sometimes a string (newer versions of Celery shown above), others times an int (older versions of Celery)

Example event json (Task started):
{
	"utcoffset": -11,
	"type": "task-started",
	"uuid": "8e42b71d-175b-47c1-52e8-0f2e82ddd9ba",
	"clock": 60,
	"timestamp": 1463965647.278163,
	"hostname": "celery@worker1.My-Mac.local",
	"pid": 10837
}

Example event json (Task succeeded):
{
	"uuid": "8e42b71d-175b-47c1-52e8-0f2e82ddd9ba",
	"clock": 61,
	"timestamp": 1463965647.282055,
	"hostname": "celery@worker1.My-Mac.local",
	"pid": 10837,
	"utcoffset": -11,
	"result": "\\"{'theFoo': u'theBar', 'fooBarCount': 23}\\"",
	"runtime": 0.004214855001919204,
	"type": "task-succeeded"
}

Example event json (Task revoked):
{
	"type": "task-revoked",
	"uuid": "9e4c1cf3-6b66-482e-5145-0eb4d12f91ab",
	"clock": 1549,
	"timestamp": 1444938011.962136,
	"hostname": "celery@worker1.My-Mac.local",
	"pid": 46221,
	"utcoffset": -11,

	"terminated": true,
	"signum": "15",
	"expired": false
}

{
	"terminated": true,
	"type": "task-revoked",
	"uuid": "9e4c1cf3-6b66-482e-5145-0eb4d12f91ab",
	"clock": 1550,
	"timestamp": 1444938011.978788,
	"hostname": "celery@worker1.My-Mac.local",
	"pid": 46221,
	"utcoffset": -11,
	"signum": 15,
	"expired": false
}
NOTES:
	* Observed that "signum" is sometimes returned as a string and sometimes as an integer

*/

/*
TaskEvent is the JSON schema for Celery task events
*/
type TaskEvent struct {
	Type      string  `json:"type"`
	Hostname  string  `json:"hostname"`
	Timestamp float32 `json:"timestamp"`
	PID       int     `json:"pid"`
	Clock     int     `json:"clock"`
	UTCOffset int     `json:"utcoffset"`

	// UUID is the id of the task
	UUID string `json:"uuid"`

	// Name is the textual name of the task executed
	Name string `json:"name,omitempty"`

	// Args is a string of the arguments passed to the task
	Args string `json:"args,omitempty"`

	// Kwargs is a string of the key-word arguments passed to the task
	Kwargs string `json:"kwargs,omitempty"`

	// Runtime is the execution time
	Runtime float32 `json:"runtime,omitempty"`

	// Retries is the number of re-tries this task has performed
	Retries int `json:"retries,omitempty"`

	// ETA is the explicit time and date to run the retry at.
	ETA interface{} `json:"eta,omitempty"`

	// Expires is the datetime or seconds in the future for the task should expire
	Expires interface{} `json:"expires,omitempty"`

	// Terminated is a flag indicating whether the task has been terminated
	Terminated bool `json:"terminated,omitempty"`

	// Signum is the signal number
	Signum interface{} `json:"signum,omitempty"`

	// Expired is a flag indicating whether the task has expired due to factors
	Expired bool `json:"expired,omitempty"`

	// Queue may be set for a Celery task, as a rule cross-reference with RoutingKey
	Queue string `json:"queue,omitempty"`
}

/*
NewTaskEvent is a factory function to create a new TaskEvent object
*/
func NewTaskEvent() *TaskEvent {
	return &TaskEvent{
		Type:      "",
		Hostname:  "",
		Timestamp: 0,
		PID:       0,
		Clock:     0,
		UTCOffset: 0,

		UUID:       "",
		Name:       "",
		Args:       "",
		Runtime:    0,
		Retries:    0,
		ETA:        "",
		Expires:    "",
		Terminated: false,
		Signum:     "",
		Expired:    false,
		Queue:      "",
	}
}

/*
TaskEventsList is an array of task events
*/
type TaskEventsList []TaskEvent
