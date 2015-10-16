package celeriac

import ()

/*
baseCmd defines a base command for interacting with worker tasks
*/
type baseCmd struct {
	Command     string `json:"method"`
	Destination string `json:"destination"`
}

//------------------------------------------------------------------------------

type revokeTaskArgs struct {
	TaskID    string `json:"task_id"`
	Terminate bool   `json:"terminate"`
	Signal    string `josn:"signal"`
}

// RevokeTaskCmd is a wrapper to a command
type RevokeTaskCmd struct {
	baseCmd
	Arguments revokeTaskArgs `json:"arguments"`
}

/*
NewRevokeTaskCmd creates a new command for revoking a task by given id

If a task is revoked, the workers will ignore the task and not execute it after all.
*/
func NewRevokeTaskCmd(taskID string, terminateProcess bool) *RevokeTaskCmd {
	return &RevokeTaskCmd{
		baseCmd: baseCmd{
			Command:     "revoke",
			Destination: "",
		},
		Arguments: revokeTaskArgs{
			TaskID:    taskID,
			Terminate: terminateProcess,
			Signal:    "SIGTERM",
		},
	}
}

//------------------------------------------------------------------------------

// PingCmd is a wrapper to a command
type PingCmd struct {
	baseCmd
}

// NewPingCmd creates a new command for pinging workers
func NewPingCmd() *PingCmd {
	return &PingCmd{
		baseCmd: baseCmd{
			Command:     "ping",
			Destination: "",
		},
	}
}

//------------------------------------------------------------------------------

// Set rate limit for task by type
type rateLimitTaskArgs struct {
	TaskName  string `json:"task_name"`
	RateLimit string `json:"rate_limit"`
}

// RateLimitTaskCmd is a wrapper to a command
type RateLimitTaskCmd struct {
	baseCmd
	Arguments rateLimitTaskArgs `json:"arguments"`
}

/*
NewRateLimitTaskCmd creates a new command for rate limiting a task

taskName: Name of task to change rate limit for
rateLimit: The rate limit as tasks per second, or a rate limit string (`"100/m"`, etc.
            see :attr:`celery.task.base.Task.rate_limit` for more information)
*/
func NewRateLimitTaskCmd(taskName string, rateLimit string) *RateLimitTaskCmd {
	return &RateLimitTaskCmd{
		baseCmd: baseCmd{
			Command:     "rate_limit",
			Destination: "",
		},
		Arguments: rateLimitTaskArgs{
			TaskName:  taskName,
			RateLimit: rateLimit,
		},
	}
}

//------------------------------------------------------------------------------

// Set time limit for task by type
type timeLimitTaskArgs struct {
	TaskName  string `json:"task_name"`
	HardLimit string `json:"hard"`
	SoftLimit string `json:"soft"`
}

// TimeLimitTaskCmd is a wrapper to a command
type TimeLimitTaskCmd struct {
	baseCmd
	Arguments timeLimitTaskArgs `json:"arguments"`
}

/*
NewTimeLimitTaskCmd creates a new command for rate limiting a task

taskName: Name of task to change rate limit for
hardLimit: New hard time limit (in seconds)
softLimit: New soft time limit (in seconds)
*/
func NewTimeLimitTaskCmd(taskName string, hardLimit string, softLimit string) *TimeLimitTaskCmd {
	return &TimeLimitTaskCmd{
		baseCmd: baseCmd{
			Command:     "time_limit",
			Destination: "",
		},
		Arguments: timeLimitTaskArgs{
			TaskName:  taskName,
			HardLimit: hardLimit,
			SoftLimit: softLimit,
		},
	}
}
