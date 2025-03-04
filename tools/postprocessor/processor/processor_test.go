//go:build integration

package processor_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/modes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestProcessor(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}

func (s *ProcessorTestSuite) TestStaticEndpointReplacement() {
	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	resp, err := s.client.Get(fmt.Sprintf("%s/api/v1/books/30", s.srv.URL))
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	cfg := &config.Config{
		Settings: config.Settings{Debug: false},
		Endpoints: []config.EndpointRule{
			{
				Method: "GET",
				Path:   "/api/v1/books/30",
				Status: http.StatusCreated,
				StaticResponse: `{
                    "id": 122,
                    "title": "John Doe",
                    "dueDate": "2025-02-19T17:12:52.127Z",
                    "completed": false
                }`,
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))
	require.Equal(s.T(), modes.Simulate, s.hf.GetMode().Mode) // Verify current mode is Simulate

	resp, err = s.client.Get(fmt.Sprintf("%s/api/v1/books/30", s.srv.URL))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	var result map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))

	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	assert.Equal(s.T(), float64(122), result["id"])
	assert.Equal(s.T(), "John Doe", result["title"])
	assert.Equal(s.T(), "2025-02-19T17:12:52.127Z", result["dueDate"])
	assert.Equal(s.T(), false, result["completed"])
}

func (s *ProcessorTestSuite) TestUUIDPatternMatching() {
	mux := http.NewServeMux()

	// Handler for requests with UUIDs
	mux.HandleFunc("/api/uuids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       uuid.New().String(),
			"secondId": uuid.New().String(),
			"name":     "test-resource",
		})
	})

	s.setMux(mux)

	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	// Make request with UUID data
	captureData := map[string]interface{}{
		"id":       "550e8400-e29b-41d4-a716-446655440000",
		"secondId": "987fcdeb-51e4-12d3-a456-426614174000",
		"name":     "test-resource",
	}

	resp, err := s.makeRequest("POST", "/api/uuids", captureData)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	cfg := &config.Config{
		Settings: config.Settings{Debug: true},
		Patterns: []config.Pattern{
			{
				Type: config.UUIDPattern,
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Verify the processed simulation has correct UUID patterns
	var savedSim v2.SimulationViewV5
	savedData, err := os.ReadFile(processedSimFile)
	require.NoError(s.T(), err)
	err = json.Unmarshal(savedData, &savedSim)
	require.NoError(s.T(), err)

	require.Len(s.T(), savedSim.RequestResponsePairs, 1)
	pair := savedSim.RequestResponsePairs[0]

	// Verify request body has UUID pattern
	require.Len(s.T(), pair.RequestMatcher.Body, 1)
	require.Equal(s.T(), "regex", pair.RequestMatcher.Body[0].Matcher)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))
	require.Equal(s.T(), modes.Simulate, s.hf.GetMode().Mode)

	// Test with different UUIDs
	testCases := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "different_uuids",
			data: map[string]interface{}{
				"id":       "7c9e6679-7425-40de-944b-e07fc1f90ae7",
				"secondId": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				"name":     "test-resource",
			},
		},
		{
			name: "another_uuid_set",
			data: map[string]interface{}{
				"id":       "8a7e5b3c-9f1d-44e2-b1aa-4e9a44a5d7e9",
				"secondId": "123e4567-e89b-12d3-a456-426614174000",
				"name":     "test-resource",
			},
		},
	}

	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			resp, err := s.makeRequest("POST", "/api/uuids", tc.data)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)

			var result map[string]interface{}
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

			// Verify UUIDs match input
			assert.Equal(t, tc.data["id"], result["id"])
			assert.Equal(t, tc.data["external_id"], result["external_id"])
			assert.Equal(t, tc.data["name"], result["name"])
		})
	}
}

