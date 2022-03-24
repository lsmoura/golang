// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by "go-option -type=wrapper"; DO NOT EDIT.
// Install go-option by "go get -u github.com/searKing/golang/tools/go-option"

package otlpsql

// A WrapperOption sets options.
type WrapperOption interface {
	apply(*wrapper)
}

// EmptyWrapperOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyWrapperOption struct{}

func (EmptyWrapperOption) apply(*wrapper) {}

// WrapperOptionFunc wraps a function that modifies wrapper into an
// implementation of the WrapperOption interface.
type WrapperOptionFunc func(*wrapper)

func (f WrapperOptionFunc) apply(do *wrapper) {
	f(do)
}
func (o *wrapper) ApplyOptions(options ...WrapperOption) *wrapper {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

// WithAllowRoot sets AllowRoot in wrapper.
func WithAllowRoot(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.AllowRoot = v
	})
}

// WithPing sets Ping in wrapper.
func WithPing(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.Ping = v
	})
}

// WithRowsNext sets RowsNext in wrapper.
func WithRowsNext(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.RowsNext = v
	})
}

// WithRowsClose sets RowsClose in wrapper.
func WithRowsClose(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.RowsClose = v
	})
}

// WithRowsAffected sets RowsAffected in wrapper.
func WithRowsAffected(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.RowsAffected = v
	})
}

// WithLastInsertID sets LastInsertID in wrapper.
func WithLastInsertID(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.LastInsertID = v
	})
}

// WithQuery sets Query in wrapper.
func WithQuery(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.Query = v
	})
}

// WithQueryParams sets QueryParams in wrapper.
func WithQueryParams(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.QueryParams = v
	})
}

// WithInstanceName sets InstanceName in wrapper.
func WithInstanceName(v string) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.InstanceName = v
	})
}

// WithDisableErrSkip sets DisableErrSkip in wrapper.
func WithDisableErrSkip(v bool) WrapperOption {
	return WrapperOptionFunc(func(o *wrapper) {
		o.DisableErrSkip = v
	})
}
