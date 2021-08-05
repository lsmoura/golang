// Code generated by "go-option -type Pick"; DO NOT EDIT.

package resolver

// A PickOption sets options.
type PickOption interface {
	apply(*Pick)
}

// EmptyPickOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyPickOption struct{}

func (EmptyPickOption) apply(*Pick) {}

// PickOptionFunc wraps a function that modifies Pick into an
// implementation of the PickOption interface.
type PickOptionFunc func(*Pick)

func (f PickOptionFunc) apply(do *Pick) {
	f(do)
}

// sample code for option, default for nothing to change
func _PickOptionWithDefault() PickOption {
	return PickOptionFunc(func(*Pick) {
		// nothing to change
	})
}
func (o *Pick) ApplyOptions(options ...PickOption) *Pick {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
