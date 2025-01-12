package config

import "net/http"

type Config interface {
	Routes() []Route
	Patterns() []Pattern
}

type Route interface {
	Path() string
	Methods() []string
	Request() Request
	Response() Response
}

type Request interface {
	Body() string
	Headers() http.Header
	AdditionalConfig() map[string]any
}

type Response interface {
	Body() string
	Headers() http.Header
	AdditionalConfig() map[string]any
}

type Pattern interface {
	Type() string
	Formats() []string
	Pattern() *string
	ReplaceWith() string
}
