package db

import "backend/internal/utils/concurrency"

// Database interface defines the basic methods needed for any database implementation
// It includes methods for connecting, getting thread pool, and closing the database connection
type Database[E DatabaseConfigurations] interface {
	Name() string
	Connect(config E) error
	GetThreadPool() *concurrency.ThreadPool
	Close() error
}

// DatabaseConfigurations is an interface that defines the method to retrieve authentication configurations
type DatabaseConfigurations interface {
	GetAuthenticationConfigurations() AuthenticationConfigurations // Returns authentication configurations
}

// AuthenticationConfigurations holds the authentication details required to connect to the database
type AuthenticationConfigurations struct {
	Host     string
	Port     int
	Username string
	Password string
}

// GetAuthenticationConfigurations returns the authentication configurations for the database connection
func (config *AuthenticationConfigurations) GetAuthenticationConfigurations() AuthenticationConfigurations {
	return *config
}

// URIConfigurations holds the authentication configurations and the URI for the database connection
type URIConfigurations struct {
	AuthConfig AuthenticationConfigurations
	URI        string
}

// GetAuthenticationConfigurations returns the authentication configurations for URI-based database connection
func (config *URIConfigurations) GetAuthenticationConfigurations() AuthenticationConfigurations {
	return config.AuthConfig
}
