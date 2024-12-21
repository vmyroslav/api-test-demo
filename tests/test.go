package tests

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"testing"
)

func NewHttpClient(t testing.TB) *http.Client {
	t.Helper()

	//if os.Getenv("HTTP_PROXY") == "" {
	//	return http.DefaultClient
	//}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   "localhost:8500",
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
