package app

import (
	"encoding/json"
	"os"
	"time"
)

// Config stores the configuration of the microservice
type Config struct {
	// PublicURL is the host:port the public API server listens on
	PublicURL string `json:"public_url"`
	// PrivateURL is the host:port the private API server listens on
	PrivateURL string `json:"private_url"`
	// BookServiceURL is the host:port of the book microservice
	BookServiceURL string `json:"book_service_url"`
	// UserServiceURL is the host:port of the users microservice
	UserServiceURL string `json:"user_service_url"`
	// DSN is the database connection string
	DSN string `json:"dsn"`
	// BookReturnDeadline is the time span that a user has to return a book after it has been taken
	BookReturnDeadline time.Duration `json:"book_return_deadline"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &Config{}
	err = json.NewDecoder(file).Decode(result)
	return result, err
}