func (s *ProcessorTestSuite) TestUUIDSpecificReplacement() {
	mux := http.NewServeMux()

	// Handler for requests with UUIDs
	mux.HandleFunc("/api/uuids", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       uuid.New().String(),
			"secondId": uuid.New().String(),
			"name":     "test-resource",
		})
	})

	s.setMux(mux)

	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	resp, err := s.makeRequest(http.MethodGet, "/api/uuids", nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	// Configure specific UUID replacements
	cfg := &config.Config{
		Settings: config.Settings{Debug: true},
		Patterns: []config.Pattern{
			{
				Type:        config.UUIDPattern,
				ReplaceWith: "00000000-0000-0000-0000-000000000001",
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Verify the processed simulation
	var savedSim v2.SimulationViewV5
	savedData, err := os.ReadFile(processedSimFile)
	require.NoError(s.T(), err)
	err = json.Unmarshal(savedData, &savedSim)
	require.NoError(s.T(), err)

	require.Len(s.T(), savedSim.RequestResponsePairs, 1)
	pair := savedSim.RequestResponsePairs[0]

	// Verify request body has exact matcher for first UUID and regex for second
	require.Len(s.T(), pair.RequestMatcher.Body, 1)

	var bodyPattern map[string]interface{}
	err = json.Unmarshal([]byte(pair.RequestMatcher.Body[0].Value.(string)), &bodyPattern)
	require.NoError(s.T(), err)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))
	require.Equal(s.T(), modes.Simulate, s.hf.GetMode().Mode)

	resp, err = s.makeRequest("GET", "/api/uuids", nil)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))

	// Verify response
	assert.Equal(s.T(), "00000000-0000-0000-0000-000000000001", result["id"])
	assert.Equal(s.T(), "00000000-0000-0000-0000-000000000001", result["secondId"])
	assert.Equal(s.T(), "test-resource", result["name"])
}

func (s *ProcessorTestSuite) TestPrefixPattern() {
	var (
		mux        = http.NewServeMux()
		fakePrefix = "fake-api-test-"
		randomFunc = func(length int) string {
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

			seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			b := make([]byte, length)
			for i := range b {
				b[i] = charset[seededRand.Intn(len(charset))]
			}

			return string(b)
		}
		id        = uuid.New().String()
		randTitle = randomFunc(5)
		randName  = randomFunc(5)
		author    = "Author"
	)

	// Handler for requests with UUIDs
	mux.HandleFunc("/api/prefixes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(reqBody); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})

	s.setMux(mux)

	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	// Make request with prefix data
	captureData := map[string]interface{}{
		"uuid":   id,
		"name":   fmt.Sprintf("%s%s", fakePrefix, randName),
		"title":  fmt.Sprintf("%s%s", fakePrefix, randTitle),
		"author": author,
	}

	resp, err := s.makeRequest("POST", "/api/prefixes", captureData)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	cfg := &config.Config{
		Settings: config.Settings{Debug: true},
		Patterns: []config.Pattern{
			{
				Type:    config.PrefixPattern,
				Pattern: "fake-api-test-",
				Length:  5,
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))
	require.Equal(s.T(), modes.Simulate, s.hf.GetMode().Mode) // Verify current mode is Simulate

	// Make request with prefix data
	testData := map[string]interface{}{
		"uuid":   id,
		"name":   fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"title":  fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"author": author,
	}

	fmt.Printf("Sending request data: %+v", testData)

	resp, err = s.makeRequest("POST", "/api/prefixes", testData)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))

	assert.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	assert.Equal(s.T(), testData["uuid"], result["uuid"])
	assert.Equal(s.T(), testData["name"], result["name"])
	assert.Equal(s.T(), testData["title"], result["title"])
	assert.Equal(s.T(), testData["author"], result["author"])
}

func (s *ProcessorTestSuite) TestDatetimePattern() {
	var (
		mux      = http.NewServeMux()
		taskName = "test-task"
	)

	// Handler for datetime requests
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reqBody)
	})

	s.setMux(mux)

	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	// Create capture request data with fixed test data
	captureData := map[string]interface{}{
		"name":       taskName,
		"created_at": "2024-02-23T10:30:00Z",
		"updated_at": "2024-02-23T10:30:00Z",
		"due_date":   "2024-02-23",
	}

	resp, err := s.makeRequest("POST", "/api/tasks", captureData)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	cfg := &config.Config{
		Settings: config.Settings{Debug: true},
		Patterns: []config.Pattern{
			{
				Type: config.DatetimePattern,
				Formats: []string{
					"2006-01-02T15:04:05Z",
					"2006-01-02",
				},
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))

	// Test data with different dates but same format
	testData := map[string]interface{}{
		"name":       taskName,
		"created_at": "2024-02-24T15:45:00Z",
		"updated_at": "2024-02-24T15:45:00Z",
		"due_date":   "2024-02-25",
	}

	resp, err = s.makeRequest("POST", "/api/tasks", testData)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var result map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&result))

	// Verify the response
	assert.Equal(s.T(), taskName, result["name"])
	assert.Equal(s.T(), testData["created_at"], result["created_at"])
	assert.Equal(s.T(), testData["updated_at"], result["updated_at"])
	assert.Equal(s.T(), testData["due_date"], result["due_date"])
}

