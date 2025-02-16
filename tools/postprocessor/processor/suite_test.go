package processor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/processor"

	hoverfly "github.com/SpectoLabs/hoverfly/core"
	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

type ProcessorTestSuite struct {
	suite.Suite
	hf     *hoverfly.Hoverfly
	srv    *httptest.Server
	client *http.Client
	tmpDir string
	logger *slog.Logger

	mu sync.Mutex
}

func (s *ProcessorTestSuite) SetupTest() {
	s.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	s.tmpDir = s.T().TempDir()
	s.srv = s.startTestServer()
	s.hf = startHoverfly(s.T())
	s.client = createHoverflyClient()
}

func (s *ProcessorTestSuite) TearDownTest() {
	if s.srv != nil {
		s.srv.Close()
	}

	if s.hf != nil {
		s.hf.StopProxy()
	}
}

func (s *ProcessorTestSuite) startTestServer() *httptest.Server {
	mux := http.NewServeMux()

	// Handler for requests with UUIDs
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	})

	// Handler for static endpoint
	mux.HandleFunc("/api/v1/books/30", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        30,
			"title":     "Original Title",
			"dueDate":   time.Now().Format(time.RFC3339),
			"completed": true,
		})
	})

	// Handler for prefix requests
	mux.HandleFunc("/api/prefixes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Parse request body and responde with the same data and add id property equal to 1
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		req["id"] = 1
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
	})

	return httptest.NewServer(mux)
}

func (s *ProcessorTestSuite) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		fmt.Sprintf("%s%s", s.srv.URL, path),
		strings.NewReader(string(data)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return s.client.Do(req)
}

func (s *ProcessorTestSuite) restartHoverfly() {
	if s.hf != nil {
		s.hf.StopProxy()
	}

	s.hf = startHoverfly(s.T())
	s.client = createHoverflyClient()
}

func (s *ProcessorTestSuite) exportSimulation() v2.SimulationViewV5 {
	simulation, err := s.hf.GetSimulation()
	require.NoError(s.T(), err)

	// Save original simulation for reference
	origSimFile := filepath.Join(s.tmpDir, "original_simulation.json")
	err = writeHoverflySimulation(&simulation, origSimFile)
	require.NoError(s.T(), err)

	return simulation
}

func (s *ProcessorTestSuite) processSimulation(sim *v2.SimulationViewV5, cfg *config.Config) string {
	proc := processor.New(cfg, s.logger)
	require.NoError(s.T(), proc.Process(sim))

	// Save processed simulation
	processedSimFile := filepath.Join(s.tmpDir, "processed_simulation.json")
	require.NoError(s.T(), writeHoverflySimulation(sim, processedSimFile))

	return processedSimFile
}

func (s *ProcessorTestSuite) setMux(mux *http.ServeMux) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.srv.Close()
	s.srv = httptest.NewServer(mux)
}

func createHoverflyClient() *http.Client {
	proxyURL := fmt.Sprintf("http://localhost:%d", pport)
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(mustParseURL(proxyURL)),
		},
		Timeout: 5 * time.Second,
	}
}

func findFreePort() int {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return 50000 + time.Now().Nanosecond()%10000
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

func mustParseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	return parsed
}

var pport = 8500

func startHoverfly(t *testing.T) *hoverfly.Hoverfly {
	cfg := hoverfly.InitSettings()
	pport = findFreePort()
	cfg.AdminPort = strconv.Itoa(findFreePort())
	cfg.ProxyPort = strconv.Itoa(pport)

	hf := hoverfly.NewHoverflyWithConfiguration(cfg)
	err := hf.StartProxy()
	require.NoError(t, err)

	return hf
}

func writeHoverflySimulation(sim *v2.SimulationViewV5, filePath string) error {
	payload, err := json.MarshalIndent(sim, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling simulation: %w", err)
	}

	if err := os.WriteFile(filePath, payload, 0o644); err != nil {
		return fmt.Errorf("writing simulation file: %w", err)
	}

	return nil
}
