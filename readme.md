Celeriac
=======+

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
import (
	"github.com/svcavallar/celeriac.v1"

func main() {
	taskBrokerURI := "amqp://svcworker:svcworker@127.0.0.1:5672"

	// Connect to RabbitMQ task queue
	service.TaskQueueMgr, err = celeriac.NewTaskQueueMgr(taskBrokerURI)
	if err != nil {
		log.Infof("Failed to connect to task queue: %v", err)
		return err
	}

	log.Infof("Service connected to task queue - (URL: %s)", taskBrokerURI)

	// Get the task events from the task events channel
	go func() {
		for {
			select {
			default:
				ev := <-service.TaskQueueMgr.Monitor.EventsChannel

				if ev != nil {

					if x, ok := ev.(*celeriac.WorkerEvent); ok {
						log.Printf("Task monitor: Worker event - %s", x.Type)
					} else if x, ok := ev.(*celeriac.TaskEvent); ok {
						log.Printf("Task monitor: Task event - %s [ID]: %s", x.Type, x.UUID)
					}

				}
			}
		}
	}()
}
```
