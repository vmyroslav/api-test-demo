package processor

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/matching/matchers"
)

func TestRegexUUIDMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name:    "exact simple match",
			pattern: `{"id":"123"}`,
			input:   `{"id":"123"}`,
			want:    true,
		},
		{
			name:    "exact match with .*",
			pattern: `{"id":".*"}`,
			input:   `{"id":"123"}`,
			want:    true,
		},
		{
			name:    "simple flexible match",
			pattern: `{"id":".*","name":"test"}`,
			input:   `{"id":"123","name":"test"}`,
			want:    true,
		},
		{
			name:    "uuid match",
			pattern: `{"id":"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"}`,
			input:   `{"id":"550e8400-e29b-41d4-a716-446655440000"}`,
			want:    true,
		},
		{
			name:    "real world example - same order",
			pattern: `{"external_id":".*","id":".*","name":"test-capture"}`,
			input:   `{"external_id":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","id":"550e8400-e29b-41d4-a716-446655440000","name":"test-capture"}`,
			want:    true,
		},
		{
			name:    "alternative field order pattern",
			pattern: `{"name":"test","id":".*"}|{"id":".*","name":"test"}`,
			input:   `{"name":"test","id":"123"}`,
			want:    true,
		},
		{
			name:    "exact match pattern - alternative order",
			pattern: `{"name":"test","id":"123"}`,
			input:   `{"name":"test","id":"123"}`,
			want:    true,
		},
		{
			name:    "non-matching name",
			pattern: `{"id":".*","name":"test"}`,
			input:   `{"id":"123","name":"other"}`,
			want:    false,
		},
		{
			name:    "missing field",
			pattern: `{"id":".*","name":"test"}`,
			input:   `{"id":"123"}`,
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchers.RegexMatch(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("Pattern: %q\nInput:   %q\nExpected: %v\nGot: %v", tc.pattern, tc.input, tc.want, got)
			}
		})
	}
}

func TestPrefixRegexMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name:    "simple prefix match",
			pattern: `{"name":"fake-api-test-.*"}`,
			input:   `{"name":"fake-api-test-12345"}`,
			want:    true,
		},
		{
			name:    "prefix with exact length match",
			pattern: `{"name":"fake-api-test-.{5}"}`,
			input:   `{"name":"fake-api-test-12345"}`,
			want:    true,
		},
		{
			name:    "multiple fields with prefix",
			pattern: `{"title":"fake-api-test-.{5}","name":"fake-api-test-.{5}"}`,
			input:   `{"title":"fake-api-test-abcde","name":"fake-api-test-12345"}`,
			want:    true,
		},
		{
			name: "escaped prefix match",
			// Note: hyphens need to be escaped in regex
			pattern: `{"name":"fake\-api\-test\-.{5}"}`,
			input:   `{"name":"fake-api-test-12345"}`,
			want:    true,
		},
		{
			name:    "full request pattern",
			pattern: `{"author":"Author","name":"fake-api-test-.{5}","title":"fake-api-test-.{5}","uuid":"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"}`,
			input:   `{"author":"Author","name":"fake-api-test-12345","title":"fake-api-test-abcde","uuid":"550e8400-e29b-41d4-a716-446655440000"}`,
			want:    true,
		},
		{
			name:    "flexible field order",
			pattern: `{"name":"fake-api-test-.{5}","title":".*"}|{"title":".*","name":"fake-api-test-.{5}"}`,
			input:   `{"title":"test","name":"fake-api-test-12345"}`,
			want:    true,
		},
		{
			name:    "non-matching prefix",
			pattern: `{"name":"fake-api-test-.{5}"}`,
			input:   `{"name":"different-prefix-12345"}`,
			want:    false,
		},
		{
			name:    "wrong length after prefix",
			pattern: `{"name":"fake-api-test-.{5}"}`,
			input:   `{"name":"fake-api-test-123456"}`, // 6 chars instead of 5
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchers.RegexMatch(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("Pattern: %q\nInput:   %q\nExpected: %v\nGot: %v",
					tc.pattern, tc.input, tc.want, got)
			}
		})
	}
}

func TestDatetimeRegexMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name:    "simple datetime json",
			pattern: `{"created_at":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)"}`,
			input:   `{"created_at":"2024-02-23T10:30:00Z"}`,
			want:    true,
		},
		{
			name:    "simple date json",
			pattern: `{"created_at":"(\d{4}-\d{2}-\d{2})"}`,
			input:   `{"created_at":"2024-02-23"}`,
			want:    true,
		},
		{
			name:    "actual processed pattern",
			pattern: `{"created_at":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z|\d{4}-\d{2}-\d{2})"}`,
			input:   `{"created_at":"2024-02-23T10:30:00Z"}`,
			want:    true,
		},
		{
			name:    "complete test data",
			pattern: `{"created_at":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z|\d{4}-\d{2}-\d{2})","due_date":"(\d{4}-\d{2}-\d{2})","name":"test-task","updated_at":"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z|\d{4}-\d{2}-\d{2})"}`,
			input:   `{"created_at":"2024-02-23T10:30:00Z","due_date":"2024-02-23","name":"test-task","updated_at":"2024-02-23T10:30:00Z"}`,
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchers.RegexMatch(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("RegexMatch(%q, %q) = %v, want %v", tc.pattern, tc.input, got, got)
			}
		})
	}
}

func TestActivityPatternRegexMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name: "exact activity payload match with proper regex",
			// This uses a proper regex pattern that will match the date field with any value
			pattern: `{"completed":false,"dueDate":"[^"]*","id":1,"title":"fake-api-test-[a-zA-Z0-9]{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-abcde"}`,
			want:    true,
		},
		{
			name:    "activity with wrong title format - too long",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-abcdef"}`,
			want:    false,
		},
		{
			name:    "activity with wrong title format - too short",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-abc"}`,
			want:    false,
		},
		{
			name:    "activity with wrong title prefix",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"different-prefix-abcde"}`,
			want:    false,
		},
		{
			name:    "activity with different id",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":2,"title":"fake-api-test-abcde"}`,
			want:    false,
		},
		{
			name:    "activity with different completed value",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":true,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-abcde"}`,
			want:    false,
		},
		{
			name:    "activity with missing field",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-abcde"}`,
			want:    false,
		},
		{
			name:    "activity with actual client format - activity-12345",
			pattern: `{"completed":false,"dueDate":"*","id":1,"title":"fake-api-test-.{5}"}`,
			input:   `{"completed":false,"dueDate":"2025-02-25T10:30:00Z","id":1,"title":"fake-api-test-activity-12345"}`,
			want:    false, // This will fail because the pattern expects exactly 5 chars after the prefix
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchers.RegexMatch(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("\nPattern: %q\nInput:   %q\nExpected: %v\nGot: %v", tc.pattern, tc.input, tc.want, got)
			}
		})
	}
}

func TestRandomTitlePatternMatch(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		input   string
		want    bool
	}{
		{
			name:    "client randomTitleForResource(activity) pattern",
			pattern: `{"title":"fake-api-test-activity-.{5}"}`,
			input:   `{"title":"fake-api-test-activity-12345"}`,
			want:    true,
		},
		{
			name:    "current pattern vs client format",
			pattern: `{"title":"fake-api-test-.{5}"}`,
			input:   `{"title":"fake-api-test-activity-12345"}`,
			want:    false,
		},
		{
			name:    "more general pattern to match any client format",
			pattern: `{"title":"fake-api-test-.*"}`,
			input:   `{"title":"fake-api-test-activity-12345"}`,
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchers.RegexMatch(tc.pattern, tc.input)
			if got != tc.want {
				t.Errorf("\nPattern: %q\nInput:   %q\nExpected: %v\nGot: %v", tc.pattern, tc.input, tc.want, got)
			}
		})
	}
}
