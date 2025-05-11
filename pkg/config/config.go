package config

// Config is an interface that each service-specific config must implement
type Config interface {
	Validate() error
}
