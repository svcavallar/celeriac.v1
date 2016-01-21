Celeriac
=======

Golang client library for adding support for interacting and monitoring Celery workers and tasks.

It provides functionality to place tasks on the task queue, as well
as monitor task and worker events.

Dependencies
------------
This library depends upon the following packages:
- github.com/streadway/amqp
- github.com/Sirupsen/logrus
- github.com/nu7hatch/gouuid
- github.com/pquerna/ffjson

Usage
-----
Installation: `go get github.com/svcavallar/celeriac.v1`

This imports a new namespace called `celeriac`

```go
package main

import (
	"log"
	"os"

	"github.com/svcavallar/celeriac.v1"
)

func main() {
	taskBrokerURI := "amqp://user:pass@localhost:5672/vhost"

	// Connect to RabbitMQ task queue
	TaskQueueMgr, err := celeriac.NewTaskQueueMgr(taskBrokerURI)
	if err != nil {
		log.Printf("Failed to connect to task queue: %v", err)
		os.Exit(-1)
	}

	log.Printf("Service connected to task queue - (URL: %s)", taskBrokerURI)

	// Get the task events from the task events channel
	for {
		select {
		default:
			ev := <-TaskQueueMgr.Monitor.EventsChannel

			if ev != nil {

				if x, ok := ev.(*celeriac.WorkerEvent); ok {
					log.Printf("Task monitor: Worker event - %s", x.Type)
				} else if x, ok := ev.(*celeriac.TaskEvent); ok {
					log.Printf("Task monitor: Task event - %s [ID]: %s", x.Type, x.UUID)
				}

			}
		}
	}
}
```

Dispatching Tasks
-----------------

### By Name
This will create and dispatch a task incorporating the supplied data. The task will automatically be allocated and identified by a UUID returned in the task object. The UUID is represented in the form of "6ba7b810-9dad-11d1-80b4-00c04fd430c8".

	// Dispatch a new task
	taskName := "root.test.task"
	taskData := map[string]interface{}{
		"foo": "bar"
	}
	routingKey := "root.test"
	
	task, err := TaskQueueMgr.DispatchTask(taskName, taskData, routingKey)
	if err != nil {
		log.Errorf("Failed to dispatch task to queue: %v", err)
	}


### By ID & Name
This will create and dispatch a task incorporating the supplied data, and identified by the user-supplied task identifier. 

	// Dispatch a new task
	taskID := "my_task_id_123456789"
	taskName := "root.test.task"
	taskData := map[string]interface{}{
		"foo": "bar"
	}
	routingKey := "root.test"
	
	task, err := TaskQueueMgr.DispatchTaskWithID(taskID, taskName, taskData, routingKey)
	if err != nil {
		log.Errorf("Failed to dispatch task to queue: %v", err)
	}

