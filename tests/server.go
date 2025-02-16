package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
)

type TestServer struct {
	server *httptest.Server
	mux    *http.ServeMux

	attempts atomic.Int32
}

func NewTestServer() *TestServer {
	ts := &TestServer{}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/Activities/", ts.handleGetActivity)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		ts.ResetAttempts()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Attempts reset"))
	})

	ts.server = httptest.NewServer(mux)
	ts.mux = mux

	return ts
}

func (ts *TestServer) URL() string {
	return ts.server.URL
}

func (ts *TestServer) Close() {
	ts.server.Close()
}

func (ts *TestServer) ResetAttempts() {
	ts.attempts.Store(0)
}

func (ts *TestServer) handleGetActivity(w http.ResponseWriter, r *http.Request) {
	attempt := ts.attempts.Add(1)

	// Extract the id from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idStr := pathParts[4]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Simulate different responses based on attempt number
	switch attempt {
	case 1:
		// First attempt - Service Unavailable
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Service Temporarily Unavailable",
			"code":  "SERVER_BUSY",
		})
	case 2:
		// Second attempt - Gateway Timeout
		w.WriteHeader(http.StatusGatewayTimeout)
	case 3:
		// Third attempt - Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"id":        int32(id),
			"title":     "Test Activity",
			"dueDate":   "2025-02-20T09:58:59.009Z",
			"completed": false,
		})
	default:
		// Any subsequent attempts - Bad Request
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Too many attempts",
		})
	}
}

func (ts *TestServer) Handler() http.Handler {
	return ts.mux
}
