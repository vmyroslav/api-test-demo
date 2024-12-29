package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Matcher represents a Hoverfly matcher
type Matcher struct {
	Matcher string `json:"matcher"`
	Value   string `json:"value"`
}

// Request represents a Hoverfly request
type Request struct {
	Path        []Matcher           `json:"path"`
	Method      []Matcher           `json:"method"`
	Destination []Matcher           `json:"destination"`
	Scheme      []Matcher           `json:"scheme"`
	Body        []Matcher           `json:"body"`
	Headers     map[string][]string `json:"headers,omitempty"`
}

// Response represents a Hoverfly response
type Response struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
	Templated   bool                `json:"templated"`
}

// Pair represents a request-response pair
type Pair struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// HoverflySimulation represents the structure of a Hoverfly simulation file
type HoverflySimulation struct {
	Data struct {
		Pairs []Pair `json:"pairs"`
	} `json:"data"`
	Meta struct {
		SchemaVersion   string            `json:"schemaVersion"`
		TimeExported    string            `json:"timeExported,omitempty"`
		HoverflyVersion string            `json:"hoverflyVersion,omitempty"`
		Templates       map[string]string `json:"templates,omitempty"`
	} `json:"meta"`
}

func ProcessSimulation(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	var simulation HoverflySimulation
	if err := json.Unmarshal(data, &simulation); err != nil {
		return err
	}

	// Keep track of counters for each endpoint
	endpointCounters := make(map[string]int)

	// Process each request-response pair
	for i := range simulation.Data.Pairs {
		pair := &simulation.Data.Pairs[i]

		// Get or create counter for this endpoint
		endpoint := fmt.Sprintf("%s-%s",
			pair.Request.Path[0].Value,
			pair.Request.Method[0].Value)
		endpointCounters[endpoint]++
		counter := endpointCounters[endpoint]

		// Process request bodies with random strings and timestamps
		if len(pair.Request.Body) > 0 && pair.Request.Body[0].Value != "" {
			requestBody := pair.Request.Body[0].Value

			// Create a glob pattern for the entire request
			if strings.Contains(requestBody, "fake-api-test") {
				// Replace exact random strings with glob pattern, including the counter
				re := regexp.MustCompile(`"(fake-api-test-[^"]*-title-)[a-zA-Z0-9]{5}"`)
				requestBody = re.ReplaceAllString(requestBody, fmt.Sprintf(`"$1*_%d"`, counter))

				// Replace timestamps if present
				timeRe := regexp.MustCompile(`"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2}))"`)
				if timeRe.MatchString(requestBody) {
					requestBody = timeRe.ReplaceAllString(requestBody, `"*"`)
				}

				// Update the request matcher
				pair.Request.Body[0].Value = requestBody
				pair.Request.Body[0].Matcher = "glob"

				// Process response for 200 status with JSON body
				if pair.Response.Status == 200 && strings.Contains(pair.Response.Body, "{") {
					var respData map[string]interface{}
					if err := json.Unmarshal([]byte(pair.Response.Body), &respData); err == nil {
						// Update response fields with request values
						if _, ok := respData["title"]; ok {
							respData["title"] = "{{ Request.Body 'jsonpath' '$.title' }}"
						}
						if _, ok := respData["dueDate"]; ok {
							respData["dueDate"] = "{{ Request.Body 'jsonpath' '$.dueDate' }}"
						}
						if _, ok := respData["firstName"]; ok {
							respData["firstName"] = "{{ Request.Body 'jsonpath' '$.firstName' }}"
						}
						if _, ok := respData["lastName"]; ok {
							respData["lastName"] = "{{ Request.Body 'jsonpath' '$.lastName' }}"
						}

						// Marshal the modified response body
						if newRespBody, err := json.Marshal(respData); err == nil {
							pair.Response.Body = string(newRespBody)
							pair.Response.Templated = true
						}
					}
				}
			}
		}
	}

	// Add template matchers if they don't exist
	if simulation.Meta.Templates == nil {
		simulation.Meta.Templates = make(map[string]string)
	}
	simulation.Meta.Templates["timestamp"] = "2006-01-02T15:04:05Z07:00"

	output, err := json.MarshalIndent(simulation, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, output, 0o644)
}
