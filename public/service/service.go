package service

// Environment defines an executable Benthos environment, including a registry
// of all available components and plugins, and provides methods for executing
// that environment as a service.
type Environment struct {
	// TODO: Add component registrar, includes plugins if nil after opts then we
	// use standard set.
	// mgr *manager.Type
}

// OptFunc is a function signature used to define options for the Environment
// constructor.
type OptFunc func(t *Environment)

// New constructs a new Benthos service environment.
func New(opts ...OptFunc) (*Environment, error) {
	t := &Environment{}
	for _, opt := range opts {
		opt(t)
	}
	return t, nil
}

//------------------------------------------------------------------------------

// RunCLI executes Benthos as a CLI service, allowing users to specify a
// configuration via flags. This is how a standard distribution of Benthos
// operates.
//
// This call blocks until either the pipeline shuts down or a termination signal
// is received.
func (t *Environment) RunCLI() {
	panic("not implemented")
}
