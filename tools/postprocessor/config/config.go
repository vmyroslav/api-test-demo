package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   string         `yaml:"version"`
	Settings  Settings       `yaml:"settings"`
	Patterns  []Pattern      `yaml:"patterns"`
	Endpoints []EndpointRule `yaml:"endpoints,omitempty"`
}

type Settings struct {
	CaseSensitive bool `yaml:"case_sensitive"`
	Debug         bool `yaml:"debug"`
	DecodeBody    bool `yaml:"decode_body"`
	Override      bool `yaml:"override"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// validate all patterns
	for _, pattern := range config.Patterns {
		if err = pattern.validate(); err != nil {
			return nil, fmt.Errorf("invalid pattern %q: %w", pattern.Type, err)
		}
	}

	// validate all endpoints
	if err = config.validateEndpoints(); err != nil {
		return nil, fmt.Errorf("invalid endpoint rule: %w", err)
	}

	return &config, nil
}

func (p *Pattern) validate() error {
	if !p.Type.IsValid() {
		return fmt.Errorf("invalid pattern type: %s", p.Type)
	}

	switch p.Type {
	case DatetimePattern:
		if len(p.Formats) == 0 {
			return fmt.Errorf("datetime pattern requires at least one format")
		}
		// validate each format by attempting to parse a sample date
		for _, format := range p.Formats {
			if _, err := time.Parse(format, time.Now().Format(format)); err != nil {
				return fmt.Errorf("invalid datetime format %q: %v", format, err)
			}
		}
	case PrefixPattern:
		if p.Pattern == "" {
			return fmt.Errorf("prefix pattern requires a pattern string")
		}
		if p.Length <= 0 {
			return fmt.Errorf("prefix pattern requires random_length > 0")
		}
	case UUIDPattern:
		// UUID pattern doesn't require additional validation
	}

	return nil
}

func (c *Config) validateEndpoints() error {
	for _, rule := range c.Endpoints {
		// Ensure method is provided
		if rule.Method == "" {
			return fmt.Errorf("endpoint rule requires method")
		}

		// Ensure path is provided
		if rule.Path == "" {
			return fmt.Errorf("endpoint rule requires path")
		}

		// Validate StaticResponse is valid JSON if provided
		if rule.StaticResponse != "" {
			var js map[string]interface{}
			if err := json.Unmarshal([]byte(rule.StaticResponse), &js); err != nil {
				return fmt.Errorf("invalid JSON in static_response for %s %s: %w",
					rule.Method, rule.Path, err)
			}
		}
	}

	return nil
}
