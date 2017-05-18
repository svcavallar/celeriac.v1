package celeriac

import (
	"encoding/json"
	"fmt"
	"time"

	// Package dependencies
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
    "strings"
    "crypto/tls"
)

/*
TaskQueueMgr defines a manager for interoperating with a Celery task queue
*/
type TaskQueueMgr struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	Monitor    *TaskMonitor
}

/*
NewTaskQueueMgr is a factory function that creates a new instance of the TaskQueueMgr
*/
func NewTaskQueueMgr(brokerURI string) (*TaskQueueMgr, error) {
	self := &TaskQueueMgr{}

	err := self.init(brokerURI)
	if err != nil {
		return nil, err
	}

	return self, nil
}

func (taskQueueMgr *TaskQueueMgr) init(brokerURI string) error {
	var err error

	// Connect to the task queue
    if strings.HasPrefix(brokerURI, "amqps") {
        tlsConfig := &tls.Config{}
        //tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
        tlsConfig.InsecureSkipVerify = true

        taskQueueMgr.connection, err = amqp.DialTLS(brokerURI, tlsConfig)
    } else {
        taskQueueMgr.connection, err = amqp.Dial(brokerURI)
    }

	if err != nil {
		log.Errorf("Failed to connect to AMQP queue: %v", err)
		return err
	}

	go func() {
		fmt.Printf("Closing: %s", <-taskQueueMgr.connection.NotifyClose(make(chan *amqp.Error)))
	}()

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

/*
Close performs appropriate cleanup of any open task queue connections
*/
func (taskQueueMgr *TaskQueueMgr) Close() {
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
	log.Infof("Dispatched task [NAME]: %s, [ID]:%s to task queue with [ROUTING KEY]:%s", taskName, task.ID, routingKey)

	return task, err
}

/*
RevokeTask attempts to notify Celery workers that the specified task needs revoking
*/
func (taskQueueMgr *TaskQueueMgr) RevokeTask(taskID string) error {
	if taskID == "" || len(taskID) == 0 {
		return ErrInvalidTaskID
	}

	log.Infof("Revoking task [ID]:%s", taskID)

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
