package service

import "context"

// AckFunc is a common function returned by inputs that must be called once for
// each message consumed. This function ensures that the source of the message
// receives the acknowledgement (if applicable).
type AckFunc func(ctx context.Context, err error) error

// Reader is an interface implemented by Benthos inputs. Calls to Read block
// until either a message has been received, or the context is cancelled. When
// successful both a message and an acknowledgement func are returned, the ack
// func MUST be called at least once for every message, otherwise the input
// could stall.
type Reader interface {
	Read(context.Context) (Message, AckFunc, error)
}

// ReaderCloser is an interface implemented by a Reader that supports closing
// down the underlying connection and resources.
type ReaderCloser interface {
	Reader
	Closer
}
