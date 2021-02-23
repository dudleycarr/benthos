package service

// MetricCounter is a representation of a single metric counter. Interactions
// with this counter are thread safe. The number of label values for each call
// should exactly match the number of label keys specified with Counter.
type MetricCounter interface {
	// Incr increments a counter by an amount. The number of label values should
	// exactly match the number of label keys specified when this metric was
	// created.
	Incr(count int64, labelValues ...string) error
}

// MetricTimer is a representation of a single metric timer. Interactions with
// this timer are thread safe. The number of label values for each call should
// exactly match the number of label keys specified with Timer.
type MetricTimer interface {
	// Timing sets a timing metric. The number of label values should exactly
	// match the number of label keys specified when this metric was created.
	Timing(delta int64, labelValues ...string) error
}

// MetricGauge is a representation of a single metric gauge. Interactions with
// this gauge are thread safe. The number of label values for each call should
// exactly match the number of label keys specified with Gauge.
type MetricGauge interface {
	// Set sets the value of a gauge metric. The number of label values should
	// exactly match the number of label keys specified when this metric was
	// created.
	Set(value int64, labelValues ...string) error

	// Incr increments a gauge by an amount. The number of label values should
	// exactly match the number of label keys specified when this metric was
	// created.
	Incr(count int64, labelValues ...string) error

	// Decr decrements a gauge by an amount. The number of label values should
	// exactly match the number of label keys specified when this metric was
	// created.
	Decr(count int64, labelValues ...string) error
}

//------------------------------------------------------------------------------

// Metrics is an interface for metrics aggregation.
type Metrics interface {
	// Counter returns a counter stat for a given name. Dynamic labels can be
	// added to this metric by specifying any number of label keys. Calls to the
	// resulting counter must have a number of label values exactly matching the
	// number of keys.
	Counter(name string, labelKeys ...string) MetricCounter

	// Timer returns a timer stat for a given name. Dynamic labels can be added
	// to this metric by specifying any number of label keys. Calls to the
	// resulting timer must have a number of label values exactly matching the
	// number of keys.
	Timer(name string, labelKeys ...string) MetricTimer

	// Gauge returns a gauge stat for a given name. Dynamic labels can be added
	// to this metric by specifying any number of label keys. Calls to the
	// resulting gauge must have a number of label values exactly matching the
	// number of keys.
	Gauge(name string, labelKeys ...string) MetricGauge

	// With adds a variadic set of labels to a metrics aggregator. Each label
	// must consist of a key and value pair. An odd number of key/value pairs
	// will therefore result in malformed metric labels, but should never panic.
	With(labelKeyValues ...string) Metrics

	Closer
}
