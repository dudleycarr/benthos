package input

import (
	"github.com/Jeffail/benthos/v3/internal/docs"
	"github.com/Jeffail/benthos/v3/lib/input/reader"
	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/message/batch"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/types"
	"github.com/Jeffail/benthos/v3/lib/util/tls"
)

//------------------------------------------------------------------------------

func init() {
	Constructors[TypeNSQ] = TypeSpec{
		constructor: fromSimpleConstructor(NewNSQ),
		Summary: `
Subscribe to an NSQ instance topic and channel.`,
		Description: `
        ### Metadata

        This input adds the following metadata fields to each message:

        ` + "``` text" + `
        - nsq_attempts
        - nsq_id
        - nsq_nsqd_address
        - nsq_timestamp
        ` + "```" + `

        You can access these metadata fields using
        [function interpolation](/docs/configuration/interpolation#metadata).`,
		FieldSpecs: docs.FieldSpecs{
			func() docs.FieldSpec {
				b := batch.FieldSpec()
				b.Deprecated = true
				return b
			}(),
			docs.FieldCommon("nsqd_tcp_addresses", "A list of nsqd addresses to connect to.").Array(),
			docs.FieldCommon("lookupd_http_addresses", "A list of nsqlookupd addresses to connect to.").Array(),
			tls.FieldSpec(),
			docs.FieldCommon("topic", "The topic to consume from."),
			docs.FieldCommon("channel", "The channel to consume from."),
			docs.FieldCommon("user_agent", "A user agent to assume when connecting."),
			docs.FieldCommon("max_in_flight", "The maximum number of pending messages to consume at any given time."),
			docs.FieldCommon("max_attempts", "The maximum number of attempts to successfully process a message."),
		},
		Categories: []Category{
			CategoryServices,
		},
	}
}

//------------------------------------------------------------------------------

// NewNSQ creates a new NSQ input type.
func NewNSQ(conf Config, mgr types.Manager, log log.Modular, stats metrics.Type) (Type, error) {
	var n reader.Async
	var err error
	if n, err = reader.NewNSQ(conf.NSQ, log, stats); err != nil {
		return nil, err
	}
	if n, err = reader.NewAsyncBatcher(conf.NSQ.Batching, n, mgr, log, stats); err != nil {
		return nil, err
	}
	n = reader.NewAsyncBundleUnacks(n)
	return NewAsyncReader(TypeNSQ, true, n, log, stats)
}

//------------------------------------------------------------------------------
