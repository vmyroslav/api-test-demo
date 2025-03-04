package config

// EndpointRule defines a static response for a specific endpoint
type EndpointRule struct {
	// Method is the HTTP method
	Method string `yaml:"method"`

	// Status is the HTTP status code to return
	Status int `yaml:"status"`

	// Path is the endpoint path, supports glob patterns
	// Example: "/api/v1/activities/*"
	Path string `yaml:"path"`

	// StaticResponse is a JSON string to return for this endpoint
	StaticResponse string `yaml:"static_response"`
}
