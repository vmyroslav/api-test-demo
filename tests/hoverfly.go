package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type TestHelper struct {
	hoverflyCmd *exec.Cmd
	t           *testing.T
}

func NewTestHelper(t *testing.T) *TestHelper {
	t.Helper()
	helper := &TestHelper{t: t}

	if useHoverfly() {
		helper.startHoverfly()
		t.Cleanup(helper.stopHoverfly)
	}

	return helper
}

func (h *TestHelper) startHoverfly() {
	simulationPath := filepath.Join("testdata", "simulations")
	err := os.MkdirAll(simulationPath, 0755)
	if err != nil {
		h.t.Fatalf("Failed to create simulations directory: %v", err)
	}

	// Start Hoverfly in the background
	h.hoverflyCmd = exec.Command("hoverfly", "-webserver")
	h.hoverflyCmd.Stdout = os.Stdout
	h.hoverflyCmd.Stderr = os.Stderr

	if err := h.hoverflyCmd.Start(); err != nil {
		h.t.Fatalf("Failed to start Hoverfly: %v", err)
	}

	// Wait for Hoverfly to be ready
	h.waitForHoverfly()

	// Set up Hoverfly mode
	if os.Getenv("CAPTURE_MODE") == "true" {
		h.setHoverflyMode("capture")
		// Register cleanup to export simulation
		h.t.Cleanup(func() {
			h.exportSimulation()
		})
	} else {
		h.setHoverflyMode("simulate")
		h.importSimulation()
	}
}

func (h *TestHelper) stopHoverfly() {
	if h.hoverflyCmd != nil && h.hoverflyCmd.Process != nil {
		h.hoverflyCmd.Process.Kill()
	}
}

func (h *TestHelper) waitForHoverfly() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			h.t.Fatal("Timeout waiting for Hoverfly to start")
		default:
			resp, err := http.Get("http://localhost:8888/api/v2/hoverfly")
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (h *TestHelper) setHoverflyMode(mode string) {
	cmd := exec.Command("hoverctl", "mode", mode)
	if err := cmd.Run(); err != nil {
		h.t.Fatalf("Failed to set Hoverfly mode to %s: %v", mode, err)
	}
}

func (h *TestHelper) importSimulation() {
	simFile := h.getSimulationPath()
	if _, err := os.Stat(simFile); err == nil {
		cmd := exec.Command("hoverctl", "import", simFile)
		if err := cmd.Run(); err != nil {
			h.t.Fatalf("Failed to import simulation: %v", err)
		}
	}
}

func (h *TestHelper) exportSimulation() {
	simFile := h.getSimulationPath()
	cmd := exec.Command("hoverctl", "export", simFile)
	if err := cmd.Run(); err != nil {
		h.t.Fatalf("Failed to export simulation: %v", err)
	}
}

func (h *TestHelper) getSimulationPath() string {
	return filepath.Join("testdata", "simulations", fmt.Sprintf("%s.json", h.t.Name()))
}

func (h *TestHelper) GetHTTPClient() *http.Client {
	if useHoverfly() {
		return &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(h.GetProxyURL()),
			},
		}
	}
	return http.DefaultClient
}

func (h *TestHelper) GetProxyURL() *url.URL {
	proxyURL, _ := url.Parse("http://localhost:8500")
	return proxyURL
}

func (h *TestHelper) GetBaseURL() string {
	if useHoverfly() {
		return "http://localhost:8500"
	}
	return os.Getenv("API_BASE_URL")
}

func useHoverfly() bool {
	return os.Getenv("USE_HOVERFLY") == "true"
}
