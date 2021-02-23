package service

// Logger is an interface provided to Benthos components for structured logging.
type Logger interface {
	// Debugf logs a debug message using fmt.Sprintf when args are specified.
	Debugf(template string, args ...interface{})

	// Infof logs an info message using fmt.Sprintf when args are specified.
	Infof(template string, args ...interface{})

	// Warnf logs a warning message using fmt.Sprintf when args are specified.
	Warnf(template string, args ...interface{})

	// Errorf logs an error message using fmt.Sprintf when args are specified.
	Errorf(template string, args ...interface{})

	// With adds a variadic set of fields to a logger. Each field must consist
	// of a string key and a value of any type. An odd number of key/value pairs
	// will therefore result in malformed log messages, but should never panic.
	With(keyValuePairs ...interface{}) Logger
}
