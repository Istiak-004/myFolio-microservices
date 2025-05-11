package config

import "fmt"

// Each service config struct must implement this for validation
func (b Base) Validate() error {
	if b.AppName == "" {
		return fmt.Errorf("APP_NAME is required")
	}
	if b.Env == "" {
		return fmt.Errorf("ENV is required")
	}
	return nil
}
