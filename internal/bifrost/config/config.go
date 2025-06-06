package config

import "github.com/tremendouscan/bifrost/internal/bifrost/options"

// Config is the running configuration structure of the Bifrost pump service.
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given Bifrost pump command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
