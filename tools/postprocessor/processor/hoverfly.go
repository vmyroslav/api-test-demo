package processor

// Matcher represents a Hoverfly matcher
type Matcher struct {
	Matcher string `json:"matcher"`
	Value   string `json:"value"`
}

// Request represents a Hoverfly request
type Request struct {
	Path        []Matcher           `json:"path"`
	Method      []Matcher           `json:"method"`
	Destination []Matcher           `json:"destination"`
	Scheme      []Matcher           `json:"scheme"`
	Body        []Matcher           `json:"body"`
	Headers     map[string][]string `json:"headers,omitempty"`
}

// Response represents a Hoverfly response
type Response struct {
	Status      int                 `json:"status"`
	Body        string              `json:"body"`
	EncodedBody bool                `json:"encodedBody"`
	Headers     map[string][]string `json:"headers"`
	Templated   bool                `json:"templated"`
}

// Pair represents a request-response pair
type Pair struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Simulation represents the structure of a Hoverfly simulation file
type Simulation struct {
	Data struct {
		Pairs []Pair `json:"pairs"`
	} `json:"data"`
	Meta struct {
		SchemaVersion   string            `json:"schemaVersion"`
		TimeExported    string            `json:"timeExported,omitempty"`
		HoverflyVersion string            `json:"hoverflyVersion,omitempty"`
		Templates       map[string]string `json:"templates,omitempty"`
	} `json:"meta"`
}
