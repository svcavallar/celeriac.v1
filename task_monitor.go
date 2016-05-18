package celeriac

import (
	"fmt"

	// Package dependencies
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

/*
TaskMonitor is a Celery task event consumer
*/
type TaskMonitor struct {
	connection                   *amqp.Connection
	channel                      *amqp.Channel
	consumerTag                  string
	done                         chan error
	processMessageChannel        bool
	deliveryChannel              <-chan amqp.Delivery
	monitorWorkerHeartbeatEvents bool

	// Public channel on which events are piped
	EventsChannel chan interface{}
}

/*
NewTaskMonitor is a factory function that creates a new Celery consumer
*/
func NewTaskMonitor(connection *amqp.Connection, channel *amqp.Channel,
	exchangeName string, exchangeType string, queueName string, bindingKey string, ctag string) (*TaskMonitor, error) {

	monitor := &TaskMonitor{
		connection:                   connection,
		channel:                      channel,
		consumerTag:                  ctag,
		done:                         make(chan error),
		EventsChannel:                make(chan interface{}),
		monitorWorkerHeartbeatEvents: false,
	}

	var err error

	log.Printf("Declaring exchange (%q)", exchangeName)
	if err = monitor.channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("Declared exchange, declaring queue %q", queueName)
	queue, err := monitor.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue declare: %s", err)
	}

	log.Printf("Declared queue (%q %d messages, %d consumers), binding to Exchange (binding key %q)",
		queue.Name, queue.Messages, queue.Consumers, bindingKey)

	if err = monitor.channel.QueueBind(
		queue.Name,   // name of the queue
		bindingKey,   // binding key
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to exchange, starting monitoring (consumer tag %q)", monitor.consumerTag)
	monitor.deliveryChannel, err = monitor.channel.Consume(
		queue.Name,          // name
		monitor.consumerTag, // consumerTag,
		false,               // noAck
		false,               // exclusive
		false,               // noLocal
		false,               // noWait
		nil,                 // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	// Spawn a Go routine to handle event monitoring
	go monitor.handle(monitor.deliveryChannel, monitor.done, monitor.EventsChannel)

	return monitor, nil
}

/*
Shutdown stops all monitoring, cleaning up any open connections
*/
func (monitor *TaskMonitor) Shutdown() error {
	// Close the events channel
	close(monitor.EventsChannel)

	// Will close() the deliveries channel
	if err := monitor.channel.Cancel(monitor.consumerTag, true); err != nil {
		return fmt.Errorf("Monitor cancel failed: %s", err)
	}

	if err := monitor.connection.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// Wait for handle() to exit
	return <-monitor.done
}

/*
SetMonitorWorkerHeartbeatEvents sets the property whether to process heartbeat events emitted
by workers.

NOTE: By default this is set to 'false' so as to minimize unnecessary "noisy heartbeat" events.
*/
func (monitor *TaskMonitor) SetMonitorWorkerHeartbeatEvents(processHeartbeatEvents bool) {
	monitor.monitorWorkerHeartbeatEvents = processHeartbeatEvents
}

/*
Handles Celery event messages on the task queue
*/
func (monitor *TaskMonitor) handle(deliveries <-chan amqp.Delivery, done chan error, out chan interface{}) {

	rawEvent := NewEvent()
	for d := range deliveries {

		// This code is here purely for debugging ALL messages, and for extracting ones we are unsure of the json format!
		//log.Printf("Received %d bytes: [%v] %q",
		//	len(d.Body),
		//	d.DeliveryTag,
		//	d.Body,
		//)

		// NOTE: We use "ffjson" for performance over the standard Golang json decoding package
		err := rawEvent.UnmarshalJSON(d.Body)
		if err != nil {
			fmt.Printf("Error: %v", err)
		}

		var celeryEvent interface{}

		switch rawEvent.Type {

		case ConstEventTypeWorkerOnline,
			ConstEventTypeWorkerOffline:

			var t = NewWorkerEvent()
			err := t.UnmarshalJSON(d.Body)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			celeryEvent = t
			break

		case ConstEventTypeWorkerHeartbeat:
			if monitor.monitorWorkerHeartbeatEvents {
				var t = NewWorkerEvent()
				err := t.UnmarshalJSON(d.Body)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
				celeryEvent = t
			}
			break

		case ConstEventTypeTaskSent,
			ConstEventTypeTaskReceived,
			ConstEventTypeTaskStarted,
			ConstEventTypeTaskSucceeded,
			ConstEventTypeTaskFailed,
			ConstEventTypeTaskRevoked,
			ConstEventTypeTaskRetried:

			var t = NewTaskEvent()
			err := t.UnmarshalJSON(d.Body)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			celeryEvent = t
			break

		default:
			celeryEvent = rawEvent
			break
		}

		// Publish the task event through our channel
		out <- celeryEvent

		d.Ack(true)
	}

	log.Printf("handle: deliveries channel closed")
	done <- nil

}
