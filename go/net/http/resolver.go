// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"net/http"
	"net/url"

	"github.com/searKing/golang/go/net/resolver"
	url_ "github.com/searKing/golang/go/net/url"
)

// RequestWithTarget reset Host in url.Url by resolver.Target
func RequestWithTarget(req *http.Request, target string) error {
	u2, err := url_.ResolveWithTarget(req.Context(), req.URL, target)
	if err != nil {
		return err
	}
	req.URL = u2
	return nil
}

// ProxyFuncWithTargetOrDefault builds a proxy function from the given string, which should
// represent a target that can be used as a proxy. It performs basic
// sanitization of the Target and returns any error encountered.
func ProxyFuncWithTargetOrDefault(target string, def func(req *http.Request) (*url.URL, error)) (func(req *http.Request) (*url.URL, error), error) {
	if target == "" {
		return def, nil
	}
	return func(req *http.Request) (*url.URL, error) {
		reqURL := req.URL
		if target == "" {
			return nil, nil
		}
		address, err := resolver.ResolveOneAddr(req.Context(), target)
		if err != nil {
			return nil, err
		}
		var proxy url.URL
		if reqURL.Scheme == "https" {
			proxy.Scheme = "https"
			proxy.Host = address.Addr
		} else if reqURL.Scheme == "http" {
			proxy.Scheme = "http"
			proxy.Host = address.Addr
		} else {
			return nil, nil
		}

		return &proxy, nil
	}, nil
}