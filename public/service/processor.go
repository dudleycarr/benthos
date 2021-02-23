package service

import (
	"context"
)

// ProcessorFunc is a minimal implementation of a Benthos processor, able to
// mutate a message into one or more resulting messages, or return an error if
// the message could not be processed. If zero messages are returned and the
// error is nil then the message is filtered.
type ProcessorFunc func(context.Context, Message) ([]Message, error)

// Processor is a Benthos processor implementation that works against single
// messages.
type Processor interface {
	// Process a message into one or more resulting messages, or return an error
	// if the message could not be processed. If zero messages are returned and
	// the error is nil then the message is filtered.
	Process(context.Context, Message) ([]Message, error)
}

// ProcessorCloser is an interface implemented by processors that allow shutting
// down and removing underlying resources.
type ProcessorCloser interface {
	Processor
	Closer
}

// ProcessorConstructor is a func that's provided a configuration type and
// access to a service manager and must return an instantiation of a processor
// based on the config, or an error.
type ProcessorConstructor func(label string, conf interface{}, mgr Resources) (ProcessorCloser, error)

//------------------------------------------------------------------------------

// BatchProcessorFunc is a minimal implementation of a Benthos batch processor,
// able to mutate a batch of messages or return an error if the messages could
// not be processed.
type BatchProcessorFunc func(context.Context, []Message) ([]Message, error)

// BatchProcessor is a Benthos processor implementation that works against
// batches of messages, which allows windowed processing and includes methods
// for cleanly shutting down.
type BatchProcessor interface {
	// Process a batch of messages or return an error if the messages could not
	// be processed.
	ProcessBatch(context.Context, []Message) ([]Message, error)
}

// BatchProcessorCloser is an interface implemented by processors that allow
// shutting down and removing underlying resources.
type BatchProcessorCloser interface {
	BatchProcessor
	Closer
}

// BatchProcessorConstructor is a func that's provided a configuration type and
// access to a service manager and must return an instantiation of a processor
// based on the config, or an error.
type BatchProcessorConstructor func(label string, conf interface{}, mgr Resources) (BatchProcessorCloser, error)
