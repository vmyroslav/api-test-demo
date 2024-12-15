package tests

import "net/http"

func NewHttpClient() *http.Client {
	return &http.Client{}
}

func HostURL() string {
	return "https://fakerestapi.azurewebsites.net"
}
