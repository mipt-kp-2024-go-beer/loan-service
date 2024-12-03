package app

import "time"

// Config stores the configuration of the microservice
type Config struct {
	// Host is the host the server listens on
	Host string
	// Port is the port the server listens on
	Port uint16
	// DSN is the database connection string
	DSN string
	// BookReturnDeadline is the time span that a user has to return a book after it has been taken
	BookReturnDeadline time.Duration
}
