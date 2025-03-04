package tests

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func NewHttpClient(t *testing.T) *http.Client {
	t.Helper()

	hoverflyAddr := os.Getenv("HOVERFLY_PROXY")
	if hoverflyAddr == "" {
		return http.DefaultClient
	}

	// We will run tests in parallel only if we are using Hoverfly.
	// Just to avoid redundant load for public API.
	t.Parallel()

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   hoverflyAddr,
			}),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip certificate verification when using Hoverfly
			},
		},
	}

	return client
}
