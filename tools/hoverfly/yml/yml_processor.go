package yml

func ProcessVersion(version string) Config {
	panic("implement me")
	return nil
}

type Config interface {
	Routes() []Route
}

type Route struct {
	Method  string
	Path    string
	Matcher string
}
