package config

import (
	"fmt"
)

// PatternType represents the type of pattern to match against.
// It's used to classify different kinds of data we want to match and replace.
type PatternType string

const (
	UUIDPattern     PatternType = "uuid"
	DatetimePattern PatternType = "datetime"
	PrefixPattern   PatternType = "prefix"
)

// String returns the string representation of the PatternType
func (pt PatternType) String() string {
	return string(pt)
}

// IsValid checks if the pattern type is valid
func (pt PatternType) IsValid() bool {
	switch pt {
	case UUIDPattern, DatetimePattern, PrefixPattern:
		return true
	default:
		return false
	}
}

func (pt *PatternType) UnmarshalYAML(unmarshal func(any) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}

	parsedType := PatternType(value)
	if !parsedType.IsValid() {
		return fmt.Errorf("invalid pattern type: %s", value)
	}

	*pt = parsedType
	return nil
}

func (pt PatternType) MarshalYAML() (any, error) {
	return pt.String(), nil
}

// Pattern defines a single pattern configuration that specifies what to match and how to replace it.
type Pattern struct {
	Type        PatternType `yaml:"type"`
	Pattern     string      `yaml:"pattern,omitempty"`
	Formats     []string    `yaml:"formats,omitempty"`
	ReplaceWith string      `yaml:"replace_with,omitempty"`
	Length      int         `yaml:"length,omitempty"`
}
