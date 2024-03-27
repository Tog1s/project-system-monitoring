package config

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Metrics MetricsConfig
}

type MetricsConfig struct {
	LoadAverage bool
}

func New() Config {
	return Config{}
}

func ReadFromFile(path string) (*Config, error) {
	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading config file %s: %w", path, err)
	}

	config := New()
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, fmt.Errorf("error while parse yaml file %s: %w", path, err)
	}

	return &config, nil
}
