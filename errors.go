package celeriac

import (
	"errors"
	"fmt"

	// Package dependencies
	log "github.com/sirupsen/logrus"
)

// Global Errors
var (
	// ErrInvalidTaskID is raised when an invalid task ID has been detected
	ErrInvalidTaskID = errors.New("invalid task ID specified")

	// ErrInvalidTaskName is raised when an invalid task name has been detected
	ErrInvalidTaskName = errors.New("invalid task name specified")
)

// Fail logs the error and exits the program
// Only use this to handle critical errors
func Fail(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Log only logs the error but doesn't exit the program
// Use this to log errors that should not exit the program
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
