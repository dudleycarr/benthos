package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type closableRateLimit struct {
	next   time.Duration
	err    error
	closed bool
}

func (c *closableRateLimit) Access(ctx context.Context) (time.Duration, error) {
	return c.next, c.err
}

func (c *closableRateLimit) Close(ctx context.Context) error {
	c.closed = true
	return nil
}

func TestRateLimitAirGapShutdown(t *testing.T) {
	rl := &closableRateLimit{
		next: time.Second,
	}
	agrl := newAirGapRateLimit(rl)

	tout, err := agrl.Access()
	assert.NoError(t, err)
	assert.Equal(t, time.Second, tout)

	rl.next = time.Millisecond
	rl.err = errors.New("test error")

	tout, err = agrl.Access()
	assert.EqualError(t, err, "test error")
	assert.Equal(t, time.Millisecond, tout)

	err = agrl.WaitForClose(time.Millisecond * 5)
	assert.EqualError(t, err, "action timed out")
	assert.False(t, rl.closed)

	agrl.CloseAsync()
	err = agrl.WaitForClose(time.Millisecond * 5)
	assert.NoError(t, err)
	assert.True(t, rl.closed)
}
