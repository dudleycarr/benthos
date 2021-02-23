package service

import (
	"testing"

	"github.com/Jeffail/benthos/v3/lib/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageCopyAirGap(t *testing.T) {
	p := message.NewPart([]byte("hello world"))
	p.Metadata().Set("foo", "bar")
	g1 := newAirGapMessage(p)
	g2 := CopyMessage(g1)

	b := p.Get()
	v := p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok := g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	g2.SetBytes([]byte("and now this"))
	g2.MetaSet("foo", "baz")

	b = p.Get()
	v = p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "and now this", string(b))
	assert.Equal(t, "baz", v)

	g1.SetBytes([]byte("but not this"))
	g1.MetaSet("foo", "buz")

	b = p.Get()
	v = p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "but not this", string(b))
	assert.Equal(t, "buz", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "and now this", string(b))
	assert.Equal(t, "baz", v)
}

func TestMessageCopyAirGapUnknownType(t *testing.T) {
	type notAirGapMsg struct {
		*airGapMessage
	}

	p := message.NewPart([]byte("hello world"))
	p.Metadata().Set("foo", "bar")

	g1 := notAirGapMsg{newAirGapMessage(p).(*airGapMessage)}
	g2 := CopyMessage(g1)

	b := p.Get()
	v := p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok := g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	g2.SetBytes([]byte("and now this"))
	g2.MetaSet("foo", "baz")

	b = p.Get()
	v = p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "and now this", string(b))
	assert.Equal(t, "baz", v)

	g1.SetBytes([]byte("but not this"))
	g1.MetaSet("foo", "buz")

	b = p.Get()
	v = p.Metadata().Get("foo")
	assert.Equal(t, "hello world", string(b))
	assert.Equal(t, "bar", v)

	b, ok = g1.Bytes()
	v, _ = g1.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "but not this", string(b))
	assert.Equal(t, "buz", v)

	b, ok = g2.Bytes()
	v, _ = g2.MetaGet("foo")
	require.True(t, ok)
	assert.Equal(t, "and now this", string(b))
	assert.Equal(t, "baz", v)
}

func TestMessageQuery(t *testing.T) {
	p := message.NewPart([]byte(`{"foo":"bar"}`))
	p.Metadata().Set("foo", "bar")
	p.Metadata().Set("bar", "baz")
	g1 := newAirGapMessage(p)

	b, ok := g1.Bytes()
	assert.True(t, ok)
	assert.Equal(t, `{"foo":"bar"}`, string(b))

	s, ok := g1.Structured()
	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, s)

	m, ok := g1.MetaGet("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", m)

	seen := map[string]string{}
	g1.MetaWalk(func(k, v string) bool {
		seen[k] = v
		return false
	})
	assert.Len(t, seen, 1)

	seen = map[string]string{}
	g1.MetaWalk(func(k, v string) bool {
		seen[k] = v
		return true
	})
	assert.Equal(t, map[string]string{
		"foo": "bar",
		"bar": "baz",
	}, seen)
}

func TestMessageMutate(t *testing.T) {
	p := message.NewPart([]byte(`not a json doc`))
	p.Metadata().Set("foo", "bar")
	p.Metadata().Set("bar", "baz")
	g1 := newAirGapMessage(p)

	_, ok := g1.Structured()
	assert.False(t, ok)

	g1.SetStructured(map[string]interface{}{
		"foo": "bar",
	})
	assert.Equal(t, "not a json doc", string(p.Get()))

	s, ok := g1.Structured()
	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
	}, s)

	g1.SetBytes([]byte("foo bar baz"))
	assert.Equal(t, "not a json doc", string(p.Get()))

	_, ok = g1.Structured()
	assert.False(t, ok)

	b, ok := g1.Bytes()
	assert.True(t, ok)
	assert.Equal(t, "foo bar baz", string(b))

	g1.MetaDelete("foo")

	seen := map[string]string{}
	g1.MetaWalk(func(k, v string) bool {
		seen[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"bar": "baz"}, seen)

	g1.MetaSet("foo", "new bar")

	seen = map[string]string{}
	g1.MetaWalk(func(k, v string) bool {
		seen[k] = v
		return true
	})
	assert.Equal(t, map[string]string{"foo": "new bar", "bar": "baz"}, seen)
}
