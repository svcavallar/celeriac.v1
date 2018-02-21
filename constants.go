package celeriac

const (
	// ConstPublishTaskContentType is the content type of the task data to be published
	ConstPublishTaskContentType = "application/json"

	// ConstPublishTaskContentEncoding is the content encoding type of the task data to be published
	ConstPublishTaskContentEncoding = "utf-8"

	// ConstTaskDefaultExchangeName is the default exchange name to use when publishing a task
	ConstTaskDefaultExchangeName = ""

	// ConstTaskDefaultRoutingKey is the default routing key to use when publishing a task
	ConstTaskDefaultRoutingKey = "celery"

	// ConstTaskControlExchangeName is the exchange name for dispatching task control commands
	ConstTaskControlExchangeName = "celery.pidbox"

	// ConstEventsMonitorExchangeName is the exchange name used for Celery events
	ConstEventsMonitorExchangeName = "celeryev"

	// ConstEventsMonitorExchangeType is the exchange type for the events monitor
	ConstEventsMonitorExchangeType = "topic"

	// ConstEventsMonitorQueueName is the queue name of the events monitor
	ConstEventsMonitorQueueName = "celeriac-events-monitor-queue"

	// ConstEventsMonitorBindingKey is the binding key for the events monitor
	ConstEventsMonitorBindingKey = "*.*"

	// ConstEventsMonitorConsumerTag is the consumer tag name for the events monitor
	ConstEventsMonitorConsumerTag = "celeriac-events-monitor"

	// ConstTimeFormat is the general format for all timestamps
	ConstTimeFormat = "2006-01-02T15:04:05.999999"

	// Event type names (as strings)

	// ConstEventTypeWorkerOnline is the event type when a Celery worker comes online
	ConstEventTypeWorkerOnline = "worker-online"

	// ConstEventTypeWorkerOffline is the event type when a Celery worker goes offline
	ConstEventTypeWorkerOffline = "worker-offline"

	// ConstEventTypeWorkerHeartbeat is the event type when a Celery worker is online and "alive"
	ConstEventTypeWorkerHeartbeat = "worker-heartbeat"

	// ConstEventTypeTaskSent is the event type when a Celery task is sent
	ConstEventTypeTaskSent = "task-sent"

	// ConstEventTypeTaskReceived is the event type when a Celery worker receives a task
	ConstEventTypeTaskReceived = "task-received"

	// ConstEventTypeTaskStarted is the event type when a Celery worker starts a task
	ConstEventTypeTaskStarted = "task-started"

	// ConstEventTypeTaskSucceeded is the event type when a Celery worker completes a task
	ConstEventTypeTaskSucceeded = "task-succeeded"

	// ConstEventTypeTaskFailed is the event type when a Celery worker fails to complete a task
	ConstEventTypeTaskFailed = "task-failed"

	// ConstEventTypeTaskRevoked is the event type when a Celery worker has its task revoked
	ConstEventTypeTaskRevoked = "task-revoked"

	// ConstEventTypeTaskRetried is the event type when a Celery worker retries a task
	ConstEventTypeTaskRetried = "task-retried"
)
