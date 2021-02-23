package service

import (
	"context"
	"time"

	"github.com/Jeffail/benthos/v3/lib/types"
)

// RateLimit is an interface implemented by Benthos rate limits.
type RateLimit interface {
	// Access the rate limited resource. Returns a duration or an error if the
	// rate limit check fails. The returned duration is either zero (meaning the
	// resource may be accessed) or a reasonable length of time to wait before
	// requesting again.
	Access(context.Context) (time.Duration, error)

	Close(context.Context) error
}

// RateLimitConstructor is a func that's provided a configuration type and
// access to a service manager and must return an instantiation of a rate limit
// based on the config, or an error.
type RateLimitConstructor func(label string, conf interface{}, mgr Resources) (RateLimit, error)

//------------------------------------------------------------------------------

// Implements types.RateLimit
type airGapRateLimit struct {
	c RateLimit

	ctx  context.Context
	done func()
}

func newAirGapRateLimit(c RateLimit) types.RateLimit {
	ctx, done := context.WithCancel(context.Background())
	return &airGapRateLimit{c, ctx, done}
}

func (a *airGapRateLimit) Access() (time.Duration, error) {
	return a.c.Access(context.Background())
}

func (a *airGapRateLimit) CloseAsync() {
	go func() {
		if err := a.c.Close(context.Background()); err == nil {
			a.done()
		}
	}()
}

func (a *airGapRateLimit) WaitForClose(tout time.Duration) error {
	select {
	case <-a.ctx.Done():
	case <-time.After(tout):
		return types.ErrTimeout
	}
	return nil
}
