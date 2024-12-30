package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
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

	// Process each request-response pair
	for i := range simulation.Data.Pairs {
		pair := &simulation.Data.Pairs[i]

		// Check if request body contains our patterns
		if len(pair.Request.Body) > 0 && pair.Request.Body[0].Value != "" {
			requestBody := pair.Request.Body[0].Value

			// Only process if body contains our patterns
			if strings.Contains(requestBody, "fake-api-test") || containsTimePattern(requestBody) {
				var requestData map[string]interface{}
				if err := json.Unmarshal([]byte(requestBody), &requestData); err == nil {
					modifiedFields := make(map[string]bool)

					// Process each field in the request
					for field, value := range requestData {
						if strVal, ok := value.(string); ok {
							// Handle fake-api-test pattern
							if strings.Contains(strVal, "fake-api-test") {
								requestData[field] = strings.Replace(strVal, strVal[len(strVal)-5:], "*", 1)
								modifiedFields[field] = true
							}
							// Handle timestamp pattern
							if isTimeFormat(strVal) {
								requestData[field] = "*"
								modifiedFields[field] = true
							}
						}
					}

					if len(modifiedFields) > 0 {
						// Update request body with glob pattern
						if newBody, err := json.Marshal(requestData); err == nil {
							pair.Request.Body[0].Value = string(newBody)
							pair.Request.Body[0].Matcher = "glob"

							// Only template the response if we modified the request
							if pair.Response.Status == 200 && strings.Contains(pair.Response.Body, "{") {
								var respData map[string]interface{}
								if err := json.Unmarshal([]byte(pair.Response.Body), &respData); err == nil {
									// Template only the fields that were modified in request
									for field := range modifiedFields {
										if _, exists := respData[field]; exists {
											respData[field] = fmt.Sprintf("{{ Request.Body 'jsonpath' '$.%s' }}", field)
										}
									}
									if newRespBody, err := json.Marshal(respData); err == nil {
										pair.Response.Body = string(newRespBody)
										pair.Response.Templated = true
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Add template matchers
	if simulation.Meta.Templates == nil {
		simulation.Meta.Templates = make(map[string]string)
	}
	simulation.Meta.Templates["timestamp"] = "2006-01-02T15:04:05Z07:00"

	// Remove duplicates
	uniquePairs := make([]Pair, 0)
	seen := make(map[string]bool)

	for _, pair := range simulation.Data.Pairs {
		key := createRequestSignature(pair.Request)
		if !seen[key] {
			seen[key] = true
			uniquePairs = append(uniquePairs, pair)
		}
	}

	simulation.Data.Pairs = uniquePairs

	output, err := json.MarshalIndent(simulation, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, output, 0o644)
}

func containsTimePattern(s string) bool {
	timePattern := regexp.MustCompile(`"\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})"`)
	return timePattern.MatchString(s)
}

func isTimeFormat(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

func createRequestSignature(req Request) string {
	var sig strings.Builder

	// Add method
	if len(req.Method) > 0 {
		sig.WriteString("method:" + req.Method[0].Value + ";")
	}

	// Add path
	if len(req.Path) > 0 {
		sig.WriteString("path:" + req.Path[0].Value + ";")
	}

	// Add body if exists
	if len(req.Body) > 0 && req.Body[0].Value != "" {
		sig.WriteString("body:" + req.Body[0].Value + ";")
	}

	return sig.String()
}
