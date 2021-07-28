// Code generated by "go-option -type doWithBackoff"; DO NOT EDIT.

package http

// A DoWithBackoffOption sets options.
type DoWithBackoffOption interface {
	apply(*doWithBackoff)
}

// EmptyDoWithBackoffOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyDoWithBackoffOption struct{}

func (EmptyDoWithBackoffOption) apply(*doWithBackoff) {}

// DoWithBackoffOptionFunc wraps a function that modifies doWithBackoff into an
// implementation of the DoWithBackoffOption interface.
type DoWithBackoffOptionFunc func(*doWithBackoff)

func (f DoWithBackoffOptionFunc) apply(do *doWithBackoff) {
	f(do)
}

// sample code for option, default for nothing to change
func _DoWithBackoffOptionWithDefault() DoWithBackoffOption {
	return DoWithBackoffOptionFunc(func(*doWithBackoff) {
		// nothing to change
	})
}
func (o *doWithBackoff) ApplyOptions(options ...DoWithBackoffOption) *doWithBackoff {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}
