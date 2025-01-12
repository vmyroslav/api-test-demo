package v0

type ConfigV0 struct {
	Version  string     `yaml:"version"`
	Settings Settings   `yaml:"settings"`
	Patterns []Patterns `yaml:"patterns"`
	Rules    []Rules    `yaml:"rules"`
}
type Settings struct {
	CaseSensitive bool `yaml:"case_sensitive"`
	Debug         bool `yaml:"debug"`
}
type Patterns struct {
	Type        string   `yaml:"type"`
	Formats     []string `yaml:"formats,omitempty"`
	ReplaceWith string   `yaml:"replace_with"`
	Pattern     string   `yaml:"pattern,omitempty"`
	KeepChars   int      `yaml:"keep_chars,omitempty"`
}
type Title struct {
	Type        string `yaml:"type"`
	Pattern     string `yaml:"pattern"`
	KeepChars   int    `yaml:"keep_chars"`
	ReplaceWith string `yaml:"replace_with"`
}
type DueDate struct {
	Type        string   `yaml:"type"`
	Formats     []string `yaml:"formats"`
	ReplaceWith string   `yaml:"replace_with"`
}
type Fields struct {
	Title   Title   `yaml:"title"`
	DueDate DueDate `yaml:"due_date"`
}
type Body struct {
	Fields Fields `yaml:"fields"`
}
type Request struct {
	Body Body `yaml:"body"`
}
type Response struct {
	TemplateMatchedFields bool `yaml:"template_matched_fields"`
}
type XRequestID struct {
	Type string `yaml:"type"`
}
type Headers struct {
	XRequestID XRequestID `yaml:"X-Request-ID"`
}

type Rules struct {
	Path     string   `yaml:"path"`
	Methods  []string `yaml:"methods"`
	Request  Request  `yaml:"request,omitempty"`
	Response Response `yaml:"response,omitempty"`
}