func (s *ProcessorTestSuite) TestCombinedScenarios() {
	var (
		mux        = http.NewServeMux()
		fakePrefix = "fake-api-test-"
		randomFunc = func(length int) string {
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
			b := make([]byte, length)
			for i := range b {
				b[i] = charset[seededRand.Intn(len(charset))]
			}
			return string(b)
		}
	)

	// Handler for complex resource
	mux.HandleFunc("/api/resources", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reqBody)
	})

	// Handler for static response
	mux.HandleFunc("/api/static/123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        123,
			"name":      "Original Name",
			"createdAt": time.Now().Format(time.RFC3339),
		})
	})

	s.setMux(mux)

	// Step 1: Capture Phase
	require.NoError(s.T(), s.hf.SetMode(modes.Capture))

	// First request: Complex resource with all pattern types
	firstRequestData := map[string]interface{}{
		"id":             uuid.New().String(),
		"externalId":     uuid.New().String(),
		"name":           fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"title":          fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"createdAt":      time.Now().Format(time.RFC3339),
		"dueDate":        time.Now().Format("2006-01-02"),
		"lastModifiedAt": time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}

	resp, err := s.makeRequest("POST", "/api/resources", firstRequestData)
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Second request: Static endpoint
	resp, err = s.client.Get(fmt.Sprintf("%s/api/static/123", s.srv.URL))
	require.NoError(s.T(), err)
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Step 2: Get and Process Simulation
	simulation := s.exportSimulation()

	cfg := &config.Config{
		Settings: config.Settings{Debug: true},
		Patterns: []config.Pattern{
			{
				Type: config.UUIDPattern,
			},
			{
				Type:    config.PrefixPattern,
				Pattern: fakePrefix,
				Length:  5,
			},
			{
				Type: config.DatetimePattern,
				Formats: []string{
					time.RFC3339,
					"2006-01-02",
				},
			},
		},
		Endpoints: []config.EndpointRule{
			{
				Method: "GET",
				Path:   "/api/static/123",
				StaticResponse: `{
					"id": 456,
					"name": "Static Response",
					"createdAt": "2025-01-01T00:00:00Z"
				}`,
			},
		},
	}

	processedSimFile := s.processSimulation(&simulation, cfg)

	// Step 3: Switch to Simulation Mode
	s.restartHoverfly()
	require.NoError(s.T(), s.hf.Import(processedSimFile))
	require.NoError(s.T(), s.hf.SetMode(modes.Simulate))

	// Test 1: Complex resource with different values
	testData := map[string]interface{}{
		"id":             uuid.New().String(),
		"externalId":     uuid.New().String(),
		"name":           fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"title":          fmt.Sprintf("%s%s", fakePrefix, randomFunc(5)),
		"createdAt":      time.Now().Format(time.RFC3339),
		"dueDate":        time.Now().Format("2006-01-02"),
		"lastModifiedAt": time.Now().Format(time.RFC3339),
	}

	resp, err = s.makeRequest("POST", "/api/resources", testData)
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)

	var complexResult map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&complexResult))

	// Verify all fields were handled correctly
	assert.Equal(s.T(), testData["id"], complexResult["id"], "UUID field should match")
	assert.Equal(s.T(), testData["externalId"], complexResult["externalId"], "Second UUID field should match")
	assert.Equal(s.T(), testData["name"], complexResult["name"], "Prefixed name should match")
	assert.Equal(s.T(), testData["title"], complexResult["title"], "Prefixed title should match")
	assert.Equal(s.T(), testData["createdAt"], complexResult["createdAt"], "DateTime should match")
	assert.Equal(s.T(), testData["dueDate"], complexResult["dueDate"], "Date should match")
	assert.Equal(s.T(), testData["lastModifiedAt"], complexResult["lastModifiedAt"], "DateTime should match")

	// Test 2: Static endpoint
	resp, err = s.client.Get(fmt.Sprintf("%s/api/static/123", s.srv.URL))
	require.NoError(s.T(), err)
	defer resp.Body.Close()

	require.Equal(s.T(), http.StatusOK, resp.StatusCode)

	var staticResult map[string]interface{}
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&staticResult))

	// Verify static response
	assert.Equal(s.T(), float64(456), staticResult["id"])
	assert.Equal(s.T(), "Static Response", staticResult["name"])
	assert.Equal(s.T(), "2025-01-01T00:00:00Z", staticResult["createdAt"])
}
