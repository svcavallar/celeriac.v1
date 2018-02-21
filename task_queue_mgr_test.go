package celeriac

import (
	"reflect"
	"sync"
	"testing"
	"time"

	// Package dependencies
	log "github.com/Sirupsen/logrus"
)

var (
	queueMgr     *TaskQueueMgr
	brokerURI    = "amqp://admin:password@localhost:5672"
	testLifetime = 60 * time.Second
	taskDispatchInterval = 1000 * time.Millisecond
	taskRoutingKey = "root.test"
	taskName = "root.test.mytask"
)

var wg sync.WaitGroup

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestNewTaskQueueMgr(t *testing.T) {

	queueMgr, err := NewTaskQueueMgr(brokerURI)
	if err != nil {
		log.Fatalf("Error creating new task queue manager, error: %v", err)
	}

	expect(t, err, nil)


	// Configure monitor
	queueMgr.Monitor.monitorWorkerHeartbeatEvents = true

	startTime := time.Now()

	workerEventCount := 0

	// Go routine to monitor the Celery events emitted on the celeriac events channel
	go func() {
		for {
			select {
			default:
				ev := <-queueMgr.Monitor.EventsChannel

				if ev != nil {
					if x, ok := ev.(*WorkerEvent); ok {
						log.Printf("Celery Event Channel: Worker event - %s [Hostname]: %s", x.Type, x.Hostname)
						workerEventCount++
					} else if x, ok := ev.(*TaskEvent); ok {
						log.Printf("Celery Event Channel: Task event - %s [ID]: %s", x.Type, x.UUID)
						workerEventCount++
					} else if x, ok := ev.(*Event); ok {
						log.Printf("Celery Event Channel: General event - %s [Hostname]: %s - [Data]: %v", x.Type, x.Hostname, x.Data)
						workerEventCount++
					} else {
						log.Warnf("Celery Event Channel: Unhandled event type: %v", ev)
					}

				}
			}
		}
	}()

	// Go routine to dispatch a task at a regular time interval for the lifetime of the test
	go func() {
		dispatchStartTime := time.Now()

		for {
			dispatchElapsedTime := time.Since(dispatchStartTime)
			if dispatchElapsedTime.Seconds() >= taskDispatchInterval.Seconds() {
				// Dispatch a new task
				taskData := map[string]interface{}{
					"foo": "bar",
				}

				_, err := queueMgr.DispatchTask(taskName, taskData, taskRoutingKey)
				if err != nil {
					log.Errorf("Failed to dispatch task to queue: %v", err)
				}

				dispatchStartTime = time.Now()
			}
		}
	}()

	// Get the task events from the task events channel
	for {
		elapsedTime := time.Since(startTime)
		if elapsedTime.Seconds() >= testLifetime.Seconds() {
			log.Printf("Processed %d worker events", workerEventCount)

			break
		}
	}
}
