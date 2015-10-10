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
	brokerURI    = "amqp://svcworker:svcworker@localhost:5672/"
	testLifetime = 10 * time.Second
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

	startTime := time.Now()

	taskEventCount := 0

	// Get the task events from the task events channel
	for {
		elapsedTime := time.Since(startTime)
		if elapsedTime.Seconds() >= testLifetime.Seconds() {
			log.Printf("Processed %d task events", taskEventCount)
			break
		}

		// Monitor the Celery events emitted on the celeriac events channel
		go func() {
			for {
				select {
				default:
					ev := <-queueMgr.Monitor.EventsChannel

					if ev != nil {

						if x, ok := ev.(*WorkerEvent); ok {
							log.Printf("Task monitor: Worker event - %s [Hostname]: %s", x.Type, x.Hostname)
						} else if x, ok := ev.(*TaskEvent); ok {
							log.Printf("Task monitor: Task event - %s [ID]: %s", x.Type, x.UUID)
						}

					}
				}
			}
		}()
	}

}
