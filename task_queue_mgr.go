package celeriac

import (
	"encoding/json"
	"time"

	// Package dependencies
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
    "strings"
    "crypto/tls"
)

/*
TaskQueueMgr defines a manager for interacting with a Celery task queue
*/
type TaskQueueMgr struct {
	brokerURI	string
	connection *amqp.Connection
	channel    *amqp.Channel
	Monitor    *TaskMonitor
	errorChannel chan *amqp.Error
	closed       bool
}

/*
NewTaskQueueMgr is a factory function that creates a new instance of the TaskQueueMgr
*/
func NewTaskQueueMgr(brokerURI string) (*TaskQueueMgr, error) {
	self := &TaskQueueMgr{
		brokerURI: brokerURI,
		errorChannel: make(chan *amqp.Error),
	}

	err := self.connect()
	if err != nil {
		return nil, err
	}

	// Setup broker reconnection monitor
	go self.brokerReconnector()

	return self, nil
}

func (taskQueueMgr *TaskQueueMgr) connect() error {
	for {
		var err error

		// Connect to the task queue
		if strings.HasPrefix(taskQueueMgr.brokerURI, "amqps") {
			tlsConfig := &tls.Config{}
			//tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
			tlsConfig.InsecureSkipVerify = true

			taskQueueMgr.connection, err = amqp.DialTLS(taskQueueMgr.brokerURI, tlsConfig)
		} else {
			taskQueueMgr.connection, err = amqp.Dial(taskQueueMgr.brokerURI)
		}

		if err != nil {
			log.Errorf("Failed to connect to AMQP queue: %v. Retrying...", err)
			time.Sleep(1000 * time.Millisecond)
		} else {
			taskQueueMgr.errorChannel = make(chan *amqp.Error)

			// Be informed when the connection is closed so we can reconnect automatically
			taskQueueMgr.connection.NotifyClose(taskQueueMgr.errorChannel)

			log.Printf("Established AMQP connection, getting Channel")
			taskQueueMgr.channel, err = taskQueueMgr.connection.Channel()
			if err != nil {
				log.Errorf("Failed to open AMQP channel: %v", err)
				return err
			}

			// Create the task monitor
			// Currently the monitor has one queue for all events
			taskQueueMgr.Monitor, err = NewTaskMonitor(taskQueueMgr.connection,
				taskQueueMgr.channel,
				ConstEventsMonitorExchangeName,
				ConstEventsMonitorExchangeType,
				ConstEventsMonitorQueueName,
				ConstEventsMonitorBindingKey,
				ConstEventsMonitorConsumerTag)

			if err != nil {
				log.Errorf("%s", err)
				return err
			}

			return nil
		}
	}
}

func (taskQueueMgr *TaskQueueMgr) brokerReconnector() {
	for {
		err := <-taskQueueMgr.errorChannel
		if !taskQueueMgr.closed {
			log.Errorf("Connection closed. Reconnecting... (%v)", err)

			taskQueueMgr.connect()
		}
	}
}

/*
Close performs appropriate cleanup of any open task queue connections
*/
func (taskQueueMgr *TaskQueueMgr) Close() {
	taskQueueMgr.closed = true

	// Stop monitoring
	taskQueueMgr.Monitor.Shutdown()

	// Close connections
	if taskQueueMgr.connection != nil {
		taskQueueMgr.connection.Close()
	}
}

/*
publish publishes data onto an AMQP channel via the specified exchange name and routing key
*/
func (taskQueueMgr *TaskQueueMgr) publish(data interface{}, exchangeName string, routingKey string) error {
	// Non-blocking channel where if there is no error its simply ignored
	select {
		case err := <-taskQueueMgr.errorChannel:
			if err != nil {
				taskQueueMgr.connect()
			}
		default:
	}

	bodyData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode:    amqp.Persistent,
		Timestamp:       time.Now(),
		ContentType:     ConstPublishTaskContentType,
		ContentEncoding: ConstPublishTaskContentEncoding,
		Body:            bodyData,
	}

	return taskQueueMgr.channel.Publish(exchangeName, routingKey, false, false, msg)
}

/*
DispatchTask places a new task on the Celery task queue
Creates a new Task based on the supplied task name and data
*/
func (taskQueueMgr *TaskQueueMgr) DispatchTask(taskName string, taskData map[string]interface{}, routingKey string) (*Task, error) {
	var err error

	task, err := taskQueueMgr.DispatchTaskWithID("", taskName, taskData, routingKey)
	return task, err
}

/*
DispatchTaskWithID places a new task with the specified ID on the Celery task queue
Creates a new Task based on the supplied task name and data
*/
func (taskQueueMgr *TaskQueueMgr) DispatchTaskWithID(taskID string, taskName string, taskData map[string]interface{}, routingKey string) (*Task, error) {
	var err error

	task, err := NewTaskWithID(taskID, taskName, nil, taskData)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
		panic(err)
	}

	if len(routingKey) == 0 || routingKey == "" {
		routingKey = ConstTaskDefaultRoutingKey
	}

	err = taskQueueMgr.publish(task, ConstTaskDefaultExchangeName, routingKey)
	log.Infof("Dispatched task [NAME]: %s, [ID]: %s to task queue with [ROUTING KEY]: %s", taskName, task.ID, routingKey)

	return task, err
}

/*
RevokeTask attempts to notify Celery workers that the specified task needs revoking
*/
func (taskQueueMgr *TaskQueueMgr) RevokeTask(taskID string) error {
	if taskID == "" || len(taskID) == 0 {
		return ErrInvalidTaskID
	}

	log.Infof("Revoking task [ID]: %s", taskID)

	rt := NewRevokeTaskCmd(taskID, true)
	return taskQueueMgr.publish(rt, ConstTaskControlExchangeName, ConstTaskDefaultRoutingKey)
}

/*
Ping attempts to ping Celery workers
*/
func (taskQueueMgr *TaskQueueMgr) Ping() error {
	log.Infof("Sending ping to workers")

	rt := NewPingCmd()
	return taskQueueMgr.publish(rt, ConstTaskControlExchangeName, ConstTaskDefaultRoutingKey)
}

/*
RateLimitTask attempts to set rate limit tasks by type
*/
func (taskQueueMgr *TaskQueueMgr) RateLimitTask(taskName string, rateLimit string) error {
	if taskName == "" || len(taskName) == 0 {
		return ErrInvalidTaskName
	}

	rt := NewRateLimitTaskCmd(taskName, rateLimit)
	return taskQueueMgr.publish(rt, ConstTaskControlExchangeName, ConstTaskDefaultRoutingKey)
}

/*
TimeLimitTask attempts to set time limits for task by type
*/
func (taskQueueMgr *TaskQueueMgr) TimeLimitTask(taskName string, hardLimit string, softLimit string) error {
	if taskName == "" || len(taskName) == 0 {
		return ErrInvalidTaskName
	}

	rt := NewTimeLimitTaskCmd(taskName, hardLimit, softLimit)
	return taskQueueMgr.publish(rt, ConstTaskControlExchangeName, ConstTaskDefaultRoutingKey)
}
