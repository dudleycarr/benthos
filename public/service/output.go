package service

import "context"

// BatchWriter is an interface implemented by Benthos outputs that support
// synchronous batch writes. Calls to WriteBatch block until either the messages
// have been successfully or unsuccessfully sent, or the context is cancelled.
type BatchWriter interface {
	WriteBatch(context.Context, []Message) error
}

// BatchWriterCloser is an interface implemented by a BatchWriter that supports
// closing down the underlying connection and resources.
type BatchWriterCloser interface {
	BatchWriter
	Closer
}

// Writer is an interface implemented by Benthos outputs that support
// synchronous single message writes. Calls to Write block until either the
// message has been successfully or unsuccessfully sent, or the context is
// cancelled.
type Writer interface {
	Write(context.Context, Message) error
}

// WriterCloser is an interface implemented by a Writer that supports closing
// down the underlying connection and resources.
type WriterCloser interface {
	Writer
	Closer
}
