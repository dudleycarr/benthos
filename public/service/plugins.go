package service

import (
	"github.com/Jeffail/benthos/v3/internal/interop"
	"github.com/Jeffail/benthos/v3/lib/cache"
	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/types"
)

// RegisterCache attempts to register a new cache plugin by providing an
// optional constructor for a configuration struct for the plugin as well as a
// constructor for the cache itself. The constructor will be called for each
// instantiation of the component within a config.
//
// If you plugin does not require configuration then it is valid to set confCtor
// to nil.
func RegisterCache(name string, confCtor ConfigConstructor, cCtor CacheConstructor) error {
	cache.RegisterPlugin(name, cache.PluginConfigConstructor(confCtor), func(
		config interface{},
		mgr types.Manager,
		logger log.Modular,
		metrics metrics.Type,
	) (types.Cache, error) {
		label := interop.GetLabel(mgr)
		c, err := cCtor(label, config, resourcesFromManager(mgr))
		if err != nil {
			return nil, err
		}
		return newAirGapCache(c), nil
	})
	return nil
}

// RegisterRateLimit attempts to register a new rate limit plugin by providing
// an optional constructor for a configuration struct for the plugin as well as
// a constructor for the rate limit itself. The constructor will be called for
// each instantiation of the component within a config.
func RegisterRateLimit(name string, confCtor ConfigConstructor, cCtor RateLimitConstructor) error {
	panic("not implemented")
}

// RegisterProcessorFunc attempts to register a new processor plugin by
// providing a func to operate on each message. Depending on a configuration
// this func could be called by multiple parallel components, and therefore
// should be thread safe.
func RegisterProcessorFunc(name string, fn ProcessorFunc) error {
	panic("not implemented")
}

// RegisterProcessor attempts to register a new processor plugin by providing an
// optional constructor for a configuration struct for the plugin as well as a
// constructor for the processor itself. The constructor will be called for each
// instantiation of the component within a config.
func RegisterProcessor(name string, confCtor ConfigConstructor, cCtor ProcessorConstructor) error {
	panic("not implemented")
}

// RegisterBatchProcessorFunc attempts to register a new batch processor plugin
// by providing a func to operate on each message. Depending on a configuration
// this func could be called by multiple parallel components, and therefore
// should be thread safe.
func RegisterBatchProcessorFunc(name string, fn BatchProcessorFunc) error {
	panic("not implemented")
}

// RegisterBatchProcessor attempts to register a new processor plugin by
// providing an optional constructor for a configuration struct for the plugin
// as well as a constructor for the processor itself. The constructor will be
// called for each instantiation of the component within a config.
func RegisterBatchProcessor(name string, confCtor ConfigConstructor, cCtor BatchProcessorConstructor) error {
	panic("not implemented")
}
