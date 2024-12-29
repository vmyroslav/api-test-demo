package tests

import (
	"crypto/tls"
	"math/rand"
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
		Timeout: 5 * time.Second,
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

// ToPtr converts a value of any type to a pointer.
func ToPtr[T any](v T) *T {
	return &v
}

// RandomString generates a random string of the given length.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}
