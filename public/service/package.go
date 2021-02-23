// Package service provides a high level API for creating a Benthos service,
// registering custom plugin components, customizing the availability of native
// components, and running it.
package service

import (
	"context"
)

// Closer is implemented by components that support stopping and cleaning up
// their underlying resources.
type Closer interface {
	// Close the component, blocks until either the underlying resources are
	// cleaned up or the context is cancelled. Returns an error if the context
	// is cancelled.
	Close(ctx context.Context) error
}

// Resources provides access to service-wide resources.
type Resources interface {
	// Metrics returns an aggregator preset with the current component context.
	// Metrics() Metrics

	// Logger returns a logger preset with the current component context.
	// Logger() Logger

	// AccessInput(ctx context.Context, name string, fn func(i Input)) error
	// AccessOutput(ctx context.Context, name string, fn func(o Output)) error
	AccessProcessor(ctx context.Context, name string, fn func(p Processor)) error
	AccessCache(ctx context.Context, name string, fn func(c Cache)) error
	AccessRatelimit(ctx context.Context, name string, fn func(r RateLimit)) error
}

// ConfigConstructor returns a struct containing configuration fields for a
// plugin implementation with default values. This struct will be
// marshalled/unmarshalled as YAML using gopkg.in/yaml.v3.
type ConfigConstructor func() interface{}
