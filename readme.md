# Celeriac

Golang client library for adding support for interacting and monitoring Celery workers and tasks.

It provides functionality to place tasks on the task queue, as well
as monitor both task and worker events.

## Dependencies

This library depends upon the following packages:

- github.com/streadway/amqp
- github.com/Sirupsen/logrus
- github.com/nu7hatch/gouuid
- github.com/mailru/easyjson

## Install `easyjson`

```bash
$ go get -u github.com/mailru/easyjson/...
```

## Usage

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

	// Go routine to monitor the Celery events emitted on the celeriac events channel
	go func() {
        for {
            select {
            default:
                ev := <-TaskQueueMgr.Monitor.EventsChannel

                if ev != nil {

                    if x, ok := ev.(*celeriac.WorkerEvent); ok {
                        log.Printf("Celery Event Channel: Worker event - %s [Hostname]: %s", x.Type, x.Hostname)
                    } else if x, ok := ev.(*celeriac.TaskEvent); ok {
                        log.Printf("Celery Event Channel: Task event - %s [ID]: %s", x.Type, x.UUID)
                    } else if x, ok := ev.(*celeriac.Event); ok {
                        log.Printf("Celery Event Channel: General event - %s [Hostname]: %s - [Data]: %v", x.Type, x.Hostname, x.Data)
                    } else {
                        log.Printf("Celery Event Channel: Unhandled event: %v", ev)
                    }
                }
            }
        }
	}()
}
```

## Dispatching Tasks

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

## Modifying `task_event.go`

If you modify the properties of any of the structs in `task_event.go` you will need to re-generate the `easyjson` version of this file. This is easily achieved by issuing the following command:

```bash
$ easyjson -all task_event.go
```
## Processing Redis Backend Result Automatically
If you are using a Redis backend for storing results you can easily process new/updated entries by subscribing to Redis keyspace events. This will save polling for results, and is made convenient to integrate by using my golang helper package `go-redis-event-sink`, available at the repo [https://github.com/svcavallar/go-redis-event-sink](https://github.com/svcavallar/go-redis-event-sink)

An example test on how to use are provided within the repository. Essentially, for Celery, just provide it with the celery task naming mask patten to watch: `celery-task-meta-*`