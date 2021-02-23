package service

import (
	"context"
	"errors"

	"github.com/Jeffail/benthos/v3/lib/message"
	"github.com/Jeffail/benthos/v3/lib/types"
)

// Message represents a single discrete message passing through a Benthos
// pipeline. It is safe to mutate the message via Set methods, but the
// underlying byte data should not be edited directly.
type Message interface {
	Context() context.Context
	SetContext(ctx context.Context)

	// Messages can store arbitrary key/value metadata values.
	MetaGet(string) (string, bool)
	MetaSet(string, string)
	MetaDelete(string)

	// Iterate each metadata value by providing a closure, which will be called
	// with the key and value. To stop iterating, return false from the closure.
	MetaWalk(func(string, string) bool)

	Bytes() ([]byte, bool)
	Structured() (interface{}, bool)

	SetBytes([]byte)
	SetStructured(interface{})
}

//------------------------------------------------------------------------------

// CopyMessage creates a copy of a message that is safe to mutate without
// mutating the original. Both messages will share a context, and therefore a
// tracing ID, if one has been associated with them.
//
// Note that this does not perform a deep copy of the byte or structured
// contents of the message, and therefore it is not safe to perform inline
// mutations on those values.
func CopyMessage(msg Message) Message {
	if agMsg, ok := msg.(*airGapMessage); ok {
		return &airGapMessage{
			m:          agMsg.m.Copy(),
			partCopied: true,
		}
	}

	// If the message implementation is something we don't recognize then we're
	// forced to perform a manual copy. This isn't ideal as it means testing for
	// JSON parsability, which means actually attempting a parse.
	part := message.NewPart(nil)

	s, ok := msg.Structured()
	if ok {
		part.SetJSON(s)
	} else {
		b, _ := msg.Bytes()
		part.Set(b)
	}

	msg.MetaWalk(func(k, v string) bool {
		part.Metadata().Set(k, v)
		return true
	})

	return &airGapMessage{
		m:          message.WithContext(msg.Context(), part),
		partCopied: true,
	}
}

//------------------------------------------------------------------------------

// Converts a types.Part into a Message and also ensures the underlying message
// is lazily cloned in the event of any mutations.
type airGapMessage struct {
	m          types.Part
	partCopied bool
}

func newAirGapMessage(m types.Part) Message {
	return &airGapMessage{m, false}
}

func (a *airGapMessage) ensureCopied() {
	if !a.partCopied {
		a.m = a.m.Copy()
		a.partCopied = true
	}
}

func (a *airGapMessage) Context() context.Context {
	return message.GetContext(a.m)
}

func (a *airGapMessage) Bytes() ([]byte, bool) {
	b := a.m.Get()
	return b, len(b) > 0
}

func (a *airGapMessage) Structured() (interface{}, bool) {
	i, err := a.m.JSON()
	if err != nil {
		return nil, false
	}
	return i, true
}

func (a *airGapMessage) SetBytes(b []byte) {
	a.ensureCopied()
	a.m.Set(b)
}

func (a *airGapMessage) SetStructured(i interface{}) {
	a.ensureCopied()
	a.m.SetJSON(i)
}

func (a *airGapMessage) SetContext(ctx context.Context) {
	a.m = message.WithContext(ctx, a.m)
}

func (a *airGapMessage) MetaGet(key string) (string, bool) {
	v := a.m.Metadata().Get(key)
	return v, len(v) > 0
}

func (a *airGapMessage) MetaSet(key, value string) {
	a.ensureCopied()
	a.m.Metadata().Set(key, value)
}

func (a *airGapMessage) MetaDelete(key string) {
	a.ensureCopied()
	a.m.Metadata().Delete(key)
}

func (a *airGapMessage) MetaWalk(fn func(string, string) bool) {
	_ = a.m.Metadata().Iter(func(k, v string) error {
		if !fn(k, v) {
			return errors.New("stop")
		}
		return nil
	})
}
