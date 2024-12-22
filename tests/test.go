package tests

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func NewHttpClient(t testing.TB) *http.Client {
	t.Helper()

	var hoverflyAddr, ok = os.LookupEnv("HOVERFLY_PROXY")
	if !ok {
		return http.DefaultClient
	}

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

func HostURL() string {
	return "https://fakerestapi.azurewebsites.net"
}
